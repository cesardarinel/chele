package handlers

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/cesardarinel/chele-api/internal/middleware"
	"github.com/cesardarinel/chele-api/internal/models"
	"github.com/golang-jwt/jwt/v5"
	"github.com/jmoiron/sqlx"
)

type AuthHandler struct {
	DB      *sqlx.DB
	JWTSecret string
}

func NewAuthHandler(db *sqlx.DB, secret string) *AuthHandler {
	return &AuthHandler{DB: db, JWTSecret: secret}
}

type loginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req loginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		jsonError(w, "invalid request", http.StatusBadRequest)
		return
	}

	var user models.User
	err := h.DB.Get(&user, "SELECT * FROM auth_user WHERE username=? AND is_active=1", req.Username)
	if err != nil {
		jsonError(w, "invalid credentials", http.StatusUnauthorized)
		return
	}

	if !middleware.CheckDjangoPassword(req.Password, user.Password) {
		jsonError(w, "invalid credentials", http.StatusUnauthorized)
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id":  user.ID,
		"username": user.Username,
		"email":    user.Email,
		"exp":      time.Now().Add(24 * time.Hour).Unix(),
	})
	tokenStr, err := token.SignedString([]byte(h.JWTSecret))
	if err != nil {
		jsonError(w, "token error", http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"token": tokenStr,
		"user":  user,
	})
}

type registerRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Name     string `json:"name"`
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req registerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		jsonError(w, "invalid request", http.StatusBadRequest)
		return
	}
	if req.Username == "" || req.Email == "" || req.Password == "" {
		jsonError(w, "username, email, password required", http.StatusBadRequest)
		return
	}

	var count int
	h.DB.Get(&count, "SELECT COUNT(*) FROM auth_user WHERE username=?", req.Username)
	if count > 0 {
		jsonError(w, "username already exists", http.StatusConflict)
		return
	}

	hash := middleware.MakeDjangoPassword(req.Password)
	names := strings.SplitN(req.Name, " ", 2)
	firstName, lastName := names[0], ""
	if len(names) > 1 {
		lastName = names[1]
	}

	result, err := h.DB.Exec(
		`INSERT INTO auth_user
		 (password,last_login,is_superuser,username,last_name,email,is_staff,is_active,date_joined,first_name)
		 VALUES (?,NULL,0,?,?,?,0,1,datetime('now'),?)`,
		hash, req.Username, req.Email, firstName, lastName,
	)
	if err != nil {
		jsonError(w, "registration failed: "+err.Error(), http.StatusInternalServerError)
		return
	}

	uid, _ := result.LastInsertId()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id":  uid,
		"username": req.Username,
		"email":    req.Email,
		"exp":      time.Now().Add(24 * time.Hour).Unix(),
	})
	tokenStr, _ := token.SignedString([]byte(h.JWTSecret))

	writeJSON(w, http.StatusCreated, map[string]interface{}{
		"token":    tokenStr,
		"user_id":  uid,
		"username": req.Username,
		"email":    req.Email,
	})
}

func (h *AuthHandler) Me(w http.ResponseWriter, r *http.Request) {
	uid := r.Context().Value(middleware.UserIDKey).(int)
	var user models.User
	if err := h.DB.Get(&user, "SELECT * FROM auth_user WHERE id=?", uid); err != nil {
		jsonError(w, "user not found", http.StatusNotFound)
		return
	}
	jsonOK(w, user)
}

package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/cesardarinel/chele-api/internal/config"
	"github.com/cesardarinel/chele-api/internal/middleware"
	"github.com/cesardarinel/chele-api/internal/testutil"
	"github.com/golang-jwt/jwt/v5"
)

func setupAuthHandler(t *testing.T) (*AuthHandler, *config.Config) {
	t.Helper()
	db := testutil.SetupDB()
	t.Cleanup(func() { db.Close() })
	cfg := &config.Config{JWTSecert: "test-secret"}
	h := NewAuthHandler(db, cfg.JWTSecert)
	return h, cfg
}

func seedUser(t *testing.T, h *AuthHandler) (string, string, int) {
	t.Helper()
	password := "testpass123"
	hash := middleware.MakeDjangoPassword(password)
	username := "test@example.com"
	result, err := h.DB.Exec(
		"INSERT INTO auth_user (password,username,email,is_superuser,is_staff,is_active,first_name,last_name,date_joined) VALUES (?,?,?,0,0,1,'Test','User',datetime('now'))",
		hash, username, username,
	)
	if err != nil {
		t.Fatal(err)
	}
	uid64, _ := result.LastInsertId()
	return username, password, int(uid64)
}

func TestLogin_Success(t *testing.T) {
	h, _ := setupAuthHandler(t)
	username, password, _ := seedUser(t, h)

	body, _ := json.Marshal(map[string]string{"username": username, "password": password})
	req := httptest.NewRequest("POST", "/api/auth/login", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	h.Login(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}
	var resp map[string]interface{}
	json.NewDecoder(w.Body).Decode(&resp)
	if resp["token"] == nil {
		t.Error("expected token in response")
	}
}

func TestLogin_WrongPassword(t *testing.T) {
	h, _ := setupAuthHandler(t)
	username, _, _ := seedUser(t, h)

	body, _ := json.Marshal(map[string]string{"username": username, "password": "wrongpass"})
	req := httptest.NewRequest("POST", "/api/auth/login", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	h.Login(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", w.Code)
	}
}

func TestLogin_UserNotFound(t *testing.T) {
	h, _ := setupAuthHandler(t)
	body, _ := json.Marshal(map[string]string{"username": "noone@test.com", "password": "pass"})
	req := httptest.NewRequest("POST", "/api/auth/login", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	h.Login(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", w.Code)
	}
}

func TestRegister_CreatesUser(t *testing.T) {
	h, _ := setupAuthHandler(t)
	body, _ := json.Marshal(map[string]string{
		"username": "newuser@test.com",
		"email":    "newuser@test.com",
		"password": "pass123",
		"name":     "New User",
	})
	req := httptest.NewRequest("POST", "/api/auth/register", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	h.Register(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d: %s", w.Code, w.Body.String())
	}
	var count int
	h.DB.Get(&count, "SELECT COUNT(*) FROM auth_user WHERE username=?", "newuser@test.com")
	if count != 1 {
		t.Error("user was not created")
	}
}

func TestRegister_Duplicate(t *testing.T) {
	h, _ := setupAuthHandler(t)
	seedUser(t, h)

	body, _ := json.Marshal(map[string]string{
		"username": "test@example.com",
		"email":    "test@example.com",
		"password": "pass123",
	})
	req := httptest.NewRequest("POST", "/api/auth/register", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	h.Register(w, req)

	if w.Code != http.StatusConflict {
		t.Errorf("expected 409, got %d", w.Code)
	}
}

func TestRegister_Validation(t *testing.T) {
	h, _ := setupAuthHandler(t)
	body, _ := json.Marshal(map[string]string{"username": "", "email": "", "password": ""})
	req := httptest.NewRequest("POST", "/api/auth/register", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	h.Register(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}

func TestMe_WithValidToken(t *testing.T) {
	h, cfg := setupAuthHandler(t)
	_, _, uid := seedUser(t, h)

	tokenStr, err := generateTestToken(uid, cfg.JWTSecert)
	if err != nil {
		t.Fatal(err)
	}

	// Set the context manually as the middleware would
	ctx := context.WithValue(context.Background(), middleware.UserIDKey, uid)
	req := httptest.NewRequest("GET", "/api/auth/me", nil).WithContext(ctx)
	req.Header.Set("Authorization", "Bearer "+tokenStr)
	w := httptest.NewRecorder()
	h.Me(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}
	var resp map[string]interface{}
	json.NewDecoder(w.Body).Decode(&resp)
	if resp["username"] != "test@example.com" {
		t.Errorf("unexpected username: %v", resp["username"])
	}
}

func generateTestToken(userID int, secret string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userID,
	})
	return token.SignedString([]byte(secret))
}

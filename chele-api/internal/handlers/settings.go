package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/cesardarinel/chele-api/internal/middleware"
	"github.com/cesardarinel/chele-api/internal/models"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type SettingsHandler struct{ DB *sqlx.DB }

func NewSettingsHandler(db *sqlx.DB) *SettingsHandler { return &SettingsHandler{DB: db} }

func (h *SettingsHandler) Profile(w http.ResponseWriter, r *http.Request) {
	uid := r.Context().Value(middleware.UserIDKey).(int)
	var user models.User
	if err := h.DB.Get(&user, "SELECT id,username,email,first_name,last_name FROM auth_user WHERE id=?", uid); err != nil {
		jsonError(w, "not found", http.StatusNotFound)
		return
	}
	jsonOK(w, user)
}

func (h *SettingsHandler) UpdateProfile(w http.ResponseWriter, r *http.Request) {
	uid := r.Context().Value(middleware.UserIDKey).(int)
	var req struct {
		FirstName *string `json:"first_name"`
		LastName  *string `json:"last_name"`
		Email     *string `json:"email"`
	}
	json.NewDecoder(r.Body).Decode(&req)
	if req.FirstName != nil { h.DB.Exec("UPDATE auth_user SET first_name=? WHERE id=?", *req.FirstName, uid) }
	if req.LastName != nil { h.DB.Exec("UPDATE auth_user SET last_name=? WHERE id=?", *req.LastName, uid) }
	if req.Email != nil { h.DB.Exec("UPDATE auth_user SET email=? WHERE id=?", *req.Email, uid) }
	jsonOK(w, map[string]string{"status": "updated"})
}

func (h *SettingsHandler) BudgetSettings(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var budget models.Budget
	if err := h.DB.Get(&budget, "SELECT * FROM budgets_budget WHERE id=?", id); err != nil {
		jsonError(w, "not found", http.StatusNotFound)
		return
	}
	var members []models.BudgetMembership
	h.DB.Select(&members, "SELECT * FROM budgets_budgetmembership WHERE budget_id=?", id)
	writeJSON(w, http.StatusOK, map[string]interface{}{
		"budget":  budget,
		"members": members,
	})
}

func (h *SettingsHandler) Invite(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var req struct{ Email string `json:"email"` }
	json.NewDecoder(r.Body).Decode(&req)

	var user models.User
	if err := h.DB.Get(&user, "SELECT * FROM auth_user WHERE email=?", req.Email); err != nil {
		jsonError(w, "user not found", http.StatusNotFound)
		return
	}
	mid := strings.ReplaceAll(uuid.New().String(), "-", "")
	h.DB.Exec(
		`INSERT INTO budgets_budgetmembership (id,user_id,budget_id,role,invited_at,accepted_at)
		 VALUES (?,?,?,'editor',datetime('now'),NULL)`,
		mid, user.ID, id,
	)
	jsonOK(w, map[string]string{"status": "invited"})
}

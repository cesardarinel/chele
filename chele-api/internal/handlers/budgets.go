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

type BudgetHandler struct {
	DB *sqlx.DB
}

func NewBudgetHandler(db *sqlx.DB) *BudgetHandler {
	return &BudgetHandler{DB: db}
}

func (h *BudgetHandler) List(w http.ResponseWriter, r *http.Request) {
	uid := r.Context().Value(middleware.UserIDKey).(int)
	var budgets []models.Budget
	h.DB.Select(&budgets,
		`SELECT b.* FROM budgets_budget b
		 JOIN budgets_budgetmembership m ON m.budget_id=b.id
		 WHERE m.user_id=? ORDER BY b.name`, uid)
	jsonOK(w, budgets)
}

func (h *BudgetHandler) Get(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var budget models.Budget
	if err := h.DB.Get(&budget, "SELECT * FROM budgets_budget WHERE id=?", id); err != nil {
		jsonError(w, "not found", http.StatusNotFound)
		return
	}
	jsonOK(w, budget)
}

func (h *BudgetHandler) Create(w http.ResponseWriter, r *http.Request) {
	uid := r.Context().Value(middleware.UserIDKey).(int)
	var req struct {
		Name string `json:"name"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		jsonError(w, "invalid request", http.StatusBadRequest)
		return
	}
	id := strings.ReplaceAll(uuid.New().String(), "-", "")
	h.DB.Exec(
		`INSERT INTO budgets_budget (id,name,description,owner_id,created_at,updated_at)
		 VALUES (?,?,'',?,datetime('now'),datetime('now'))`,
		id, req.Name, uid,
	)
	mid := strings.ReplaceAll(uuid.New().String(), "-", "")
	h.DB.Exec(
		`INSERT INTO budgets_budgetmembership (id,user_id,budget_id,role,invited_at,accepted_at)
		 VALUES (?,?,?,'owner',datetime('now'),datetime('now'))`,
		mid, uid, id,
	)

	// Default groups and categories
	defaultGroups := []struct {
		name    string
		isInc   bool
		cats    []string
	}{
		{"Ingresos", true, []string{"Sueldo", "Freelance", "Inversiones", "Otros ingresos"}},
		{"Gastos Fijos", false, []string{"Alquiler", "Servicios", "Suscripciones", "True Expenses"}},
		{"Gastos Diarios", false, []string{"Comida", "Transporte", "Salidas", "Salud"}},
		{"Ahorro", false, []string{"Fondo de emergencia", "Vacaciones"}},
	}
	for order, g := range defaultGroups {
		gid := strings.ReplaceAll(uuid.New().String(), "-", "")
		h.DB.Exec(
			`INSERT INTO budgets_categorygroup (id,budget_id,name,sort_order,is_income,created_at,updated_at)
			 VALUES (?,?,?,?,?,datetime('now'),datetime('now'))`,
			gid, id, g.name, order, g.isInc,
		)
		for j, c := range g.cats {
			cid := strings.ReplaceAll(uuid.New().String(), "-", "")
			h.DB.Exec(
				`INSERT INTO budgets_category (id,budget_id,group_id,name,sort_order,is_hidden,notes,created_at,updated_at)
				 VALUES (?,?,?,?,?,0,'',datetime('now'),datetime('now'))`,
				cid, id, gid, c, j,
			)
		}
	}

	writeJSON(w, http.StatusCreated, map[string]string{"id": id})
}

func (h *BudgetHandler) Update(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var req struct {
		Name        *string `json:"name"`
		Description *string `json:"description"`
	}
	json.NewDecoder(r.Body).Decode(&req)
	if req.Name != nil {
		h.DB.Exec("UPDATE budgets_budget SET name=? WHERE id=?", *req.Name, id)
	}
	if req.Description != nil {
		h.DB.Exec("UPDATE budgets_budget SET description=? WHERE id=?", *req.Description, id)
	}
	jsonOK(w, map[string]string{"status": "updated"})
}

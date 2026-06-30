package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/cesardarinel/chele-api/internal/models"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type CategoryHandler struct{ DB *sqlx.DB }

func NewCategoryHandler(db *sqlx.DB) *CategoryHandler { return &CategoryHandler{DB: db} }

func (h *CategoryHandler) ListGroups(w http.ResponseWriter, r *http.Request) {
	budgetID := r.URL.Query().Get("budget_id")
	if budgetID == "" {
		jsonError(w, "budget_id required", http.StatusBadRequest)
		return
	}
	var groups []models.CategoryGroup
	h.DB.Select(&groups, "SELECT * FROM budgets_categorygroup WHERE budget_id=? ORDER BY sort_order", budgetID)
	jsonOK(w, groups)
}

func (h *CategoryHandler) CreateGroup(w http.ResponseWriter, r *http.Request) {
	var req struct {
		BudgetID string `json:"budget_id"`
		Name     string `json:"name"`
		IsIncome bool   `json:"is_income"`
	}
	json.NewDecoder(r.Body).Decode(&req)
	id := strings.ReplaceAll(uuid.New().String(), "-", "")
	h.DB.Exec(
		`INSERT INTO budgets_categorygroup (id,budget_id,name,sort_order,is_income,created_at,updated_at)
		 VALUES (?,?,?,0,?,datetime('now'),datetime('now'))`,
		id, req.BudgetID, req.Name, req.IsIncome,
	)
	writeJSON(w, http.StatusCreated, map[string]string{"id": id})
}

func (h *CategoryHandler) UpdateGroup(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var req struct{ Name *string `json:"name"` }
	json.NewDecoder(r.Body).Decode(&req)
	if req.Name != nil {
		h.DB.Exec("UPDATE budgets_categorygroup SET name=? WHERE id=?", *req.Name, id)
	}
	jsonOK(w, map[string]string{"status": "updated"})
}

func (h *CategoryHandler) DeleteGroup(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if _, err := h.DB.Exec("DELETE FROM budgets_category WHERE group_id=?", id); err != nil {
		jsonError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if _, err := h.DB.Exec("DELETE FROM budgets_categorygroup WHERE id=?", id); err != nil {
		jsonError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	jsonOK(w, map[string]string{"status": "deleted"})
}

func (h *CategoryHandler) List(w http.ResponseWriter, r *http.Request) {
	budgetID := r.URL.Query().Get("budget_id")
	if budgetID == "" {
		jsonError(w, "budget_id required", http.StatusBadRequest)
		return
	}
	var cats []models.Category
	h.DB.Select(&cats, "SELECT * FROM budgets_category WHERE budget_id=? AND is_hidden=0 ORDER BY sort_order", budgetID)
	jsonOK(w, cats)
}

func (h *CategoryHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req struct {
		BudgetID string `json:"budget_id"`
		GroupID  string `json:"group_id"`
		Name     string `json:"name"`
	}
	json.NewDecoder(r.Body).Decode(&req)
	id := strings.ReplaceAll(uuid.New().String(), "-", "")
	_, err := h.DB.Exec(
		`INSERT INTO budgets_category (id,budget_id,group_id,name,sort_order,is_hidden,notes,created_at,updated_at)
		 VALUES (?,?,?,?,0,0,'',datetime('now'),datetime('now'))`,
		id, req.BudgetID, req.GroupID, req.Name,
	)
	if err != nil {
		jsonError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusCreated, map[string]string{"id": id})
}

func (h *CategoryHandler) Update(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var req struct{ Name *string `json:"name"` }
	json.NewDecoder(r.Body).Decode(&req)
	if req.Name != nil {
		h.DB.Exec("UPDATE budgets_category SET name=? WHERE id=?", *req.Name, id)
	}
	jsonOK(w, map[string]string{"status": "updated"})
}

func (h *CategoryHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	h.DB.Exec("UPDATE budgets_category SET is_hidden=1 WHERE id=?", id)
	jsonOK(w, map[string]string{"status": "hidden"})
}

type MonthlyBudgetHandler struct{ DB *sqlx.DB }

func NewMonthlyBudgetHandler(db *sqlx.DB) *MonthlyBudgetHandler { return &MonthlyBudgetHandler{DB: db} }

func (h *MonthlyBudgetHandler) List(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	query := "SELECT * FROM budgets_monthlybudget WHERE 1=1"
	var args []interface{}
	if v := q.Get("category_id"); v != "" {
		query += " AND category_id=?"
		args = append(args, v)
	}
	if v := q.Get("month"); v != "" {
		query += " AND month=?"
		args = append(args, v)
	}
	if v := q.Get("year"); v != "" {
		query += " AND year=?"
		args = append(args, v)
	}
	var mbs []models.MonthlyBudget
	h.DB.Select(&mbs, query, args...)
	jsonOK(w, mbs)
}

func (h *MonthlyBudgetHandler) Assign(w http.ResponseWriter, r *http.Request) {
	var req struct {
		CategoryID string `json:"category_id"`
		Month      int    `json:"month"`
		Year       int    `json:"year"`
		Amount     float64 `json:"amount"`
	}
	json.NewDecoder(r.Body).Decode(&req)
	var count int
	h.DB.Get(&count, "SELECT COUNT(*) FROM budgets_monthlybudget WHERE category_id=? AND month=? AND year=?",
		req.CategoryID, req.Month, req.Year)
	if count > 0 {
		h.DB.Exec("UPDATE budgets_monthlybudget SET budgeted=? WHERE category_id=? AND month=? AND year=?",
			req.Amount, req.CategoryID, req.Month, req.Year)
	} else {
		id := strings.ReplaceAll(uuid.New().String(), "-", "")
		h.DB.Exec(
			`INSERT INTO budgets_monthlybudget (id,category_id,month,year,budgeted,created_at,updated_at)
			 VALUES (?,?,?,?,?,datetime('now'),datetime('now'))`,
			id, req.CategoryID, req.Month, req.Year, req.Amount,
		)
	}
	jsonOK(w, map[string]string{"status": "assigned"})
}

func (h *MonthlyBudgetHandler) Move(w http.ResponseWriter, r *http.Request) {
	var req struct {
		FromCategory string  `json:"from_category"`
		ToCategory   string  `json:"to_category"`
		Amount       float64 `json:"amount"`
		Month        int     `json:"month"`
		Year         int     `json:"year"`
	}
	json.NewDecoder(r.Body).Decode(&req)

	h.DB.Exec(`UPDATE budgets_monthlybudget SET budgeted=budgeted-? WHERE category_id=? AND month=? AND year=?`,
		req.Amount, req.FromCategory, req.Month, req.Year)
	h.DB.Exec(`UPDATE budgets_monthlybudget SET budgeted=budgeted+? WHERE category_id=? AND month=? AND year=?`,
		req.Amount, req.ToCategory, req.Month, req.Year)

	jsonOK(w, map[string]string{"status": "moved"})
}

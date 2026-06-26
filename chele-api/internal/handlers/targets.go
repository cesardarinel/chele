package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/cesardarinel/chele-api/internal/models"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type TargetHandler struct{ DB *sqlx.DB }

func NewTargetHandler(db *sqlx.DB) *TargetHandler { return &TargetHandler{DB: db} }

func (h *TargetHandler) List(w http.ResponseWriter, r *http.Request) {
	categoryID := r.URL.Query().Get("category_id")
	budgetID := r.URL.Query().Get("budget_id")
	query := "SELECT g.* FROM goals_goal g JOIN budgets_category c ON c.id=g.category_id WHERE 1=1"
	var args []interface{}
	if categoryID != "" {
		query += " AND g.category_id=?"
		args = append(args, categoryID)
	}
	if budgetID != "" {
		query += " AND c.budget_id=?"
		args = append(args, budgetID)
	}
	query += " ORDER BY g.created_at DESC"
	var targets []models.Goal
	h.DB.Select(&targets, query, args...)
	jsonOK(w, targets)
}

func (h *TargetHandler) Get(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var t models.Goal
	if err := h.DB.Get(&t, "SELECT * FROM goals_goal WHERE id=?", id); err != nil {
		jsonError(w, "not found", http.StatusNotFound)
		return
	}
	jsonOK(w, t)
}

func (h *TargetHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req struct {
		CategoryID  string   `json:"category_id"`
		GoalType    string   `json:"goal_type"`
		Amount      float64  `json:"amount"`
		TargetDate  *string  `json:"target_date"`
		Frequency   *int     `json:"frequency"`
		RefillUpTo  *bool    `json:"refill_up_to"`
		SnoozeMonth *int     `json:"snooze_month"`
		SnoozeYear  *int     `json:"snooze_year"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		jsonError(w, "invalid request", http.StatusBadRequest)
		return
	}
	freq := 12
	if req.Frequency != nil {
		freq = *req.Frequency
	}
	refill := false
	if req.RefillUpTo != nil {
		refill = *req.RefillUpTo
	}
	id := strings.ReplaceAll(uuid.New().String(), "-", "")
	_, err := h.DB.Exec(
		`INSERT INTO goals_goal (id,category_id,goal_type,amount,target_date,frequency,is_completed,refill_up_to,snooze_month,snooze_year,created_at,updated_at)
		 VALUES (?,?,?,?,?,?,0,?,?,?,datetime('now'),datetime('now'))`,
		id, req.CategoryID, req.GoalType, req.Amount, req.TargetDate, freq, refill, req.SnoozeMonth, req.SnoozeYear,
	)
	if err != nil {
		jsonError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusCreated, map[string]string{"id": id})
}

func (h *TargetHandler) Update(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var req struct {
		GoalType    *string  `json:"goal_type"`
		Amount      *float64 `json:"amount"`
		TargetDate  *string  `json:"target_date"`
		Frequency   *int     `json:"frequency"`
		RefillUpTo  *bool    `json:"refill_up_to"`
		SnoozeMonth *int     `json:"snooze_month"`
		SnoozeYear  *int     `json:"snooze_year"`
		IsCompleted *bool    `json:"is_completed"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		jsonError(w, "invalid request", http.StatusBadRequest)
		return
	}
	if req.GoalType != nil {
		h.DB.Exec("UPDATE goals_goal SET goal_type=? WHERE id=?", *req.GoalType, id)
	}
	if req.Amount != nil {
		h.DB.Exec("UPDATE goals_goal SET amount=? WHERE id=?", *req.Amount, id)
	}
	if req.TargetDate != nil {
		h.DB.Exec("UPDATE goals_goal SET target_date=? WHERE id=?", *req.TargetDate, id)
	}
	if req.Frequency != nil {
		h.DB.Exec("UPDATE goals_goal SET frequency=? WHERE id=?", *req.Frequency, id)
	}
	if req.RefillUpTo != nil {
		v := 0
		if *req.RefillUpTo {
			v = 1
		}
		h.DB.Exec("UPDATE goals_goal SET refill_up_to=? WHERE id=?", v, id)
	}
	if req.SnoozeMonth != nil {
		h.DB.Exec("UPDATE goals_goal SET snooze_month=? WHERE id=?", *req.SnoozeMonth, id)
	}
	if req.SnoozeYear != nil {
		h.DB.Exec("UPDATE goals_goal SET snooze_year=? WHERE id=?", *req.SnoozeYear, id)
	}
	if req.IsCompleted != nil {
		v := 0
		if *req.IsCompleted {
			v = 1
		}
		h.DB.Exec("UPDATE goals_goal SET is_completed=? WHERE id=?", v, id)
	}
	jsonOK(w, map[string]string{"status": "updated"})
}

func (h *TargetHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	h.DB.Exec("DELETE FROM goals_goal WHERE id=?", id)
	jsonOK(w, map[string]string{"status": "deleted"})
}

func parseIntOr(s string, def int) int {
	if s == "" {
		return def
	}
	n, err := strconv.Atoi(s)
	if err != nil {
		return def
	}
	return n
}

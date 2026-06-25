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

type PayeeHandler struct{ DB *sqlx.DB }

func NewPayeeHandler(db *sqlx.DB) *PayeeHandler { return &PayeeHandler{DB: db} }

func (h *PayeeHandler) List(w http.ResponseWriter, r *http.Request) {
	budgetID := r.URL.Query().Get("budget_id")
	query := "SELECT * FROM payees_payee WHERE 1=1"
	var args []interface{}
	if budgetID != "" { query += " AND budget_id=?"; args = append(args, budgetID) }
	query += " ORDER BY name"
	var payees []models.Payee
	h.DB.Select(&payees, query, args...)
	jsonOK(w, payees)
}

func (h *PayeeHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req struct {
		BudgetID string `json:"budget_id"`
		Name     string `json:"name"`
	}
	json.NewDecoder(r.Body).Decode(&req)
	id := strings.ReplaceAll(uuid.New().String(), "-", "")
	h.DB.Exec(
		`INSERT INTO payees_payee (id,budget_id,name,created_at,updated_at)
		 VALUES (?,?,?,datetime('now'),datetime('now'))`,
		id, req.BudgetID, req.Name,
	)
	writeJSON(w, http.StatusCreated, map[string]string{"id": id})
}

func (h *PayeeHandler) Update(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var req struct{ Name *string `json:"name"` }
	json.NewDecoder(r.Body).Decode(&req)
	if req.Name != nil { h.DB.Exec("UPDATE payees_payee SET name=? WHERE id=?", *req.Name, id) }
	jsonOK(w, map[string]string{"status": "updated"})
}

func (h *PayeeHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	h.DB.Exec("DELETE FROM payees_payee WHERE id=?", id)
	jsonOK(w, map[string]string{"status": "deleted"})
}

func (h *PayeeHandler) Merge(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var req struct{ TargetID string `json:"target_id"` }
	json.NewDecoder(r.Body).Decode(&req)
	h.DB.Exec("UPDATE transactions_transaction SET payee_id=? WHERE payee_id=?", req.TargetID, id)
	h.DB.Exec("UPDATE schedules_schedule SET payee_id=? WHERE payee_id=?", req.TargetID, id)
	h.DB.Exec("DELETE FROM payees_payee WHERE id=?", id)
	jsonOK(w, map[string]string{"status": "merged"})
}

type GoalHandler struct{ DB *sqlx.DB }

func NewGoalHandler(db *sqlx.DB) *GoalHandler { return &GoalHandler{DB: db} }

func (h *GoalHandler) List(w http.ResponseWriter, r *http.Request) {
	budgetID := r.URL.Query().Get("budget_id")
	var goals []models.Goal
	if budgetID == "" {
		h.DB.Select(&goals, "SELECT * FROM goals_goal ORDER BY created_at DESC")
	} else {
		h.DB.Select(&goals,
			"SELECT g.* FROM goals_goal g JOIN budgets_category c ON c.id=g.category_id WHERE c.budget_id=? ORDER BY g.created_at DESC",
			budgetID)
	}
	jsonOK(w, goals)
}

func (h *GoalHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req struct {
		CategoryID string  `json:"category_id"`
		GoalType   string  `json:"goal_type"`
		Amount     float64 `json:"amount"`
		TargetDate *string `json:"target_date"`
		Frequency  int     `json:"frequency"`
	}
	json.NewDecoder(r.Body).Decode(&req)
	id := strings.ReplaceAll(uuid.New().String(), "-", "")
	h.DB.Exec(
		`INSERT INTO goals_goal (id,category_id,goal_type,amount,target_date,frequency,is_completed,created_at,updated_at)
		 VALUES (?,?,?,?,?,?,0,datetime('now'),datetime('now'))`,
		id, req.CategoryID, req.GoalType, req.Amount, req.TargetDate, req.Frequency,
	)
	// Apply goal to monthly budget
	var mb models.MonthlyBudget
	err := h.DB.Get(&mb, "SELECT * FROM budgets_monthlybudget WHERE category_id=? AND month=CAST(strftime('%m',date('now')) AS INTEGER) AND year=CAST(strftime('%Y',date('now')) AS INTEGER)", req.CategoryID)
	if err != nil {
		mbID := strings.ReplaceAll(uuid.New().String(), "-", "")
		h.DB.Exec(
			`INSERT INTO budgets_monthlybudget (id,category_id,month,year,budgeted,created_at,updated_at)
			 VALUES (?,?,CAST(strftime('%m',date('now')) AS INTEGER),CAST(strftime('%Y',date('now')) AS INTEGER),?,datetime('now'),datetime('now'))`,
			mbID, req.CategoryID, req.Amount,
		)
	} else {
		h.DB.Exec("UPDATE budgets_monthlybudget SET budgeted=budgeted+? WHERE id=?", req.Amount, mb.ID)
	}
	writeJSON(w, http.StatusCreated, map[string]string{"id": id})
}

func (h *GoalHandler) Update(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var req struct {
		Amount      *float64 `json:"amount"`
		GoalType    *string  `json:"goal_type"`
		IsCompleted *bool    `json:"is_completed"`
	}
	json.NewDecoder(r.Body).Decode(&req)
	if req.Amount != nil { h.DB.Exec("UPDATE goals_goal SET amount=? WHERE id=?", *req.Amount, id) }
	if req.GoalType != nil { h.DB.Exec("UPDATE goals_goal SET goal_type=? WHERE id=?", *req.GoalType, id) }
	if req.IsCompleted != nil { h.DB.Exec("UPDATE goals_goal SET is_completed=? WHERE id=?", *req.IsCompleted, id) }
	jsonOK(w, map[string]string{"status": "updated"})
}

func (h *GoalHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	h.DB.Exec("DELETE FROM goals_goal WHERE id=?", id)
	jsonOK(w, map[string]string{"status": "deleted"})
}

type RuleHandler struct{ DB *sqlx.DB }

func NewRuleHandler(db *sqlx.DB) *RuleHandler { return &RuleHandler{DB: db} }

func (h *RuleHandler) List(w http.ResponseWriter, r *http.Request) {
	budgetID := r.URL.Query().Get("budget_id")
	query := "SELECT * FROM rules_rule WHERE 1=1"
	var args []interface{}
	if budgetID != "" { query += " AND budget_id=?"; args = append(args, budgetID) }
	query += " ORDER BY sort_order"
	var rules []models.Rule
	h.DB.Select(&rules, query, args...)
	jsonOK(w, rules)
}

func (h *RuleHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req struct {
		BudgetID          string  `json:"budget_id"`
		Name              string  `json:"name"`
		ConditionField    string  `json:"condition_field"`
		ConditionOperator string  `json:"condition_operator"`
		ConditionValue    string  `json:"condition_value"`
		ActionCategoryID  *string `json:"action_category_id"`
		ActionPayeeID     *string `json:"action_payee_id"`
		ActionNotes       string  `json:"action_notes"`
	}
	json.NewDecoder(r.Body).Decode(&req)
	id := strings.ReplaceAll(uuid.New().String(), "-", "")
	h.DB.Exec(
		`INSERT INTO rules_rule (id,budget_id,name,condition_field,condition_operator,condition_value,action_category_id,action_payee_id,action_notes,sort_order,is_active,created_at,updated_at)
		 VALUES (?,?,?,?,?,?,?,?,?,0,1,datetime('now'),datetime('now'))`,
		id, req.BudgetID, req.Name, req.ConditionField, req.ConditionOperator, req.ConditionValue,
		req.ActionCategoryID, req.ActionPayeeID, req.ActionNotes,
	)
	writeJSON(w, http.StatusCreated, map[string]string{"id": id})
}

func (h *RuleHandler) Update(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var req struct {
		Name              *string `json:"name"`
		ConditionValue    *string `json:"condition_value"`
		ActionCategoryID  *string `json:"action_category_id"`
		ActionPayeeID     *string `json:"action_payee_id"`
		IsActive          *bool   `json:"is_active"`
	}
	json.NewDecoder(r.Body).Decode(&req)
	if req.Name != nil { h.DB.Exec("UPDATE rules_rule SET name=? WHERE id=?", *req.Name, id) }
	if req.ConditionValue != nil { h.DB.Exec("UPDATE rules_rule SET condition_value=? WHERE id=?", *req.ConditionValue, id) }
	if req.ActionCategoryID != nil { h.DB.Exec("UPDATE rules_rule SET action_category_id=? WHERE id=?", *req.ActionCategoryID, id) }
	if req.ActionPayeeID != nil { h.DB.Exec("UPDATE rules_rule SET action_payee_id=? WHERE id=?", *req.ActionPayeeID, id) }
	if req.IsActive != nil { h.DB.Exec("UPDATE rules_rule SET is_active=? WHERE id=?", *req.IsActive, id) }
	jsonOK(w, map[string]string{"status": "updated"})
}

func (h *RuleHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	h.DB.Exec("DELETE FROM rules_rule WHERE id=?", id)
	jsonOK(w, map[string]string{"status": "deleted"})
}

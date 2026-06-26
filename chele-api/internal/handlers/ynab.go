package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/cesardarinel/chele-api/internal/models"
	"github.com/cesardarinel/chele-api/internal/service"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

func newUUID() string {
	return strings.ReplaceAll(uuid.New().String(), "-", "")
}

type YNABHandler struct {
	DB      *sqlx.DB
	Target  *service.TargetService
}

func NewYNABHandler(db *sqlx.DB) *YNABHandler {
	return &YNABHandler{DB: db, Target: service.NewTargetService(db)}
}

// ---- Ready to Assign ----

func (h *YNABHandler) ReadyToAssign(w http.ResponseWriter, r *http.Request) {
	budgetID := chi.URLParam(r, "id")
	month := parseIntOr(r.URL.Query().Get("month"), 0)
	year := parseIntOr(r.URL.Query().Get("year"), 0)
	if month == 0 || year == 0 {
		jsonError(w, "month and year required", http.StatusBadRequest)
		return
	}

	var totalOnBudget float64
	h.DB.Get(&totalOnBudget,
		"SELECT COALESCE(SUM(balance),0) FROM accounts_account WHERE budget_id=? AND on_budget=1", budgetID)

	var totalAssigned float64
	h.DB.Get(&totalAssigned,
		`SELECT COALESCE(SUM(budgeted),0) FROM budgets_monthlybudget
		 WHERE category_id IN (SELECT id FROM budgets_category WHERE budget_id=?) AND month=? AND year=?`,
		budgetID, month, year)

	rta := totalOnBudget - totalAssigned
	if rta < 0 {
		rta = 0
	}

	jsonOK(w, map[string]float64{
		"ready_to_assign": rta,
		"total_on_budget": totalOnBudget,
		"total_assigned":  totalAssigned,
	})
}

// ---- Cover Overspending ----

func (h *YNABHandler) Cover(w http.ResponseWriter, r *http.Request) {
	var req struct {
		FromCategory string  `json:"from_category"`
		ToCategory   string  `json:"to_category"`
		Amount       float64 `json:"amount"`
		Month        int     `json:"month"`
		Year         int     `json:"year"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		jsonError(w, "invalid request", http.StatusBadRequest)
		return
	}

	// Validate source has enough
	var fromAvail float64
	h.DB.Get(&fromAvail,
		`SELECT COALESCE(budgeted,0) - COALESCE((SELECT SUM(amount) FROM transactions_transaction WHERE category_id=? AND CAST(strftime('%m',date) AS INTEGER)=? AND CAST(strftime('%Y',date) AS INTEGER)=?),0) as avail
		 FROM budgets_monthlybudget WHERE category_id=? AND month=? AND year=?`,
		req.FromCategory, req.Month, req.Year, req.FromCategory, req.Month, req.Year)
	if fromAvail < req.Amount {
		jsonError(w, "Fondos insuficientes en la categoría origen", http.StatusBadRequest)
		return
	}

	h.DB.Exec(
		`UPDATE budgets_monthlybudget SET budgeted=budgeted-? WHERE category_id=? AND month=? AND year=?`,
		req.Amount, req.FromCategory, req.Month, req.Year)
	h.DB.Exec(
		`UPDATE budgets_monthlybudget SET budgeted=budgeted+? WHERE category_id=? AND month=? AND year=?`,
		req.Amount, req.ToCategory, req.Month, req.Year)

	jsonOK(w, map[string]string{"status": "covered"})
}

// ---- Spotlight ----

func (h *YNABHandler) Spotlight(w http.ResponseWriter, r *http.Request) {
	budgetID := chi.URLParam(r, "id")
	month := parseIntOr(r.URL.Query().Get("month"), 0)
	year := parseIntOr(r.URL.Query().Get("year"), 0)

	var uncategorized int
	h.DB.Get(&uncategorized,
		"SELECT COUNT(*) FROM transactions_transaction WHERE budget_id=? AND category_id IS NULL", budgetID)

	type overspend struct {
		CategoryID   string  `json:"category_id"`
		CategoryName string  `json:"category_name"`
		Amount       float64 `json:"amount"`
	}
	var overspends []overspend
	rows, _ := h.DB.Queryx(
		`SELECT c.id,c.name,
			COALESCE(mb.budgeted,0) - COALESCE(ABS(t.spent),0) as avail
		 FROM budgets_category c
		 LEFT JOIN budgets_monthlybudget mb ON mb.category_id=c.id AND mb.month=? AND mb.year=?
		 LEFT JOIN (SELECT category_id, SUM(amount) as spent FROM transactions_transaction WHERE budget_id=? AND CAST(strftime('%m',date) AS INTEGER)=? AND CAST(strftime('%Y',date) AS INTEGER)=? GROUP BY category_id) t ON t.category_id=c.id
		 WHERE c.budget_id=? AND c.is_hidden=0
		 HAVING avail < 0`, month, year, budgetID, month, year, budgetID)
	if rows != nil {
		for rows.Next() {
			var o overspend
			rows.Scan(&o.CategoryID, &o.CategoryName, &o.Amount)
			overspends = append(overspends, o)
		}
		rows.Close()
	}

	underfunded := h.Target.ListUnderfunded(budgetID, month, year)

	jsonOK(w, map[string]interface{}{
		"uncategorized_count": uncategorized,
		"overspends":          overspends,
		"underfunded":         underfunded,
		"total_alerts":        len(overspends) + len(underfunded),
	})
}

// ---- Inspector ----

func (h *YNABHandler) Inspector(w http.ResponseWriter, r *http.Request) {
	categoryID := chi.URLParam(r, "id")
	month := parseIntOr(r.URL.Query().Get("month"), 0)
	year := parseIntOr(r.URL.Query().Get("year"), 0)
	if month == 0 || year == 0 {
		jsonError(w, "month and year required", http.StatusBadRequest)
		return
	}

	var catName string
	h.DB.Get(&catName, "SELECT name FROM budgets_category WHERE id=?", categoryID)

	avail := h.Target.GetCategoryAvailable(categoryID, month, year)

	var budgeted float64
	h.DB.Get(&budgeted,
		"SELECT COALESCE(budgeted,0) FROM budgets_monthlybudget WHERE category_id=? AND month=? AND year=?",
		categoryID, month, year)

	var activity float64
	h.DB.Get(&activity,
		"SELECT COALESCE(SUM(amount),0) FROM transactions_transaction WHERE category_id=? AND CAST(strftime('%m',date) AS INTEGER)=? AND CAST(strftime('%Y',date) AS INTEGER)=?",
		categoryID, month, year)

	var avg3m float64
	h.DB.Get(&avg3m,
		`SELECT COALESCE(AVG(ABS(amount)),0) FROM transactions_transaction WHERE category_id=? AND date >= date('now','-3 months')`,
		categoryID)

	var target models.Goal
	targetInfo := map[string]interface{}{}
	err := h.DB.Get(&target, "SELECT * FROM goals_goal WHERE category_id=? AND is_completed=0", categoryID)
	if err == nil {
		underfunded := h.Target.CalculateUnderfunded(
			target.GoalType, target.Amount, target.RefillUpTo,
			target.SnoozeMonth, target.SnoozeYear, categoryID, month, year,
			target.Frequency, target.TargetDate)
		targetInfo = map[string]interface{}{
			"type":           target.GoalType,
			"amount":         target.Amount,
			"refill_up_to":   target.RefillUpTo,
			"underfunded":    underfunded,
			"is_snoozed":     target.SnoozeMonth != nil && target.SnoozeYear != nil &&
				*target.SnoozeMonth == month && *target.SnoozeYear == year,
		}
	}

	isOverspent := avail < 0
	overspentType := ""
	if isOverspent {
		overspentType = "cash"
	}

	jsonOK(w, map[string]interface{}{
		"category_id":      categoryID,
		"category_name":    catName,
		"available":        avail,
		"assigned":         budgeted,
		"activity":         activity,
		"target":           targetInfo,
		"average_spending_3m": avg3m,
		"overspent": map[string]interface{}{
			"is_overspent": isOverspent,
			"amount":       avail,
			"type":         overspentType,
		},
	})
}

// ---- Cost to Be Me ----

func (h *YNABHandler) CostToBeMe(w http.ResponseWriter, r *http.Request) {
	budgetID := chi.URLParam(r, "id")

	var costToBeMe float64
	h.DB.Get(&costToBeMe,
		`SELECT COALESCE(SUM(
			CASE g.goal_type
				WHEN 'monthly' THEN g.amount
				WHEN 'yearly' THEN g.amount/12
				WHEN 'true_expense' THEN g.amount/CASE WHEN g.frequency>0 THEN g.frequency ELSE 12 END
				ELSE g.amount
			END
		),0) FROM goals_goal g
		 JOIN budgets_category c ON c.id=g.category_id
		 WHERE c.budget_id=? AND g.is_completed=0`, budgetID)

	var expectedIncome float64
	h.DB.Get(&expectedIncome,
		`SELECT COALESCE(AVG(monthly),0) FROM (
			SELECT SUM(amount) as monthly FROM transactions_transaction
			WHERE budget_id=? AND amount>0 AND date >= date('now','-3 months')
			GROUP BY CAST(strftime('%Y-%m',date) AS TEXT)
		)`, budgetID)

	diff := expectedIncome - costToBeMe

	jsonOK(w, map[string]interface{}{
		"cost_to_be_me":    costToBeMe,
		"expected_income":  expectedIncome,
		"difference":       diff,
		"is_over_budget":   diff < 0,
	})
}

// ---- Auto-Assign ----

func (h *YNABHandler) AutoAssign(w http.ResponseWriter, r *http.Request) {
	budgetID := chi.URLParam(r, "id")
	var req struct {
		Month int `json:"month"`
		Year  int `json:"year"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		jsonError(w, "invalid request", http.StatusBadRequest)
		return
	}

	// Get ready to assign
	var totalOnBudget float64
	h.DB.Get(&totalOnBudget,
		"SELECT COALESCE(SUM(balance),0) FROM accounts_account WHERE budget_id=? AND on_budget=1", budgetID)
	var totalAssigned float64
	h.DB.Get(&totalAssigned,
		`SELECT COALESCE(SUM(budgeted),0) FROM budgets_monthlybudget
		 WHERE category_id IN (SELECT id FROM budgets_category WHERE budget_id=?) AND month=? AND year=?`,
		budgetID, req.Month, req.Year)
	available := totalOnBudget - totalAssigned
	if available <= 0 {
		jsonOK(w, map[string]interface{}{"assigned_categories": 0, "total_assigned": 0, "remaining": 0})
		return
	}

	underfunded := h.Target.ListUnderfunded(budgetID, req.Month, req.Year)

	assigned := 0
	var totalAssignedAmt float64
	for _, u := range underfunded {
		if available <= 0 {
			break
		}
		assignAmt := u.Deficit
		if assignAmt > available {
			assignAmt = available
		}

		var count int
		h.DB.Get(&count,
			"SELECT COUNT(*) FROM budgets_monthlybudget WHERE category_id=? AND month=? AND year=?",
			u.CategoryID, req.Month, req.Year)
		if count > 0 {
			h.DB.Exec("UPDATE budgets_monthlybudget SET budgeted=budgeted+? WHERE category_id=? AND month=? AND year=?",
				assignAmt, u.CategoryID, req.Month, req.Year)
		} else {
			id := newUUID()
			h.DB.Exec(
				`INSERT INTO budgets_monthlybudget (id,category_id,month,year,budgeted,created_at,updated_at)
				 VALUES (?,?,?,?,?,datetime('now'),datetime('now'))`,
				id, u.CategoryID, req.Month, req.Year, assignAmt)
		}
		assigned++
		totalAssignedAmt += assignAmt
		available -= assignAmt
	}

	jsonOK(w, map[string]interface{}{
		"assigned_categories": assigned,
		"total_assigned":      totalAssignedAmt,
		"remaining":           available,
	})
}

// ---- Rollover ----

func (h *YNABHandler) Rollover(w http.ResponseWriter, r *http.Request) {
	budgetID := chi.URLParam(r, "id")
	fromMonth := parseIntOr(r.URL.Query().Get("from_month"), 0)
	fromYear := parseIntOr(r.URL.Query().Get("from_year"), 0)

	type catRollover struct {
		CategoryID      string  `db:"id" json:"category_id"`
		Name            string  `db:"name" json:"name"`
		RolloverBalance float64 `json:"rollover_balance"`
	}

	var cats []struct {
		ID   string `db:"id"`
		Name string `db:"name"`
	}
	h.DB.Select(&cats, "SELECT id,name FROM budgets_category WHERE budget_id=? AND is_hidden=0 ORDER BY sort_order", budgetID)

	var result []catRollover
	for _, c := range cats {
		avail := h.Target.GetCategoryAvailable(c.ID, fromMonth, fromYear)
		if avail > 0 {
			result = append(result, catRollover{
				CategoryID:      c.ID,
				Name:            c.Name,
				RolloverBalance: avail,
			})
		}
	}

	jsonOK(w, map[string]interface{}{"categories": result})
}

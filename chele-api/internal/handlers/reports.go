package handlers

import (
	"net/http"

	"github.com/jmoiron/sqlx"
)

type ReportHandler struct{ DB *sqlx.DB }

func NewReportHandler(db *sqlx.DB) *ReportHandler { return &ReportHandler{DB: db} }

func (h *ReportHandler) NetWorth(w http.ResponseWriter, r *http.Request) {
	budgetID := r.URL.Query().Get("budget_id")
	if budgetID == "" { jsonError(w, "budget_id required", http.StatusBadRequest); return }

	var assets float64
	h.DB.Get(&assets, "SELECT COALESCE(SUM(balance),0) FROM accounts_account WHERE budget_id=?", budgetID)
	var liabilities float64
	h.DB.Get(&liabilities, "SELECT COALESCE(SUM(balance),0) FROM credit_cards_creditcard WHERE budget_id=?", budgetID)
	if liabilities < 0 {
		liabilities = -liabilities
	}

	writeJSON(w, http.StatusOK, map[string]float64{
		"assets":     assets,
		"liabilities": liabilities,
		"net_worth":  assets - liabilities,
	})
}

func (h *ReportHandler) CashFlow(w http.ResponseWriter, r *http.Request) {
	budgetID := r.URL.Query().Get("budget_id")
	if budgetID == "" { jsonError(w, "budget_id required", http.StatusBadRequest); return }

	type cfMonth struct {
		Month   string  `json:"month"`
		Income  float64 `json:"income"`
		Expenses float64 `json:"expenses"`
		Net     float64 `json:"net"`
	}
	var months []cfMonth
	// Last 6 months
	for i := 5; i >= 0; i-- {
		var income, expenses float64
		h.DB.Get(&income,
			`SELECT COALESCE(SUM(amount),0) FROM transactions_transaction
			 WHERE budget_id=? AND amount>0 AND date>=date('now','-? months','start of month') AND date<date('now','-? months','start of month','+1 month')`,
			budgetID, i, i)
		h.DB.Get(&expenses,
			`SELECT COALESCE(SUM(amount),0) FROM transactions_transaction
			 WHERE budget_id=? AND amount<0 AND date>=date('now','-? months','start of month') AND date<date('now','-? months','start of month','+1 month')`,
			budgetID, i, i)
		months = append(months, cfMonth{
			Month:    "", // simplified
			Income:   income,
			Expenses: -expenses,
			Net:      income + expenses,
		})
	}
	jsonOK(w, months)
}

func (h *ReportHandler) BudgetVsReality(w http.ResponseWriter, r *http.Request) {
	budgetID := r.URL.Query().Get("budget_id")
	if budgetID == "" { jsonError(w, "budget_id required", http.StatusBadRequest); return }

	type catItem struct {
		Name       string  `json:"name"`
		Budgeted   float64 `json:"budgeted"`
		Spent      float64 `json:"spent"`
		Difference float64 `json:"difference"`
	}
	type grpItem struct {
		Group      string    `json:"group"`
		Categories []catItem `json:"categories"`
	}

	var groups []struct {
		ID   string `db:"id"`
		Name string `db:"name"`
	}
	h.DB.Select(&groups, "SELECT id,name FROM budgets_categorygroup WHERE budget_id=? ORDER BY sort_order", budgetID)

	var result []grpItem
	for _, g := range groups {
		var cats []struct {
			ID   string `db:"id"`
			Name string `db:"name"`
		}
		h.DB.Select(&cats, "SELECT id,name FROM budgets_category WHERE group_id=? AND is_hidden=0", g.ID)
		var items []catItem
		for _, c := range cats {
			var budgeted float64
			h.DB.Get(&budgeted,
				`SELECT COALESCE(budgeted,0) FROM budgets_monthlybudget
				 WHERE category_id=? AND month=CAST(strftime('%m',date('now')) AS INTEGER) AND year=CAST(strftime('%Y',date('now')) AS INTEGER)`, c.ID)
			var spent float64
			h.DB.Get(&spent,
				`SELECT COALESCE(SUM(amount),0) FROM transactions_transaction
				 WHERE category_id=? AND strftime('%m',date)=strftime('%m',date('now')) AND strftime('%Y',date)=strftime('%Y',date('now'))`, c.ID)
			if spent < 0 { spent = -spent }
			items = append(items, catItem{
				Name:       c.Name,
				Budgeted:   budgeted,
				Spent:      spent,
				Difference: budgeted - spent,
			})
		}
		result = append(result, grpItem{Group: g.Name, Categories: items})
	}
	jsonOK(w, result)
}

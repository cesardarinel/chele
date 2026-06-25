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

type AccountHandler struct {
	DB      *sqlx.DB
	Balance *service.BalanceService
}

func NewAccountHandler(db *sqlx.DB) *AccountHandler {
	return &AccountHandler{DB: db, Balance: service.NewBalanceService(db)}
}

func (h *AccountHandler) List(w http.ResponseWriter, r *http.Request) {
	budgetID := r.URL.Query().Get("budget_id")
	if budgetID == "" {
		jsonError(w, "budget_id required", http.StatusBadRequest)
		return
	}
	var accounts []models.Account
	h.DB.Select(&accounts, "SELECT * FROM accounts_account WHERE budget_id=? ORDER BY name", budgetID)

	var onB, offB []models.Account
	var tOn, tOff float64
	for _, a := range accounts {
		if a.OnBudget {
			onB = append(onB, a)
			tOn += a.Balance
		} else {
			offB = append(offB, a)
			tOff += a.Balance
		}
	}

	var ccs []models.CreditCard
	h.DB.Select(&ccs, "SELECT * FROM credit_cards_creditcard WHERE budget_id=?", budgetID)
	var ccDebt []models.CreditCard
	var tCC float64
	for _, cc := range ccs {
		if cc.Balance < 0 {
			ccDebt = append(ccDebt, cc)
			tCC += cc.Balance
		}
	}

	var loans []models.Loan
	h.DB.Select(&loans, "SELECT * FROM loans_loan WHERE budget_id=? AND status='active'", budgetID)
	var tLoan float64
	for _, l := range loans {
		tLoan += l.RemainingBalance
	}
	tLoan = -tLoan

	jsonOK(w, models.AccountSummary{
		OnBudgetAccounts:  onB,
		OffBudgetAccounts: offB,
		TotalOnBudget:     tOn,
		TotalOffBudget:    tOff,
		GrandTotal:        tOn + tOff,
		CCDebtCards:       ccDebt,
		TotalCCDebt:       tCC,
		Loans:             loans,
		TotalLoanDebt:     tLoan,
		TotalDebt:         tCC + tLoan,
	})
}

func (h *AccountHandler) Get(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var account models.Account
	if err := h.DB.Get(&account, "SELECT * FROM accounts_account WHERE id=?", id); err != nil {
		jsonError(w, "not found", http.StatusNotFound)
		return
	}
	var txn []models.Transaction
	h.DB.Select(&txn, "SELECT * FROM transactions_transaction WHERE account_id=? ORDER BY date DESC", id)
	var sched []models.Schedule
	h.DB.Select(&sched, "SELECT * FROM schedules_schedule WHERE account_id=? AND is_active=1 ORDER BY next_date", id)
	writeJSON(w, http.StatusOK, map[string]interface{}{
		"account":      account,
		"transactions": txn,
		"schedules":    sched,
	})
}

func (h *AccountHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req struct {
		BudgetID  string  `json:"budget_id"`
		Name      string  `json:"name"`
		OnBudget  bool    `json:"on_budget"`
		Balance   float64 `json:"balance"`
		StartDate string  `json:"start_date"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		jsonError(w, "invalid request", http.StatusBadRequest)
		return
	}
	id := strings.ReplaceAll(uuid.New().String(), "-", "")
	_, err := h.DB.Exec(
		`INSERT INTO accounts_account (id,budget_id,name,on_budget,balance,notes,created_at,updated_at)
		 VALUES (?,?,?,?,0,'',datetime('now'),datetime('now'))`,
		id, req.BudgetID, req.Name, req.OnBudget,
	)
	if err != nil {
		jsonError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusCreated, map[string]string{"id": id})
}

func (h *AccountHandler) Update(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var req struct {
		Name     *string `json:"name"`
		OnBudget *bool   `json:"on_budget"`
		Notes    *string `json:"notes"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		jsonError(w, "invalid request", http.StatusBadRequest)
		return
	}
	if req.Name != nil {
		h.DB.Exec("UPDATE accounts_account SET name=? WHERE id=?", *req.Name, id)
	}
	if req.OnBudget != nil {
		h.DB.Exec("UPDATE accounts_account SET on_budget=? WHERE id=?", *req.OnBudget, id)
	}
	if req.Notes != nil {
		h.DB.Exec("UPDATE accounts_account SET notes=? WHERE id=?", *req.Notes, id)
	}
	jsonOK(w, map[string]string{"status": "updated"})
}

func (h *AccountHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	h.DB.Exec("DELETE FROM accounts_account WHERE id=?", id)
	jsonOK(w, map[string]string{"status": "deleted"})
}

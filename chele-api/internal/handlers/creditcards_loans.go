package handlers

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/cesardarinel/chele-api/internal/models"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type CreditCardHandler struct{ DB *sqlx.DB }

func NewCreditCardHandler(db *sqlx.DB) *CreditCardHandler { return &CreditCardHandler{DB: db} }

func (h *CreditCardHandler) List(w http.ResponseWriter, r *http.Request) {
	budgetID := r.URL.Query().Get("budget_id")
	query := "SELECT * FROM credit_cards_creditcard WHERE 1=1"
	var args []interface{}
	if budgetID != "" {
		query += " AND budget_id=?"
		args = append(args, budgetID)
	}
	query += " ORDER BY name"
	var cards []models.CreditCard
	h.DB.Select(&cards, query, args...)
	jsonOK(w, cards)
}

func (h *CreditCardHandler) Get(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var card models.CreditCard
	if err := h.DB.Get(&card, "SELECT * FROM credit_cards_creditcard WHERE id=?", id); err != nil {
		jsonError(w, "not found", http.StatusNotFound)
		return
	}
	jsonOK(w, card)
}

func (h *CreditCardHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req struct {
		BudgetID    string  `json:"budget_id"`
		Name        string  `json:"name"`
		Limit       float64 `json:"limit"`
		Balance     float64 `json:"balance"`
		InterestRate float64 `json:"interest_rate"`
		ClosingDay  int     `json:"closing_day"`
		DueDay      int     `json:"due_day"`
		Notes       string  `json:"notes"`
	}
	json.NewDecoder(r.Body).Decode(&req)
	id := strings.ReplaceAll(uuid.New().String(), "-", "")
	h.DB.Exec(
		`INSERT INTO credit_cards_creditcard (id,budget_id,name,"limit",balance,interest_rate,closing_day,due_day,notes,created_at,updated_at)
		 VALUES (?,?,?,?,?,?,?,?,?,datetime('now'),datetime('now'))`,
		id, req.BudgetID, req.Name, req.Limit, req.Balance, req.InterestRate, req.ClosingDay, req.DueDay, req.Notes,
	)
	writeJSON(w, http.StatusCreated, map[string]string{"id": id})
}

func (h *CreditCardHandler) Update(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var req struct {
		Name        *string  `json:"name"`
		Limit       *float64 `json:"limit"`
		InterestRate *float64 `json:"interest_rate"`
		ClosingDay  *int     `json:"closing_day"`
		DueDay      *int     `json:"due_day"`
		Notes       *string  `json:"notes"`
	}
	json.NewDecoder(r.Body).Decode(&req)
	if req.Name != nil { h.DB.Exec(`UPDATE credit_cards_creditcard SET name=? WHERE id=?`, *req.Name, id) }
	if req.Limit != nil { h.DB.Exec(`UPDATE credit_cards_creditcard SET "limit"=? WHERE id=?`, *req.Limit, id) }
	if req.InterestRate != nil { h.DB.Exec(`UPDATE credit_cards_creditcard SET interest_rate=? WHERE id=?`, *req.InterestRate, id) }
	if req.ClosingDay != nil { h.DB.Exec(`UPDATE credit_cards_creditcard SET closing_day=? WHERE id=?`, *req.ClosingDay, id) }
	if req.DueDay != nil { h.DB.Exec(`UPDATE credit_cards_creditcard SET due_day=? WHERE id=?`, *req.DueDay, id) }
	if req.Notes != nil { h.DB.Exec(`UPDATE credit_cards_creditcard SET notes=? WHERE id=?`, *req.Notes, id) }
	jsonOK(w, map[string]string{"status": "updated"})
}

func (h *CreditCardHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	h.DB.Exec("DELETE FROM credit_cards_creditcard WHERE id=?", id)
	jsonOK(w, map[string]string{"status": "deleted"})
}

func (h *CreditCardHandler) Pay(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var req struct {
		AccountID string  `json:"account_id"`
		Amount    float64 `json:"amount"`
	}
	json.NewDecoder(r.Body).Decode(&req)

	var card models.CreditCard
	if err := h.DB.Get(&card, "SELECT * FROM credit_cards_creditcard WHERE id=?", id); err != nil {
		jsonError(w, "not found", http.StatusNotFound)
		return
	}

	// Apply payment
	h.DB.Exec("UPDATE credit_cards_creditcard SET balance=balance+? WHERE id=?", req.Amount, id)
	h.DB.Exec("UPDATE accounts_account SET balance=balance-? WHERE id=?", req.Amount, req.AccountID)

	txnID := strings.ReplaceAll(uuid.New().String(), "-", "")
	h.DB.Exec(
		`INSERT INTO transactions_transaction (id,budget_id,account_id,date,amount,notes,created_at,updated_at)
		 VALUES (?,?,?,date('now'),?,?,datetime('now'),datetime('now'))`,
		txnID, card.BudgetID, req.AccountID, -req.Amount, "Pago TC: "+card.Name,
	)

	jsonOK(w, map[string]string{"status": "paid"})
}

type LoanHandler struct{ DB *sqlx.DB }

func NewLoanHandler(db *sqlx.DB) *LoanHandler { return &LoanHandler{DB: db} }

func (h *LoanHandler) List(w http.ResponseWriter, r *http.Request) {
	budgetID := r.URL.Query().Get("budget_id")
	query := "SELECT * FROM loans_loan WHERE 1=1"
	var args []interface{}
	if budgetID != "" { query += " AND budget_id=?"; args = append(args, budgetID) }
	query += " ORDER BY name"
	var loans []models.Loan
	h.DB.Select(&loans, query, args...)
	jsonOK(w, loans)
}

func (h *LoanHandler) Get(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var loan models.Loan
	if err := h.DB.Get(&loan, "SELECT * FROM loans_loan WHERE id=?", id); err != nil {
		jsonError(w, "not found", http.StatusNotFound)
		return
	}
	var installments []models.Installment
	h.DB.Select(&installments, "SELECT * FROM loans_installment WHERE loan_id=? ORDER BY number", id)
	writeJSON(w, http.StatusOK, map[string]interface{}{
		"loan":         loan,
		"installments": installments,
	})
}

func (h *LoanHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req struct {
		BudgetID         string  `json:"budget_id"`
		AccountID        *string `json:"account_id"`
		Type             string  `json:"type"`
		Name             string  `json:"name"`
		TotalAmount      float64 `json:"total_amount"`
		InterestRate     float64 `json:"interest_rate"`
		TotalInstallments int    `json:"total_installments"`
		StartDate        string  `json:"start_date"`
		NextDueDate      string  `json:"next_due_date"`
		InstallmentAmount float64 `json:"installment_amount"`
		Notes            string  `json:"notes"`
	}
	json.NewDecoder(r.Body).Decode(&req)
	id := strings.ReplaceAll(uuid.New().String(), "-", "")
	h.DB.Exec(
		`INSERT INTO loans_loan (id,budget_id,account_id,type,name,status,total_amount,interest_rate,remaining_balance,total_installments,paid_installments,start_date,next_due_date,installment_amount,notes,created_at,updated_at)
		 VALUES (?,?,?,?,?,'active',?,?,?,?,0,?,?,?,?,datetime('now'),datetime('now'))`,
		id, req.BudgetID, req.AccountID, req.Type, req.Name, req.TotalAmount, req.InterestRate,
		req.TotalAmount, req.TotalInstallments, req.StartDate, req.NextDueDate, req.InstallmentAmount, req.Notes,
	)
	// Generate installments
	dDate := req.NextDueDate
	for i := 1; i <= req.TotalInstallments; i++ {
		instID := strings.ReplaceAll(uuid.New().String(), "-", "")
		h.DB.Exec(
			`INSERT INTO loans_installment (id,loan_id,number,amount,due_date,paid,notes)
			 VALUES (?,?,?,?,?,0,'')`,
			instID, id, i, req.InstallmentAmount, dDate,
		)
		// Advance one month
		dDate = advanceMonth(dDate)
	}
	writeJSON(w, http.StatusCreated, map[string]string{"id": id})
}

func advanceMonth(d string) string {
	t, _ := time.Parse("2006-01-02", d)
	t = t.AddDate(0, 1, 0)
	return t.Format("2006-01-02")
}

func (h *LoanHandler) PayInstallment(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "loan_id")
	var req struct{ AccountID string `json:"account_id"` }
	json.NewDecoder(r.Body).Decode(&req)

	var inst models.Installment
	err := h.DB.Get(&inst, "SELECT * FROM loans_installment WHERE loan_id=? AND paid=0 ORDER BY number LIMIT 1", id)
	if err != nil {
		jsonError(w, "no pending installments", http.StatusBadRequest)
		return
	}

	var loan models.Loan
	h.DB.Get(&loan, "SELECT * FROM loans_loan WHERE id=?", id)

	// Create transaction
	txnID := strings.ReplaceAll(uuid.New().String(), "-", "")
	h.DB.Exec(
		`INSERT INTO transactions_transaction (id,budget_id,account_id,date,amount,notes,created_at,updated_at)
		 VALUES (?,?,?,date('now'),?,?,datetime('now'),datetime('now'))`,
		txnID, loan.BudgetID, req.AccountID, -inst.Amount, "Cuota: "+loan.Name,
	)
	h.DB.Exec("UPDATE accounts_account SET balance=balance-? WHERE id=?", inst.Amount, req.AccountID)
	h.DB.Exec("UPDATE loans_installment SET paid=1,paid_date=date('now') WHERE id=?", inst.ID)
	h.DB.Exec("UPDATE loans_loan SET paid_installments=paid_installments+1,remaining_balance=remaining_balance-? WHERE id=?", inst.Amount, id)

	jsonOK(w, map[string]string{"status": "paid"})
}

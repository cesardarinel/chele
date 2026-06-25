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

type TransactionHandler struct {
	DB      *sqlx.DB
	Balance *service.BalanceService
}

func NewTransactionHandler(db *sqlx.DB) *TransactionHandler {
	return &TransactionHandler{DB: db, Balance: service.NewBalanceService(db)}
}

func (h *TransactionHandler) List(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	query := "SELECT t.* FROM transactions_transaction t WHERE 1=1"
	var args []interface{}
	if v := q.Get("budget_id"); v != "" {
		query += " AND t.budget_id=?"
		args = append(args, v)
	}
	if v := q.Get("account_id"); v != "" {
		query += " AND t.account_id=?"
		args = append(args, v)
	}
	if v := q.Get("month"); v != "" {
		query += " AND CAST(strftime('%m',t.date) AS INTEGER)=?"
		args = append(args, v)
	}
	if v := q.Get("year"); v != "" {
		query += " AND CAST(strftime('%Y',t.date) AS INTEGER)=?"
		args = append(args, v)
	}
	if v := q.Get("category_id"); v != "" {
		query += " AND t.category_id=?"
		args = append(args, v)
	}
	query += " ORDER BY t.date DESC"
	var txns []models.Transaction
	h.DB.Select(&txns, query, args...)
	jsonOK(w, txns)
}

func (h *TransactionHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req struct {
		BudgetID   string  `json:"budget_id"`
		AccountID  string  `json:"account_id"`
		Date       string  `json:"date"`
		Amount     float64 `json:"amount"`
		Direction  string  `json:"direction"`
		PayeeID    *string `json:"payee_id"`
		CategoryID *string `json:"category_id"`
		Notes      string  `json:"notes"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		jsonError(w, "invalid request", http.StatusBadRequest)
		return
	}
	if req.Direction == "expense" && req.Amount > 0 {
		req.Amount = -req.Amount
	}
	id := strings.ReplaceAll(uuid.New().String(), "-", "")
	_, err := h.DB.Exec(
		`INSERT INTO transactions_transaction
		 (id,budget_id,account_id,date,amount,payee_id,category_id,notes,reconciled,cleared,created_at,updated_at)
		 VALUES (?,?,?,?,?,?,?,?,0,0,datetime('now'),datetime('now'))`,
		id, req.BudgetID, req.AccountID, req.Date, req.Amount, req.PayeeID, req.CategoryID, req.Notes,
	)
	if err != nil {
		jsonError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	h.Balance.Apply(req.AccountID, req.Amount)
	writeJSON(w, http.StatusCreated, map[string]string{"id": id})
}

func (h *TransactionHandler) Update(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var old models.Transaction
	if err := h.DB.Get(&old, "SELECT * FROM transactions_transaction WHERE id=?", id); err != nil {
		jsonError(w, "not found", http.StatusNotFound)
		return
	}

	var req struct {
		AccountID  *string  `json:"account_id"`
		Date       *string  `json:"date"`
		Amount     *float64 `json:"amount"`
		Direction  *string  `json:"direction"`
		PayeeID    *string  `json:"payee_id"`
		CategoryID *string  `json:"category_id"`
		Notes      *string  `json:"notes"`
	}
	json.NewDecoder(r.Body).Decode(&req)

	newAmt := old.Amount
	if req.Amount != nil {
		newAmt = *req.Amount
	}
	if req.Direction != nil && *req.Direction == "expense" && newAmt > 0 {
		newAmt = -newAmt
	}

	if old.AccountID != nil {
		h.Balance.Reverse(*old.AccountID, old.Amount)
	}
	h.DB.Exec(
		`UPDATE transactions_transaction SET date=?,amount=?,payee_id=?,category_id=?,notes=?,updated_at=datetime('now')
		 WHERE id=?`,
		req.Date, newAmt, req.PayeeID, req.CategoryID, req.Notes, id,
	)
	if req.AccountID != nil {
		h.Balance.Apply(*req.AccountID, newAmt)
	} else if old.AccountID != nil {
		h.Balance.Apply(*old.AccountID, newAmt)
	}
	jsonOK(w, map[string]string{"status": "updated"})
}

func (h *TransactionHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var txn models.Transaction
	if err := h.DB.Get(&txn, "SELECT * FROM transactions_transaction WHERE id=?", id); err != nil {
		jsonError(w, "not found", http.StatusNotFound)
		return
	}
	if txn.AccountID != nil {
		h.Balance.Reverse(*txn.AccountID, txn.Amount)
	}
	h.DB.Exec("DELETE FROM transactions_transaction WHERE id=?", id)
	jsonOK(w, map[string]string{"status": "deleted"})
}

func (h *TransactionHandler) Bulk(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Action     string   `json:"action"`
		IDs        []string `json:"ids"`
		CategoryID *string  `json:"category_id"`
	}
	json.NewDecoder(r.Body).Decode(&req)
	if req.Action == "delete" {
		for _, id := range req.IDs {
			var txn models.Transaction
			if h.DB.Get(&txn, "SELECT * FROM transactions_transaction WHERE id=?", id) == nil {
				if txn.AccountID != nil {
					h.Balance.Reverse(*txn.AccountID, txn.Amount)
				}
				h.DB.Exec("DELETE FROM transactions_transaction WHERE id=?", id)
			}
		}
	} else if req.Action == "categorize" && req.CategoryID != nil {
		for _, id := range req.IDs {
			h.DB.Exec("UPDATE transactions_transaction SET category_id=? WHERE id=?", *req.CategoryID, id)
		}
	}
	jsonOK(w, map[string]string{"status": "ok"})
}

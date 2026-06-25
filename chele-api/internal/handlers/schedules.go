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

type ScheduleHandler struct {
	DB       *sqlx.DB
	Schedules *service.ScheduleService
}

func NewScheduleHandler(db *sqlx.DB) *ScheduleHandler {
	return &ScheduleHandler{DB: db, Schedules: service.NewScheduleService(db)}
}

func (h *ScheduleHandler) List(w http.ResponseWriter, r *http.Request) {
	budgetID := r.URL.Query().Get("budget_id")
	query := "SELECT * FROM schedules_schedule WHERE 1=1"
	var args []interface{}
	if budgetID != "" {
		query += " AND budget_id=?"
		args = append(args, budgetID)
	}
	query += " ORDER BY next_date"
	var scheds []models.Schedule
	h.DB.Select(&scheds, query, args...)
	jsonOK(w, scheds)
}

func (h *ScheduleHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req struct {
		BudgetID     string  `json:"budget_id"`
		AccountID    string  `json:"account_id"`
		Amount       float64 `json:"amount"`
		Direction    string  `json:"direction"`
		Frequency    string  `json:"frequency"`
		NextDate     string  `json:"next_date"`
		PayeeID      *string `json:"payee_id"`
		CategoryID   *string `json:"category_id"`
		Notes        string  `json:"notes"`
		SkipWeekends bool    `json:"skip_weekends"`
	}
	json.NewDecoder(r.Body).Decode(&req)
	id := strings.ReplaceAll(uuid.New().String(), "-", "")
	h.DB.Exec(
		`INSERT INTO schedules_schedule
		 (id,budget_id,account_id,amount,direction,frequency,next_date,payee_id,category_id,notes,skip_weekends,is_active,created_at,updated_at)
		 VALUES (?,?,?,?,?,?,?,?,?,?,?,1,datetime('now'),datetime('now'))`,
		id, req.BudgetID, req.AccountID, req.Amount, req.Direction, req.Frequency, req.NextDate,
		req.PayeeID, req.CategoryID, req.Notes, req.SkipWeekends,
	)
	writeJSON(w, http.StatusCreated, map[string]string{"id": id})
}

func (h *ScheduleHandler) Update(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var req struct {
		Amount       *float64 `json:"amount"`
		Frequency    *string  `json:"frequency"`
		NextDate     *string  `json:"next_date"`
		Direction    *string  `json:"direction"`
		IsActive     *bool    `json:"is_active"`
		SkipWeekends *bool    `json:"skip_weekends"`
		PayeeID      *string  `json:"payee_id"`
		CategoryID   *string  `json:"category_id"`
		Notes        *string  `json:"notes"`
	}
	json.NewDecoder(r.Body).Decode(&req)
	if req.Amount != nil { h.DB.Exec("UPDATE schedules_schedule SET amount=? WHERE id=?", *req.Amount, id) }
	if req.Frequency != nil { h.DB.Exec("UPDATE schedules_schedule SET frequency=? WHERE id=?", *req.Frequency, id) }
	if req.NextDate != nil { h.DB.Exec("UPDATE schedules_schedule SET next_date=? WHERE id=?", *req.NextDate, id) }
	if req.Direction != nil { h.DB.Exec("UPDATE schedules_schedule SET direction=? WHERE id=?", *req.Direction, id) }
	if req.IsActive != nil { h.DB.Exec("UPDATE schedules_schedule SET is_active=? WHERE id=?", *req.IsActive, id) }
	if req.SkipWeekends != nil { h.DB.Exec("UPDATE schedules_schedule SET skip_weekends=? WHERE id=?", *req.SkipWeekends, id) }
	if req.Notes != nil { h.DB.Exec("UPDATE schedules_schedule SET notes=? WHERE id=?", *req.Notes, id) }
	jsonOK(w, map[string]string{"status": "updated"})
}

func (h *ScheduleHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	h.DB.Exec("DELETE FROM schedules_schedule WHERE id=?", id)
	jsonOK(w, map[string]string{"status": "deleted"})
}

func (h *ScheduleHandler) Process(w http.ResponseWriter, r *http.Request) {
	budgetID := r.URL.Query().Get("budget_id")
	if budgetID == "" {
		jsonError(w, "budget_id required", http.StatusBadRequest)
		return
	}
	count, err := h.Schedules.ProcessDue(budgetID)
	if err != nil {
		jsonError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	jsonOK(w, map[string]int{"processed": count})
}

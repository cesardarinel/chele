package service

import (
	"strings"
	"testing"
	"time"

	"github.com/cesardarinel/chele-api/internal/middleware"
	"github.com/cesardarinel/chele-api/internal/testutil"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

func insertBudget(db *sqlx.DB) string {
	hash := middleware.MakeDjangoPassword("pass")
	db.MustExec("INSERT INTO auth_user (id,password,username,email,is_superuser,is_staff,is_active,first_name,last_name,date_joined) VALUES (1,?,'t@t.com','t@t.com',0,0,1,'T','U',datetime('now'))", hash)

	id := strings.ReplaceAll(uuid.New().String(), "-", "")
	db.MustExec("INSERT INTO budgets_budget (id,name,description,owner_id,created_at,updated_at) VALUES (?,'Test','',1,datetime('now'),datetime('now'))", id)
	return id
}

func insertAccount(db *sqlx.DB, budgetID string) string {
	id := strings.ReplaceAll(uuid.New().String(), "-", "")
	db.MustExec("INSERT INTO accounts_account (id,budget_id,name,on_budget,balance,notes,created_at,updated_at) VALUES (?,?,'Test',1,1000,'',datetime('now'),datetime('now'))", id, budgetID)
	return id
}

func TestProcessDue_CreatesTransaction(t *testing.T) {
	db := testutil.SetupDB()
	defer db.Close()
	svc := NewScheduleService(db)

	budgetID := insertBudget(db)
	acctID := insertAccount(db, budgetID)

	schID := strings.ReplaceAll(uuid.New().String(), "-", "")
	db.MustExec("INSERT INTO schedules_schedule (id,budget_id,account_id,amount,frequency,next_date,is_active,direction,notes,skip_weekends,apply_before_weekend,created_at,updated_at) VALUES (?,?,?,200,'monthly',date('now'),1,'expense','',0,0,datetime('now'),datetime('now'))", schID, budgetID, acctID)

	count, err := svc.ProcessDue(budgetID)
	if err != nil {
		t.Fatal(err)
	}
	if count != 1 {
		t.Errorf("expected 1 schedule processed, got %d", count)
	}

	var txnCount int
	db.Get(&txnCount, "SELECT COUNT(*) FROM transactions_transaction WHERE account_id=?", acctID)
	if txnCount != 1 {
		t.Errorf("expected 1 transaction created, got %d", txnCount)
	}
}

func TestProcessDue_OnlyDue(t *testing.T) {
	db := testutil.SetupDB()
	defer db.Close()
	svc := NewScheduleService(db)

	budgetID := insertBudget(db)
	acctID := insertAccount(db, budgetID)

	futureDate := time.Now().AddDate(0, 0, 30).Format("2006-01-02")
	schID := strings.ReplaceAll(uuid.New().String(), "-", "")
	db.MustExec("INSERT INTO schedules_schedule (id,budget_id,account_id,amount,frequency,next_date,is_active,direction,notes,skip_weekends,apply_before_weekend,created_at,updated_at) VALUES (?,?,?,200,'monthly',?,1,'expense','',0,0,datetime('now'),datetime('now'))", schID, budgetID, acctID, futureDate)

	count, err := svc.ProcessDue(budgetID)
	if err != nil {
		t.Fatal(err)
	}
	if count != 0 {
		t.Errorf("expected 0, got %d", count)
	}
}

func TestAdvanceDate(t *testing.T) {
	tests := []struct {
		name    string
		current string
		freq    string
		want    string
	}{
		{"weekly", "2026-01-01", "weekly", "2026-01-08"},
		{"biweekly", "2026-01-01", "biweekly", "2026-01-15"},
		{"monthly", "2026-01-01", "monthly", "2026-02-01"},
		{"quarterly", "2026-01-01", "quarterly", "2026-04-01"},
		{"yearly", "2026-01-01", "yearly", "2027-01-01"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := AdvanceDate(tt.current, tt.freq, false, false)
			if got != tt.want {
				t.Errorf("got %s, want %s", got, tt.want)
			}
		})
	}
}

package service

import (
	"strings"
	"testing"

	"github.com/cesardarinel/chele-api/internal/middleware"
	"github.com/cesardarinel/chele-api/internal/testutil"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

func insertUserBudgetAccount(db *sqlx.DB) (string, string) {
	hash := middleware.MakeDjangoPassword("pass")
	db.MustExec("INSERT INTO auth_user (id,password,username,email,is_superuser,is_staff,is_active,first_name,last_name,date_joined) VALUES (1,?, 'test@t.com','test@t.com',0,0,1,'T','U',datetime('now'))", hash)

	budgetID := strings.ReplaceAll(uuid.New().String(), "-", "")
	db.MustExec("INSERT INTO budgets_budget (id,name,description,owner_id,created_at,updated_at) VALUES (?,'Test','',1,datetime('now'),datetime('now'))", budgetID)

	acctID := strings.ReplaceAll(uuid.New().String(), "-", "")
	db.MustExec("INSERT INTO accounts_account (id,budget_id,name,on_budget,balance,notes,created_at,updated_at) VALUES (?,?,'Test',1,500,'',datetime('now'),datetime('now'))", acctID, budgetID)

	return budgetID, acctID
}

func TestBalanceApply(t *testing.T) {
	db := testutil.SetupDB()
	defer db.Close()
	svc := NewBalanceService(db)

	_, acctID := insertUserBudgetAccount(db)

	err := svc.Apply(acctID, 100)
	if err != nil {
		t.Fatal(err)
	}
	var balance float64
	db.Get(&balance, "SELECT balance FROM accounts_account WHERE id=?", acctID)
	if balance != 600 {
		t.Errorf("expected 600, got %.2f", balance)
	}
}

func TestBalanceReverse(t *testing.T) {
	db := testutil.SetupDB()
	defer db.Close()
	svc := NewBalanceService(db)

	_, acctID := insertUserBudgetAccount(db)

	err := svc.Reverse(acctID, 100)
	if err != nil {
		t.Fatal(err)
	}
	var balance float64
	db.Get(&balance, "SELECT balance FROM accounts_account WHERE id=?", acctID)
	if balance != 400 {
		t.Errorf("expected 400, got %.2f", balance)
	}
}

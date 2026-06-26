package service

import (
	"strings"
	"testing"

	"github.com/cesardarinel/chele-api/internal/middleware"
	"github.com/cesardarinel/chele-api/internal/testutil"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

func insertBudgetAndCategory(db *sqlx.DB) (string, string) {
	hash := middleware.MakeDjangoPassword("pass")
	db.MustExec("INSERT INTO auth_user (id,password,username,email,is_superuser,is_staff,is_active,first_name,last_name,date_joined) VALUES (1,?,'t@t.com','t@t.com',0,0,1,'T','U',datetime('now'))", hash)

	budgetID := strings.ReplaceAll(uuid.New().String(), "-", "")
	db.MustExec("INSERT INTO budgets_budget (id,name,description,owner_id,created_at,updated_at) VALUES (?,'Test','',1,datetime('now'),datetime('now'))", budgetID)

	groupID := strings.ReplaceAll(uuid.New().String(), "-", "")
	db.MustExec("INSERT INTO budgets_categorygroup (id,budget_id,name,sort_order,is_income,created_at,updated_at) VALUES (?,?,'Test Group',0,0,datetime('now'),datetime('now'))", groupID, budgetID)

	catID := strings.ReplaceAll(uuid.New().String(), "-", "")
	db.MustExec("INSERT INTO budgets_category (id,budget_id,group_id,name,sort_order,is_hidden,notes,created_at,updated_at) VALUES (?,?,?,'Test Cat',0,0,'',datetime('now'),datetime('now'))", catID, budgetID, groupID)

	return budgetID, catID
}

func insertTarget(db *sqlx.DB, catID, goalType string, amount float64, refillUpTo bool, snoozeM, snoozeY *int) string {
	id := strings.ReplaceAll(uuid.New().String(), "-", "")
	snoozeMVal := "NULL"
	snoozeYVal := "NULL"
	if snoozeM != nil {
		snoozeMVal = "?"
	}
	if snoozeY != nil {
		snoozeYVal = "?"
	}

	query := `INSERT INTO goals_goal (id,category_id,goal_type,amount,frequency,is_completed,refill_up_to,snooze_month,snooze_year,created_at,updated_at)
			   VALUES (?,?,?,?,12,0,?,` + snoozeMVal + `,` + snoozeYVal + `,datetime('now'),datetime('now'))`

	args := []interface{}{id, catID, goalType, amount, refillUpTo}
	if snoozeM != nil {
		args = append(args, *snoozeM)
	}
	if snoozeY != nil {
		args = append(args, *snoozeY)
	}
	db.MustExec(query, args...)
	return id
}

func TestTargetService_MonthlyUnderfunded(t *testing.T) {
	db := testutil.SetupDB()
	defer db.Close()
	svc := NewTargetService(db)

	_, catID := insertBudgetAndCategory(db)
	insertTarget(db, catID, "monthly", 500, false, nil, nil)

	underfunded := svc.CalculateUnderfunded("monthly", 500, false, nil, nil, catID, 6, 2026, 12, nil)
	if underfunded != 500 {
		t.Errorf("expected 500, got %.2f", underfunded)
	}
}

func TestTargetService_RefillUpTo(t *testing.T) {
	db := testutil.SetupDB()
	defer db.Close()
	svc := NewTargetService(db)

	_, catID := insertBudgetAndCategory(db)
	insertTarget(db, catID, "monthly", 200, true, nil, nil)

	underfunded := svc.CalculateUnderfunded("monthly", 200, true, nil, nil, catID, 6, 2026, 12, nil)
	if underfunded > 200 {
		t.Errorf("expected <=200 with refill_up_to, got %.2f", underfunded)
	}
}

func TestTargetService_Snoozed(t *testing.T) {
	db := testutil.SetupDB()
	defer db.Close()
	svc := NewTargetService(db)

	_, catID := insertBudgetAndCategory(db)
	snoozeM := 6
	snoozeY := 2026
	insertTarget(db, catID, "monthly", 500, false, &snoozeM, &snoozeY)

	underfunded := svc.CalculateUnderfunded("monthly", 500, false, &snoozeM, &snoozeY, catID, 6, 2026, 12, nil)
	if underfunded != 0 {
		t.Errorf("expected 0 for snoozed target, got %.2f", underfunded)
	}
}

func TestTargetService_Yearly(t *testing.T) {
	db := testutil.SetupDB()
	defer db.Close()
	svc := NewTargetService(db)

	_, catID := insertBudgetAndCategory(db)
	insertTarget(db, catID, "yearly", 1200, false, nil, nil)

	underfunded := svc.CalculateUnderfunded("yearly", 1200, false, nil, nil, catID, 6, 2026, 12, nil)
	if underfunded != 100 {
		t.Errorf("expected 100 (1200/12), got %.2f", underfunded)
	}
}

func TestTargetService_TrueExpense(t *testing.T) {
	db := testutil.SetupDB()
	defer db.Close()
	svc := NewTargetService(db)

	_, catID := insertBudgetAndCategory(db)
	insertTarget(db, catID, "true_expense", 600, false, nil, nil)

	underfunded := svc.CalculateUnderfunded("true_expense", 600, false, nil, nil, catID, 6, 2026, 12, nil)
	if underfunded != 50 {
		t.Errorf("expected 50 (600/12), got %.2f", underfunded)
	}
}

func TestTargetService_ListUnderfunded(t *testing.T) {
	db := testutil.SetupDB()
	defer db.Close()
	svc := NewTargetService(db)

	budgetID, catID := insertBudgetAndCategory(db)
	insertTarget(db, catID, "monthly", 300, false, nil, nil)

	result := svc.ListUnderfunded(budgetID, 6, 2026)
	if len(result) != 1 {
		t.Errorf("expected 1 underfunded category, got %d", len(result))
	}
	if len(result) > 0 && result[0].Deficit != 300 {
		t.Errorf("expected deficit 300, got %.2f", result[0].Deficit)
	}
}

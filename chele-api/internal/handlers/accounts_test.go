package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/cesardarinel/chele-api/internal/middleware"
	"github.com/cesardarinel/chele-api/internal/testutil"
	"github.com/google/uuid"
)

func setupAccountHandler(t *testing.T) (*AccountHandler, string) {
	t.Helper()
	db := testutil.SetupDB()
	t.Cleanup(func() { db.Close() })
	h := NewAccountHandler(db)

	hash := middleware.MakeDjangoPassword("pass")
	db.MustExec("INSERT INTO auth_user (id,password,username,email,is_superuser,is_staff,is_active,first_name,last_name,date_joined) VALUES (1,?,'t@t.com','t@t.com',0,0,1,'T','U',datetime('now'))", hash)

	budgetID := strings.ReplaceAll(uuid.New().String(), "-", "")
	db.MustExec("INSERT INTO budgets_budget (id,name,description,owner_id,created_at,updated_at) VALUES (?,'Test','',1,datetime('now'),datetime('now'))", budgetID)

	return h, budgetID
}

func TestAccountCreate(t *testing.T) {
	h, budgetID := setupAccountHandler(t)
	ctx := context.WithValue(context.Background(), middleware.UserIDKey, 1)

	body := `{"budget_id":"` + budgetID + `","name":"New Account","on_budget":true}`
	req := httptest.NewRequest("POST", "/api/accounts", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()
	h.Create(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d: %s", w.Code, w.Body.String())
	}
	var resp map[string]string
	json.NewDecoder(w.Body).Decode(&resp)
	if resp["id"] == "" {
		t.Error("expected id in response")
	}
}

func TestAccountList(t *testing.T) {
	h, budgetID := setupAccountHandler(t)
	ctx := context.WithValue(context.Background(), middleware.UserIDKey, 1)

	req := httptest.NewRequest("GET", "/api/accounts?budget_id="+budgetID, nil)
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()
	h.List(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	var resp struct {
		OnBudgetAccounts []interface{} `json:"on_budget"`
		TotalOnBudget    float64       `json:"total_on_budget"`
	}
	json.NewDecoder(w.Body).Decode(&resp)
	if resp.TotalOnBudget != 0 {
		t.Errorf("expected 0 total, got %.2f", resp.TotalOnBudget)
	}
}

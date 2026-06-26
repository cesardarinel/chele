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
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

func setupYNABTest(t *testing.T) (*YNABHandler, string, string, string) {
	t.Helper()
	db := testutil.SetupDB()
	t.Cleanup(func() { db.Close() })
	h := NewYNABHandler(db)

	hash := middleware.MakeDjangoPassword("pass")
	db.MustExec("INSERT INTO auth_user (id,password,username,email,is_superuser,is_staff,is_active,first_name,last_name,date_joined) VALUES (1,?,'t@t.com','t@t.com',0,0,1,'T','U',datetime('now'))", hash)

	budgetID := strings.ReplaceAll(uuid.New().String(), "-", "")
	db.MustExec("INSERT INTO budgets_budget (id,name,description,owner_id,created_at,updated_at) VALUES (?,'Test','',1,datetime('now'),datetime('now'))", budgetID)

	groupID := strings.ReplaceAll(uuid.New().String(), "-", "")
	db.MustExec("INSERT INTO budgets_categorygroup (id,budget_id,name,sort_order,is_income,created_at,updated_at) VALUES (?,?,'G',0,0,datetime('now'),datetime('now'))", groupID, budgetID)

	catID := strings.ReplaceAll(uuid.New().String(), "-", "")
	db.MustExec("INSERT INTO budgets_category (id,budget_id,group_id,name,sort_order,is_hidden,notes,created_at,updated_at) VALUES (?,?,?,'Cat1',0,0,'',datetime('now'),datetime('now'))", catID, budgetID, groupID)

	acctID := strings.ReplaceAll(uuid.New().String(), "-", "")
	db.MustExec("INSERT INTO accounts_account (id,budget_id,name,on_budget,balance,notes,created_at,updated_at) VALUES (?,?,'TestAcct',1,1000,'',datetime('now'),datetime('now'))", acctID, budgetID)

	db.MustExec("INSERT INTO budgets_monthlybudget (id,category_id,month,year,budgeted,created_at,updated_at) VALUES (?,?,6,2026,300,datetime('now'),datetime('now'))", strings.ReplaceAll(uuid.New().String(), "-", ""), catID)

	return h, budgetID, catID, groupID
}

func withChiURLParams(r *http.Request, params map[string]string) *http.Request {
	rctx := chi.NewRouteContext()
	for k, v := range params {
		rctx.URLParams.Add(k, v)
	}
	return r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))
}

func TestReadyToAssign(t *testing.T) {
	h, budgetID, _, _ := setupYNABTest(t)
	req := httptest.NewRequest("GET", "/api/budgets/"+budgetID+"/ready-to-assign?month=6&year=2026", nil)
	req = withChiURLParams(req, map[string]string{"id": budgetID})
	w := httptest.NewRecorder()
	h.ReadyToAssign(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	var resp map[string]float64
	json.NewDecoder(w.Body).Decode(&resp)
	if resp["total_on_budget"] != 1000 {
		t.Errorf("expected 1000 on budget, got %.2f", resp["total_on_budget"])
	}
}

func TestCover_Success(t *testing.T) {
	h, budgetID, catID, groupID := setupYNABTest(t)

	cat2ID := strings.ReplaceAll(uuid.New().String(), "-", "")
	h.DB.Exec("INSERT INTO budgets_category (id,budget_id,group_id,name,sort_order,is_hidden,notes,created_at,updated_at) VALUES (?,?,?,'Cat2',0,0,'',datetime('now'),datetime('now'))", cat2ID, budgetID, groupID)
	h.DB.Exec("INSERT INTO budgets_monthlybudget (id,category_id,month,year,budgeted,created_at,updated_at) VALUES (?,?,6,2026,500,datetime('now'),datetime('now'))", strings.ReplaceAll(uuid.New().String(), "-", ""), cat2ID)

	body := `{"from_category":"` + cat2ID + `","to_category":"` + catID + `","amount":100,"month":6,"year":2026}`
	req := httptest.NewRequest("POST", "/api/budgets/"+budgetID+"/cover", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req = withChiURLParams(req, map[string]string{"id": budgetID})
	w := httptest.NewRecorder()
	h.Cover(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}
}

func TestCover_InsufficientFunds(t *testing.T) {
	h, budgetID, catID, groupID := setupYNABTest(t)

	cat2ID := strings.ReplaceAll(uuid.New().String(), "-", "")
	h.DB.Exec("INSERT INTO budgets_category (id,budget_id,group_id,name,sort_order,is_hidden,notes,created_at,updated_at) VALUES (?,?,?,'Cat2',0,0,'',datetime('now'),datetime('now'))", cat2ID, budgetID, groupID)
	h.DB.Exec("INSERT INTO budgets_monthlybudget (id,category_id,month,year,budgeted,created_at,updated_at) VALUES (?,?,6,2026,10,datetime('now'),datetime('now'))", strings.ReplaceAll(uuid.New().String(), "-", ""), cat2ID)

	body := `{"from_category":"` + cat2ID + `","to_category":"` + catID + `","amount":100,"month":6,"year":2026}`
	req := httptest.NewRequest("POST", "/api/budgets/"+budgetID+"/cover", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req = withChiURLParams(req, map[string]string{"id": budgetID})
	w := httptest.NewRecorder()
	h.Cover(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400 for insufficient funds, got %d: %s", w.Code, w.Body.String())
	}
}

func TestSpotlight(t *testing.T) {
	h, budgetID, _, _ := setupYNABTest(t)
	req := httptest.NewRequest("GET", "/api/budgets/"+budgetID+"/spotlight?month=6&year=2026", nil)
	req = withChiURLParams(req, map[string]string{"id": budgetID})
	w := httptest.NewRecorder()
	h.Spotlight(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	var resp map[string]interface{}
	json.NewDecoder(w.Body).Decode(&resp)
	if resp["total_alerts"] == nil {
		t.Error("expected total_alerts in response")
	}
}

func TestInspector(t *testing.T) {
	h, _, catID, _ := setupYNABTest(t)
	req := httptest.NewRequest("GET", "/api/categories/"+catID+"/inspector?month=6&year=2026", nil)
	req = withChiURLParams(req, map[string]string{"id": catID})
	w := httptest.NewRecorder()
	h.Inspector(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	var resp map[string]interface{}
	json.NewDecoder(w.Body).Decode(&resp)
	if resp["category_id"] != catID {
		t.Errorf("expected category_id %s, got %v", catID, resp["category_id"])
	}
}

func TestCostToBeMe(t *testing.T) {
	h, budgetID, catID, _ := setupYNABTest(t)
	h.DB.Exec("INSERT INTO goals_goal (id,category_id,goal_type,amount,frequency,is_completed,refill_up_to,created_at,updated_at) VALUES (?,?,'monthly',500,12,0,0,datetime('now'),datetime('now'))", strings.ReplaceAll(uuid.New().String(), "-", ""), catID)

	req := httptest.NewRequest("GET", "/api/budgets/"+budgetID+"/cost-to-be-me", nil)
	req = withChiURLParams(req, map[string]string{"id": budgetID})
	w := httptest.NewRecorder()
	h.CostToBeMe(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	var resp map[string]interface{}
	json.NewDecoder(w.Body).Decode(&resp)
	if resp["cost_to_be_me"] == nil {
		t.Error("expected cost_to_be_me in response")
	}
}

func TestRollover(t *testing.T) {
	h, budgetID, _, _ := setupYNABTest(t)
	req := httptest.NewRequest("GET", "/api/budgets/"+budgetID+"/rollover?from_month=5&from_year=2026", nil)
	req = withChiURLParams(req, map[string]string{"id": budgetID})
	w := httptest.NewRecorder()
	h.Rollover(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
}

func TestAutoAssign(t *testing.T) {
	h, budgetID, catID, _ := setupYNABTest(t)
	h.DB.Exec("INSERT INTO goals_goal (id,category_id,goal_type,amount,frequency,is_completed,refill_up_to,created_at,updated_at) VALUES (?,?,'monthly',100,12,0,0,datetime('now'),datetime('now'))", strings.ReplaceAll(uuid.New().String(), "-", ""), catID)

	body := `{"month":6,"year":2026}`
	req := httptest.NewRequest("POST", "/api/budgets/"+budgetID+"/auto-assign", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req = withChiURLParams(req, map[string]string{"id": budgetID})
	w := httptest.NewRecorder()
	h.AutoAssign(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}
	var resp map[string]interface{}
	json.NewDecoder(w.Body).Decode(&resp)
	if resp["assigned_categories"] == nil {
		t.Error("expected assigned_categories in response")
	}
}

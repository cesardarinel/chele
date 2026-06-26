package handlers

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/cesardarinel/chele-api/internal/testutil"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

func setupTargetTest(t *testing.T) (*TargetHandler, string) {
	t.Helper()
	db := testutil.SetupDB()
	t.Cleanup(func() { db.Close() })
	h := NewTargetHandler(db)

	hash := "$2a$10$dummyhash"
	db.MustExec("INSERT INTO auth_user (id,password,username,email,is_superuser,is_staff,is_active,first_name,last_name,date_joined) VALUES (1,?,'t@t.com','t@t.com',0,0,1,'T','U',datetime('now'))", hash)

	budgetID := strings.ReplaceAll(uuid.New().String(), "-", "")
	db.MustExec("INSERT INTO budgets_budget (id,name,description,owner_id,created_at,updated_at) VALUES (?,'Test','',1,datetime('now'),datetime('now'))", budgetID)

	groupID := strings.ReplaceAll(uuid.New().String(), "-", "")
	db.MustExec("INSERT INTO budgets_categorygroup (id,budget_id,name,sort_order,is_income,created_at,updated_at) VALUES (?,?,'G',0,0,datetime('now'),datetime('now'))", groupID, budgetID)

	catID := strings.ReplaceAll(uuid.New().String(), "-", "")
	db.MustExec("INSERT INTO budgets_category (id,budget_id,group_id,name,sort_order,is_hidden,notes,created_at,updated_at) VALUES (?,?,?,'Cat1',0,0,'',datetime('now'),datetime('now'))", catID, budgetID, groupID)

	return h, catID
}

func TestTargetCreate(t *testing.T) {
	h, catID := setupTargetTest(t)
	body := `{"category_id":"` + catID + `","goal_type":"monthly","amount":500}`
	req := httptest.NewRequest("POST", "/api/targets", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	h.Create(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d: %s", w.Code, w.Body.String())
	}
}

func TestTargetList(t *testing.T) {
	h, catID := setupTargetTest(t)
	h.DB.Exec("INSERT INTO goals_goal (id,category_id,goal_type,amount,frequency,is_completed,refill_up_to,created_at,updated_at) VALUES (?,'monthly',300,12,0,0,datetime('now'),datetime('now'))", strings.ReplaceAll(uuid.New().String(), "-", ""), catID)

	req := httptest.NewRequest("GET", "/api/targets?category_id="+catID, nil)
	w := httptest.NewRecorder()
	h.List(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
}

func TestTargetGet(t *testing.T) {
	h, catID := setupTargetTest(t)
	targetID := strings.ReplaceAll(uuid.New().String(), "-", "")
	h.DB.Exec("INSERT INTO goals_goal (id,category_id,goal_type,amount,frequency,is_completed,refill_up_to,created_at,updated_at) VALUES (?,?,'monthly',300,12,0,0,datetime('now'),datetime('now'))", targetID, catID)

	req := httptest.NewRequest("GET", "/api/targets/"+targetID, nil)
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", targetID)
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
	w := httptest.NewRecorder()
	h.Get(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
}

func TestTargetDelete(t *testing.T) {
	h, catID := setupTargetTest(t)
	targetID := strings.ReplaceAll(uuid.New().String(), "-", "")
	h.DB.Exec("INSERT INTO goals_goal (id,category_id,goal_type,amount,frequency,is_completed,refill_up_to,created_at,updated_at) VALUES (?,?,'monthly',300,12,0,0,datetime('now'),datetime('now'))", targetID, catID)

	req := httptest.NewRequest("DELETE", "/api/targets/"+targetID, nil)
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", targetID)
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
	w := httptest.NewRecorder()
	h.Delete(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
}

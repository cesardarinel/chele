package router

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/cesardarinel/chele-api/internal/config"
	"github.com/cesardarinel/chele-api/internal/handlers"
	"github.com/cesardarinel/chele-api/internal/middleware"
	"github.com/cesardarinel/chele-api/internal/models"
	"github.com/cesardarinel/chele-api/internal/service"
	"github.com/go-chi/chi/v5"
	"github.com/jmoiron/sqlx"
)

func New(db *sqlx.DB, cfg *config.Config) *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.CORS)
	authMw := middleware.JWTAuth(cfg.JWTSecert)

	auth := handlers.NewAuthHandler(db, cfg.JWTSecert)
	r.Post("/api/auth/login", auth.Login)
	r.Post("/api/auth/register", auth.Register)
	r.With(authMw).Get("/api/auth/me", auth.Me)

	acct := handlers.NewAccountHandler(db)
	r.With(authMw).Get("/api/accounts", acct.List)
	r.With(authMw).Post("/api/accounts", acct.Create)
	r.With(authMw).Get("/api/accounts/{id}", acct.Get)
	r.With(authMw).Put("/api/accounts/{id}", acct.Update)
	r.With(authMw).Delete("/api/accounts/{id}", acct.Delete)

	bgt := handlers.NewBudgetHandler(db)
	r.With(authMw).Get("/api/budgets", bgt.List)
	r.With(authMw).Post("/api/budgets", bgt.Create)
	r.With(authMw).Get("/api/budgets/{id}", bgt.Get)
	r.With(authMw).Put("/api/budgets/{id}", bgt.Update)

	dashSvc := service.NewDashboardService(db)
	r.With(authMw).Get("/api/budgets/{id}/dashboard", func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		month := queryInt(r, "mes")
		year := queryInt(r, "anio")
		rangeVal := queryInt(r, "rango")
		if month == 0 {
			month = int(time.Now().Month())
		}
		if year == 0 {
			year = time.Now().Year()
		}
		resp, err := dashSvc.GetBudgetView(id, month, year, rangeVal)
		if err != nil {
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}
		json.NewEncoder(w).Encode(resp)
	})

	txn := handlers.NewTransactionHandler(db)
	r.With(authMw).Get("/api/transactions", txn.List)
	r.With(authMw).Post("/api/transactions", txn.Create)
	r.With(authMw).Put("/api/transactions/{id}", txn.Update)
	r.With(authMw).Delete("/api/transactions/{id}", txn.Delete)
	r.With(authMw).Post("/api/transactions/bulk", txn.Bulk)

	cat := handlers.NewCategoryHandler(db)
	r.With(authMw).Get("/api/category-groups", cat.ListGroups)
	r.With(authMw).Post("/api/category-groups", cat.CreateGroup)
	r.With(authMw).Put("/api/category-groups/{id}", cat.UpdateGroup)
	r.With(authMw).Delete("/api/category-groups/{id}", cat.DeleteGroup)
	r.With(authMw).Get("/api/categories", cat.List)
	r.With(authMw).Post("/api/categories", cat.Create)
	r.With(authMw).Put("/api/categories/{id}", cat.Update)
	r.With(authMw).Delete("/api/categories/{id}", cat.Delete)

	mb := handlers.NewMonthlyBudgetHandler(db)
	r.With(authMw).Get("/api/monthly-budgets", mb.List)
	r.With(authMw).Put("/api/monthly-budgets", mb.Assign)
	r.With(authMw).Post("/api/monthly-budgets/move", mb.Move)

	sched := handlers.NewScheduleHandler(db)
	r.With(authMw).Get("/api/schedules", sched.List)
	r.With(authMw).Post("/api/schedules", sched.Create)
	r.With(authMw).Put("/api/schedules/{id}", sched.Update)
	r.With(authMw).Delete("/api/schedules/{id}", sched.Delete)
	r.With(authMw).Post("/api/schedules/process", sched.Process)

	cc := handlers.NewCreditCardHandler(db)
	r.With(authMw).Get("/api/credit-cards", cc.List)
	r.With(authMw).Post("/api/credit-cards", cc.Create)
	r.With(authMw).Get("/api/credit-cards/{id}", cc.Get)
	r.With(authMw).Put("/api/credit-cards/{id}", cc.Update)
	r.With(authMw).Delete("/api/credit-cards/{id}", cc.Delete)
	r.With(authMw).Post("/api/credit-cards/{id}/pay", cc.Pay)

	loan := handlers.NewLoanHandler(db)
	r.With(authMw).Get("/api/loans", loan.List)
	r.With(authMw).Post("/api/loans", loan.Create)
	r.With(authMw).Get("/api/loans/{id}", loan.Get)
	r.With(authMw).Post("/api/loans/{id}/pay-installment", loan.PayInstallment)

	payee := handlers.NewPayeeHandler(db)
	r.With(authMw).Get("/api/payees", payee.List)
	r.With(authMw).Post("/api/payees", payee.Create)
	r.With(authMw).Put("/api/payees/{id}", payee.Update)
	r.With(authMw).Delete("/api/payees/{id}", payee.Delete)
	r.With(authMw).Post("/api/payees/{id}/merge", payee.Merge)

	gl := handlers.NewGoalHandler(db)
	r.With(authMw).Get("/api/goals", gl.List)
	r.With(authMw).Post("/api/goals", gl.Create)
	r.With(authMw).Put("/api/goals/{id}", gl.Update)
	r.With(authMw).Delete("/api/goals/{id}", gl.Delete)

	rule := handlers.NewRuleHandler(db)
	r.With(authMw).Get("/api/rules", rule.List)
	r.With(authMw).Post("/api/rules", rule.Create)
	r.With(authMw).Put("/api/rules/{id}", rule.Update)
	r.With(authMw).Delete("/api/rules/{id}", rule.Delete)

	rpt := handlers.NewReportHandler(db)
	r.With(authMw).Get("/api/reports/net-worth", rpt.NetWorth)
	r.With(authMw).Get("/api/reports/cash-flow", rpt.CashFlow)
	r.With(authMw).Get("/api/reports/budget-vs-reality", rpt.BudgetVsReality)

	stg := handlers.NewSettingsHandler(db)
	r.With(authMw).Get("/api/settings/profile", stg.Profile)
	r.With(authMw).Put("/api/settings/profile", stg.UpdateProfile)
	r.With(authMw).Get("/api/settings/budget/{id}", stg.BudgetSettings)
	r.With(authMw).Post("/api/settings/budget/{id}/invite", stg.Invite)

	r.With(authMw).Get("/api/sync/logs", func(w http.ResponseWriter, r *http.Request) {
		budgetID := r.URL.Query().Get("budget_id")
		q := "SELECT * FROM sync_engine_synclog WHERE 1=1"
		var args []interface{}
		if budgetID != "" { q += " AND budget_id=?"; args = append(args, budgetID) }
		q += " ORDER BY created_at DESC LIMIT 100"
		var logs []models.SyncLog
		db.Select(&logs, q, args...)
		json.NewEncoder(w).Encode(logs)
	})

	return r
}

func queryInt(r *http.Request, key string) int {
	v := r.URL.Query().Get(key)
	if v == "" { return 0 }
	n, _ := strconv.Atoi(v)
	return n
}

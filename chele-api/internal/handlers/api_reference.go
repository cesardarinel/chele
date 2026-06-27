package handlers

import (
	"encoding/json"
	"net/http"
)

type APIEndpoint struct {
	Method      string `json:"method"`
	Path        string `json:"path"`
	Auth        string `json:"auth"`
	Description string `json:"description"`
	Request     string `json:"request,omitempty"`
	Response    string `json:"response,omitempty"`
}

type APIDoc struct {
	Title       string                 `json:"title"`
	BaseURL     string                 `json:"base_url"`
	AuthScheme  string                 `json:"auth_scheme"`
	Groups      map[string][]APIEndpoint `json:"groups"`
}

func APIRefHandler(w http.ResponseWriter, r *http.Request) {
	doc := APIDoc{
		Title:      "Chele API Reference",
		BaseURL:    "http://localhost:8080/api",
		AuthScheme: "Bearer token via Authorization: Bearer <token> header",
		Groups: map[string][]APIEndpoint{
			"Auth": {
				{Method: "POST", Path: "/auth/login", Auth: "no", Description: "Authenticate user with email/password against Django's PBKDF2 hash. Returns JWT token.", Request: `{"username":"user@email.com","password":"secret123"}`, Response: `{"token":"eyJ...","user":{"id":1,"username":"...","email":"...","first_name":"...","last_name":"..."}}`},
				{Method: "POST", Path: "/auth/register", Auth: "no", Description: "Create a new user account. Returns JWT token.", Request: `{"username":"user@email.com","email":"user@email.com","password":"secret123","name":"John Doe"}`, Response: `{"token":"eyJ...","user_id":1,"username":"...","email":"..."}`},
				{Method: "GET", Path: "/auth/me", Auth: "yes", Description: "Get current authenticated user's profile.", Response: `{"id":1,"username":"...","email":"..."}`},
			},
			"Accounts": {
				{Method: "GET", Path: "/accounts?budget_id=<uuid>", Auth: "yes", Description: "List all accounts with balance summary: on/off budget totals, credit card debt, loans.", Response: `{"on_budget":[...],"off_budget":[...],"total_on_budget":5000,"grand_total":5000,"cc_debt_cards":[...],"total_cc_debt":-1000,"loans":[...],"total_debt":-15000}`},
				{Method: "POST", Path: "/accounts", Auth: "yes", Description: "Create a new account.", Request: `{"budget_id":"uuid","name":"Checking","on_budget":true}`, Response: `{"id":"uuid"}`},
				{Method: "GET", Path: "/accounts/{id}", Auth: "yes", Description: "Get account detail with transactions and active schedules.", Response: `{"account":{...},"transactions":[...],"schedules":[...]}`},
				{Method: "PUT", Path: "/accounts/{id}", Auth: "yes", Description: "Update account name and on_budget flag.", Request: `{"name":"New Name","on_budget":true}`},
				{Method: "DELETE", Path: "/accounts/{id}", Auth: "yes", Description: "Delete an account."},
			},
			"Budgets": {
				{Method: "GET", Path: "/budgets", Auth: "yes", Description: "List budgets where user is a member.", Response: `[{"id":"uuid","name":"Budget","owner_id":1,...}]`},
				{Method: "POST", Path: "/budgets", Auth: "yes", Description: "Create budget with default groups and categories.", Request: `{"name":"My Budget"}`, Response: `{"id":"uuid"}`},
				{Method: "GET", Path: "/budgets/{id}", Auth: "yes", Description: "Get budget details."},
				{Method: "PUT", Path: "/budgets/{id}", Auth: "yes", Description: "Update budget name and description.", Request: `{"name":"New Name","description":"..."}`},
				{Method: "GET", Path: "/budgets/{id}/dashboard?mes=&anio=&rango=", Auth: "yes", Description: "Full budget view: groups, categories, months, income/expenses, available to budget. rango: 0=1 month, 1=3 months, 2=5 months.", Response: `{"budget":{...},"groups":[...],"months":[...],"available_to_budget":2000,"total_balance":5000}`},
			},
			"Transactions": {
				{Method: "GET", Path: "/transactions?budget_id=&account_id=&month=&year=&category_id=", Auth: "yes", Description: "List transactions with optional filters.", Response: `[{"id":"uuid","amount":-50,"date":"2026-06-01","payee":{...},"category":{...}}]`},
				{Method: "POST", Path: "/transactions", Auth: "yes", Description: "Create transaction, updates account balance. TC overspend auto-moves to Payment category.", Request: `{"budget_id":"uuid","account_id":"uuid","date":"2026-06-01","amount":100,"direction":"expense|income","payee_id":"uuid|null","category_id":"uuid|null","notes":"..."}`, Response: `{"id":"uuid"}`},
				{Method: "PUT", Path: "/transactions/{id}", Auth: "yes", Description: "Update transaction (reverses old balance, applies new).", Request: `{"amount":200,"direction":"income"}`},
				{Method: "DELETE", Path: "/transactions/{id}", Auth: "yes", Description: "Delete transaction (reverses account balance)."},
				{Method: "POST", Path: "/transactions/bulk", Auth: "yes", Description: "Bulk delete or categorize transactions.", Request: `{"action":"delete|categorize","ids":["uuid1","uuid2"],"category_id":"uuid"}`},
			},
			"Category Groups": {
				{Method: "GET", Path: "/category-groups?budget_id=<uuid>", Auth: "yes", Description: "List category groups."},
				{Method: "POST", Path: "/category-groups", Auth: "yes", Description: "Create a category group.", Request: `{"budget_id":"uuid","name":"Group Name","is_income":false}`},
				{Method: "PUT", Path: "/category-groups/{id}", Auth: "yes", Description: "Rename a category group.", Request: `{"name":"New Name"}`},
				{Method: "DELETE", Path: "/category-groups/{id}", Auth: "yes", Description: "Delete group and its categories."},
			},
			"Categories": {
				{Method: "GET", Path: "/categories?budget_id=<uuid>", Auth: "yes", Description: "List non-hidden categories."},
				{Method: "POST", Path: "/categories", Auth: "yes", Description: "Create a category.", Request: `{"budget_id":"uuid","group_id":"uuid","name":"Category"}`},
				{Method: "PUT", Path: "/categories/{id}", Auth: "yes", Description: "Rename a category.", Request: `{"name":"New Name"}`},
				{Method: "DELETE", Path: "/categories/{id}", Auth: "yes", Description: "Soft-delete category (sets is_hidden=True)."},
			},
			"Monthly Budget": {
				{Method: "GET", Path: "/monthly-budgets?category_id=&month=&year=", Auth: "yes", Description: "List monthly budget assignments."},
				{Method: "PUT", Path: "/monthly-budgets", Auth: "yes", Description: "Assign funds to a category for a month/year.", Request: `{"category_id":"uuid","month":6,"year":2026,"amount":500}`},
				{Method: "POST", Path: "/monthly-budgets/move", Auth: "yes", Description: "Move funds between categories.", Request: `{"from_category":"uuid","to_category":"uuid","amount":100,"month":6,"year":2026}`},
			},
			"Schedules": {
				{Method: "GET", Path: "/schedules?budget_id=<uuid>", Auth: "yes", Description: "List recurring schedules.", Response: `[{"id":"uuid","amount":500,"frequency":"monthly","next_date":"2026-07-01","direction":"expense|income","skip_weekends":false,"apply_before_weekend":false}]`},
				{Method: "POST", Path: "/schedules", Auth: "yes", Description: "Create a schedule.", Request: `{"budget_id":"uuid","account_id":"uuid","amount":500,"direction":"expense|income","frequency":"weekly|biweekly|monthly|quarterly|yearly","next_date":"2026-07-01","payee_id":"uuid|null","category_id":"uuid|null","skip_weekends":false,"apply_before_weekend":false}`},
				{Method: "PUT", Path: "/schedules/{id}", Auth: "yes", Description: "Update schedule fields."},
				{Method: "DELETE", Path: "/schedules/{id}", Auth: "yes", Description: "Delete a schedule."},
				{Method: "POST", Path: "/schedules/process?budget_id=<uuid>", Auth: "yes", Description: "Process due schedules (creates transactions, advances dates).", Response: `{"processed":3}`},
			},
			"Credit Cards": {
				{Method: "GET", Path: "/credit-cards?budget_id=<uuid>", Auth: "yes", Description: "List credit cards."},
				{Method: "POST", Path: "/credit-cards", Auth: "yes", Description: "Create credit card + Payment category for auto-move.", Request: `{"budget_id":"uuid","name":"Visa","limit":50000,"balance":-1000,"interest_rate":0.96,"closing_day":15,"due_day":5}`},
				{Method: "GET", Path: "/credit-cards/{id}", Auth: "yes", Description: "Get credit card detail."},
				{Method: "PUT", Path: "/credit-cards/{id}", Auth: "yes", Description: "Update credit card."},
				{Method: "DELETE", Path: "/credit-cards/{id}", Auth: "yes", Description: "Delete credit card."},
				{Method: "POST", Path: "/credit-cards/{id}/pay", Auth: "yes", Description: "Pay credit card (applies interest, creates transaction).", Request: `{"account_id":"uuid","amount":500}`},
			},
			"Loans": {
				{Method: "GET", Path: "/loans?budget_id=<uuid>", Auth: "yes", Description: "List loans."},
				{Method: "POST", Path: "/loans", Auth: "yes", Description: "Create loan with installment generation.", Request: `{"budget_id":"uuid","type":"personal|hipotecario|...","name":"Car Loan","total_amount":300000,"interest_rate":0.12,"total_installments":36,"installment_amount":9500,"next_due_date":"2026-02-01"}`, Response: `{"id":"uuid"}`},
				{Method: "GET", Path: "/loans/{id}", Auth: "yes", Description: "Get loan with installments.", Response: `{"loan":{...},"installments":[...]}`},
				{Method: "POST", Path: "/loans/{id}/pay-installment", Auth: "yes", Description: "Pay next installment (creates transaction).", Request: `{"account_id":"uuid"}`},
			},
			"Payees": {
				{Method: "GET", Path: "/payees?budget_id=<uuid>", Auth: "yes", Description: "List payees."},
				{Method: "POST", Path: "/payees", Auth: "yes", Description: "Create a payee.", Request: `{"budget_id":"uuid","name":"Payee Name"}`},
				{Method: "PUT", Path: "/payees/{id}", Auth: "yes", Description: "Rename a payee."},
				{Method: "DELETE", Path: "/payees/{id}", Auth: "yes", Description: "Delete a payee."},
				{Method: "POST", Path: "/payees/{id}/merge", Auth: "yes", Description: "Merge payee into another.", Request: `{"target_id":"uuid"}`},
			},
			"Goals": {
				{Method: "GET", Path: "/goals?budget_id=<uuid>", Auth: "yes", Description: "List goals."},
				{Method: "POST", Path: "/goals", Auth: "yes", Description: "Create a goal (applies to monthly budget).", Request: `{"category_id":"uuid","goal_type":"monthly|yearly|target_balance|target_date|true_expense","amount":500,"target_date":"2026-12-31","frequency":12}`},
				{Method: "PUT", Path: "/goals/{id}", Auth: "yes", Description: "Update a goal."},
				{Method: "DELETE", Path: "/goals/{id}", Auth: "yes", Description: "Delete a goal."},
			},
			"Rules": {
				{Method: "GET", Path: "/rules?budget_id=<uuid>", Auth: "yes", Description: "List automation rules."},
				{Method: "POST", Path: "/rules", Auth: "yes", Description: "Create auto-categorize rule.", Request: `{"budget_id":"uuid","name":"Rule","condition_field":"payee|amount|notes","condition_operator":"contains|equals|starts_with","condition_value":"text","action_category_id":"uuid"}`},
				{Method: "PUT", Path: "/rules/{id}", Auth: "yes", Description: "Update a rule."},
				{Method: "DELETE", Path: "/rules/{id}", Auth: "yes", Description: "Delete a rule."},
			},
			"Reports": {
				{Method: "GET", Path: "/reports/net-worth?budget_id=<uuid>", Auth: "yes", Description: "Net worth (assets - liabilities).", Response: `{"assets":50000,"liabilities":10000,"net_worth":40000}`},
				{Method: "GET", Path: "/reports/cash-flow?budget_id=<uuid>", Auth: "yes", Description: "Last 6 months income/expenses.", Response: `[{"month":"6/2026","income":5000,"expenses":3200,"net":1800}]`},
				{Method: "GET", Path: "/reports/budget-vs-reality?budget_id=<uuid>", Auth: "yes", Description: "Budget vs actual per category for current month."},
			},
			"YNAB Engine": {
				{Method: "GET", Path: "/budgets/{id}/ready-to-assign?month=&year=", Auth: "yes", Description: "Calculate Ready to Assign = on_budget_balance - total_assigned.", Response: `{"ready_to_assign":2000,"total_on_budget":5000,"total_assigned":3000}`},
				{Method: "POST", Path: "/budgets/{id}/auto-assign", Auth: "yes", Description: "Auto-distribute funds to underfunded categories (priority order).", Request: `{"month":6,"year":2026}`, Response: `{"assigned_categories":3,"total_assigned":900,"remaining":100}`},
				{Method: "POST", Path: "/budgets/{id}/cover", Auth: "yes", Description: "Move funds to cover an overspent category.", Request: `{"from_category":"uuid","to_category":"uuid","amount":50,"month":6,"year":2026}`},
				{Method: "GET", Path: "/budgets/{id}/spotlight?month=&year=", Auth: "yes", Description: "Aggregated alerts: uncategorized, overspends, underfunded.", Response: `{"uncategorized_count":3,"overspends":[...],"underfunded":[...],"total_alerts":2}`},
				{Method: "GET", Path: "/budgets/{id}/cost-to-be-me", Auth: "yes", Description: "Total target cost vs expected income.", Response: `{"cost_to_be_me":3500,"expected_income":4200,"difference":700,"is_over_budget":false}`},
				{Method: "GET", Path: "/budgets/{id}/rollover?from_month=&from_year=", Auth: "yes", Description: "Category rollover balances to next month.", Response: `{"categories":[{"category_id":"uuid","name":"Groceries","rollover_balance":40}]}`},
				{Method: "GET", Path: "/categories/{id}/inspector?month=&year=", Auth: "yes", Description: "Category detail for inspector panel.", Response: `{"category_id":"uuid","category_name":"Groceries","available":150,"assigned":200,"activity":-50,"target":{...},"average_spending_3m":180}`},
			},
			"Settings": {
				{Method: "GET", Path: "/settings/profile", Auth: "yes", Description: "Get user profile."},
				{Method: "PUT", Path: "/settings/profile", Auth: "yes", Description: "Update user profile.", Request: `{"first_name":"John","last_name":"Doe","email":"new@email.com"}`},
				{Method: "GET", Path: "/settings/budget/{id}", Auth: "yes", Description: "Get budget settings with members."},
				{Method: "POST", Path: "/settings/budget/{id}/invite", Auth: "yes", Description: "Invite user by email.", Request: `{"email":"user@email.com"}`},
			},
			"Sync": {
				{Method: "GET", Path: "/sync/logs?budget_id=<uuid>", Auth: "yes", Description: "Get sync logs (last 100)."},
			},
		},
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	json.NewEncoder(w).Encode(doc)
}

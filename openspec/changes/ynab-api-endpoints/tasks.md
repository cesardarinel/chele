## 1. Target Model

- [ ] 1.1 Add `Target` struct to `internal/models/models.go` with all fields matching Django's Target model
- [ ] 1.2 Add `Target` table to `internal/testutil/testutil.go` schema

## 2. Target Service

- [ ] 2.1 Implement `TargetService` with `CalculateUnderfunded(target, category, month, year)` supporting monthly/yearly/target_balance/target_date/true_expense
- [ ] 2.2 Implement refill_up_to calculation: underfunded = max(0, amount - category_available_balance)
- [ ] 2.3 Implement set_aside_another calculation: underfunded = amount
- [ ] 2.4 Implement snooze filtering: skip targets with snooze_month/snooze_year matching current month
- [ ] 2.5 Implement yearly proration: monthly = amount / 12
- [ ] 2.6 Implement target_date proration: monthly = (amount - already_assigned) / months_remaining

## 3. Target Handlers

- [ ] 3.1 Create `internal/handlers/targets.go` with List, Create, Get, Update, Delete methods
- [ ] 3.2 Implement CREATE with all target fields and validation
- [ ] 3.3 Implement UPDATE supporting partial updates (amount, goal_type, snooze, etc.)
- [ ] 3.4 Register routes in `internal/router/router.go`: `GET/POST/PUT/DELETE /api/targets`

## 4. Ready to Assign Endpoint

- [ ] 4.1 Create handler `GET /api/budgets/:id/ready-to-assign` with month/year query params
- [ ] 4.2 Implement calculation: `SELECT SUM(balance) FROM accounts_account WHERE budget_id=? AND on_budget=1` minus `SELECT SUM(budgeted) FROM budgets_monthlybudget WHERE category_id IN (SELECT id FROM budgets_category WHERE budget_id=?) AND month=? AND year=?`
- [ ] 4.3 Register route in router

## 5. Auto-assign Endpoint

- [ ] 5.1 Create `internal/service/autoassign.go` with `AssignAll(budgetID, month, year)` service
- [ ] 5.2 Implement priority ordering: cash overspends > true_expense > monthly targets > other types
- [ ] 5.3 Implement partial assignment: stop when Ready to Assign is exhausted
- [ ] 5.4 Create handler `POST /api/budgets/:id/auto-assign`
- [ ] 5.5 Register route

## 6. Cover Overspending Endpoint

- [ ] 6.1 Create handler `POST /api/budgets/:id/cover`
- [ ] 6.2 Implement atomic UPDATE: decrement source MonthlyBudget, increment target MonthlyBudget
- [ ] 6.3 Add validation: source must have sufficient available balance
- [ ] 6.4 Register route

## 7. Spotlight Alerts Endpoint

- [ ] 7.1 Create handler `GET /api/budgets/:id/spotlight`
- [ ] 7.2 Implement uncategorized transactions count query
- [ ] 7.3 Implement uncovered overspends query (categories with negative available balance)
- [ ] 7.4 Implement underfunded targets query (using TargetService)
- [ ] 7.5 Aggregate and return combined response
- [ ] 7.6 Register route

## 8. Inspector Endpoint

- [ ] 8.1 Create handler `GET /api/categories/:id/inspector`
- [ ] 8.2 Load MonthlyBudget for given category/month/year
- [ ] 8.3 Load Target for category (if any)
- [ ] 8.4 Calculate 3-month average spending from transactions
- [ ] 8.5 Return combined response
- [ ] 8.6 Register route

## 9. Cost to Be Me Endpoint

- [ ] 9.1 Create handler `GET /api/budgets/:id/cost-to-be-me`
- [ ] 9.2 Implement total monthly target cost calculation (sum of all active target monthly requirements)
- [ ] 9.3 Implement expected income calculation (average of last 3 months income transactions, amount > 0)
- [ ] 9.4 Compare and return difference with is_over_budget flag
- [ ] 9.5 Register route

## 10. Rollover Endpoint

- [ ] 10.1 Create handler `GET /api/budgets/:id/rollover`
- [ ] 10.2 For each category with MonthlyBudget data, calculate: `available + budgeted - spent = rollover_balance`
- [ ] 10.3 Return array of category_id, name, rollover_balance
- [ ] 10.4 Include refill_up_to underfunded adjustment in response
- [ ] 10.5 Register route

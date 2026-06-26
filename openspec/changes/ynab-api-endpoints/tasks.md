## 1. Target Model

- [x] 1.1 Add `Target` struct to `internal/models/models.go` with all fields matching Django's Target model
- [x] 1.2 Add `Target` table to `internal/testutil/testutil.go` schema

## 2. Target Service

- [x] 2.1 Implement `TargetService` with `CalculateUnderfunded(target, category, month, year)` supporting monthly/yearly/target_balance/target_date/true_expense
- [x] 2.2 Implement refill_up_to calculation: underfunded = max(0, amount - category_available_balance)
- [x] 2.3 Implement set_aside_another calculation: underfunded = amount
- [x] 2.4 Implement snooze filtering: skip targets with snooze_month/snooze_year matching current month
- [x] 2.5 Implement yearly proration: monthly = amount / 12
- [x] 2.6 Implement target_date proration: monthly = (amount - already_assigned) / months_remaining

## 3. Target Handlers

- [x] 3.1 Create `internal/handlers/targets.go` with List, Create, Get, Update, Delete methods
- [x] 3.2 Implement CREATE with all target fields and validation
- [x] 3.3 Implement UPDATE supporting partial updates (amount, goal_type, snooze, etc.)
- [x] 3.4 Register routes in `internal/router/router.go`: `GET/POST/PUT/DELETE /api/targets`

## 4. Ready to Assign Endpoint

- [x] 4.1 Create handler `GET /api/budgets/:id/ready-to-assign` with month/year query params
- [x] 4.2 Implement calculation: on-budget balance minus assigned funds
- [x] 4.3 Register route in router

## 5. Auto-assign Endpoint

- [x] 5.1 Create auto-assign logic in YNABHandler
- [x] 5.2 Implement priority ordering: cash overspends > true_expense > monthly targets > other types
- [x] 5.3 Implement partial assignment: stop when Ready to Assign is exhausted
- [x] 5.4 Create handler `POST /api/budgets/:id/auto-assign`
- [x] 5.5 Register route

## 6. Cover Overspending Endpoint

- [x] 6.1 Create handler `POST /api/budgets/:id/cover`
- [x] 6.2 Implement atomic UPDATE: decrement source MonthlyBudget, increment target MonthlyBudget
- [x] 6.3 Add validation: source must have sufficient available balance
- [x] 6.4 Register route

## 7. Spotlight Alerts Endpoint

- [x] 7.1 Create handler `GET /api/budgets/:id/spotlight`
- [x] 7.2 Implement uncategorized transactions count query
- [x] 7.3 Implement uncovered overspends query (categories with negative available balance)
- [x] 7.4 Implement underfunded targets query (using TargetService)
- [x] 7.5 Aggregate and return combined response
- [x] 7.6 Register route

## 8. Inspector Endpoint

- [x] 8.1 Create handler `GET /api/categories/:id/inspector`
- [x] 8.2 Load MonthlyBudget for given category/month/year
- [x] 8.3 Load Target for category (if any)
- [x] 8.4 Calculate 3-month average spending from transactions
- [x] 8.5 Return combined response
- [x] 8.6 Register route

## 9. Cost to Be Me Endpoint

- [x] 9.1 Create handler `GET /api/budgets/:id/cost-to-be-me`
- [x] 9.2 Implement total monthly target cost calculation (sum of all active target monthly requirements)
- [x] 9.3 Implement expected income calculation (average of last 3 months income transactions, amount > 0)
- [x] 9.4 Compare and return difference with is_over_budget flag
- [x] 9.5 Register route

## 10. Rollover Endpoint

- [x] 10.1 Create handler `GET /api/budgets/:id/rollover`
- [x] 10.2 For each category with MonthlyBudget data, calculate: `available + budgeted - spent = rollover_balance`
- [x] 10.3 Return array of category_id, name, rollover_balance
- [x] 10.4 Include refill_up_to underfunded adjustment in response
- [x] 10.5 Register route

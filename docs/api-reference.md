# Chele API Reference

**Base URL:** `http://localhost:8080/api`
**Auth:** Bearer token via `Authorization: Bearer <token>` header (except `/auth/login` and `/auth/register`)

---

## Auth

### POST /api/auth/login
- **Auth:** No
- **Description:** Authenticate user and get JWT token. Validates against Django's PBKDF2 password hash.
- **Request body:**
  ```json
  { "username": "user@email.com", "password": "secret123" }
  ```
- **Response 200:**
  ```json
  { "token": "eyJ...", "user": { "id": 1, "username": "...", "email": "...", "first_name": "...", "last_name": "..." } }
  ```
- **Response 401:** `{ "error": "invalid credentials" }`

### POST /api/auth/register
- **Auth:** No
- **Description:** Create a new user account
- **Request body:**
  ```json
  { "username": "user@email.com", "email": "user@email.com", "password": "secret123", "name": "John Doe" }
  ```
- **Response 201:**
  ```json
  { "token": "eyJ...", "user_id": 1, "username": "...", "email": "..." }
  ```

### GET /api/auth/me
- **Auth:** Yes
- **Description:** Get current authenticated user's profile
- **Response 200:** `{ "id": 1, "username": "...", "email": "...", "first_name": "...", "last_name": "...", ... }`

---

## Accounts

### GET /api/accounts?budget_id=<uuid>
- **Auth:** Yes
- **Description:** List all accounts with balance summary (on/off budget, totals, CC debt, loans)
- **Query params:** `budget_id` (required)
- **Response 200:** Full account summary with totals

### POST /api/accounts
- **Auth:** Yes
- **Description:** Create a new account
- **Request body:**
  ```json
  { "budget_id": "uuid", "name": "Checking", "on_budget": true }
  ```
- **Response 201:** `{ "id": "uuid" }`

### GET /api/accounts/{id}
- **Auth:** Yes
- **Description:** Get account detail with transactions and active schedules
- **Response 200:** `{ "account": {...}, "transactions": [...], "schedules": [...] }`

### PUT /api/accounts/{id}
- **Auth:** Yes
- **Description:** Update account name and on_budget flag
- **Request body:** `{ "name": "...", "on_budget": true }`

### DELETE /api/accounts/{id}
- **Auth:** Yes
- **Description:** Delete an account

---

## Budgets

### GET /api/budgets
- **Auth:** Yes
- **Description:** List budgets where the user is a member
- **Response 200:** `[{ "id": "uuid", "name": "...", "owner_id": 1, ... }]`

### POST /api/budgets
- **Auth:** Yes
- **Description:** Create a budget with default category groups and categories
- **Request body:** `{ "name": "My Budget" }`
- **Response 201:** `{ "id": "uuid" }`

### GET /api/budgets/{id}
- **Auth:** Yes
- **Description:** Get budget details

### PUT /api/budgets/{id}
- **Auth:** Yes
- **Description:** Update budget name and description

### GET /api/budgets/{id}/dashboard?mes=&anio=&rango=
- **Auth:** Yes
- **Description:** Full budget view with categories, months, income/expenses, available to budget
- **Query params:** `mes` (month 1-12), `anio` (year), `rango` (0=1m, 1=3m, 2=5m)
- **Response 200:** Complete dashboard data structure

---

## Transactions

### GET /api/transactions?budget_id=&account_id=&month=&year=&category_id=
- **Auth:** Yes
- **Description:** List transactions with optional filters
- **Response 200:** `[{ "id": "uuid", "amount": -50.00, "date": "2026-06-01", ... }]`

### POST /api/transactions
- **Auth:** Yes
- **Description:** Create a transaction (updates account balance). If TC overspend, auto-moves funds to Payment category.
- **Request body:**
  ```json
  { "budget_id": "uuid", "account_id": "uuid", "date": "2026-06-01", "amount": 100, "direction": "expense|income", "payee_id": "uuid|null", "category_id": "uuid|null", "notes": "..." }
  ```
- **Response 201:** `{ "id": "uuid" }`

### PUT /api/transactions/{id}
- **Auth:** Yes
- **Description:** Update transaction (reverses old balance, applies new)
- **Request body:** `{ "amount": 200, "direction": "income", ... }`

### DELETE /api/transactions/{id}
- **Auth:** Yes
- **Description:** Delete transaction (reverses account balance)

### POST /api/transactions/bulk
- **Auth:** Yes
- **Description:** Bulk delete or categorize transactions
- **Request body:**
  ```json
  { "action": "delete|categorize", "ids": ["uuid1", "uuid2"], "category_id": "uuid" }
  ```

---

## Categories & Groups

### GET /api/category-groups?budget_id=<uuid>
- **Auth:** Yes
- **Description:** List category groups for a budget

### POST /api/category-groups
- **Auth:** Yes
- **Description:** Create a category group

### PUT /api/category-groups/{id}
- **Auth:** Yes
- **Description:** Rename a category group

### DELETE /api/category-groups/{id}
- **Auth:** Yes
- **Description:** Delete a category group and its categories

### GET /api/categories?budget_id=<uuid>
- **Auth:** Yes
- **Description:** List non-hidden categories

### POST /api/categories
- **Auth:** Yes
- **Description:** Create a category

### PUT /api/categories/{id}
- **Auth:** Yes
- **Description:** Rename a category

### DELETE /api/categories/{id}
- **Auth:** Yes
- **Description:** Soft-delete a category (sets `is_hidden=True`)

---

## Monthly Budget

### GET /api/monthly-budgets?category_id=&month=&year=
- **Auth:** Yes
- **Description:** List monthly budget assignments

### PUT /api/monthly-budgets
- **Auth:** Yes
- **Description:** Assign funds to a category for a given month/year
- **Request body:** `{ "category_id": "uuid", "month": 6, "year": 2026, "amount": 500 }`

### POST /api/monthly-budgets/move
- **Auth:** Yes
- **Description:** Move budgeted funds between categories
- **Request body:** `{ "from_category": "uuid", "to_category": "uuid", "amount": 100, "month": 6, "year": 2026 }`

---

## Schedules

### GET /api/schedules?budget_id=<uuid>
- **Auth:** Yes
- **Description:** List schedules (programaciones)
- **Response 200:** `[{ "id": "uuid", "amount": 500, "frequency": "monthly", "next_date": "2026-07-01", "direction": "expense|income", ... }]`

### POST /api/schedules
- **Auth:** Yes
- **Description:** Create a recurring schedule
- **Request body:**
  ```json
  { "budget_id": "uuid", "account_id": "uuid", "amount": 500, "direction": "expense|income", "frequency": "weekly|biweekly|monthly|quarterly|yearly", "next_date": "2026-07-01", "payee_id": "uuid|null", "category_id": "uuid|null", "notes": "...", "skip_weekends": false, "apply_before_weekend": false }
  ```

### PUT /api/schedules/{id}
- **Auth:** Yes
- **Description:** Update schedule fields

### DELETE /api/schedules/{id}
- **Auth:** Yes
- **Description:** Delete a schedule

### POST /api/schedules/process?budget_id=<uuid>
- **Auth:** Yes
- **Description:** Manually process all due schedules (creates transactions, updates balances, advances next_date)
- **Response 200:** `{ "processed": 3 }`

---

## Credit Cards

### GET /api/credit-cards?budget_id=<uuid>
- **Auth:** Yes
- **Description:** List credit cards

### POST /api/credit-cards
- **Auth:** Yes
- **Description:** Create a credit card (also creates a "Pago {name}" category for TC auto-move)
- **Request body:**
  ```json
  { "budget_id": "uuid", "name": "Visa", "limit": 50000, "balance": -1000, "interest_rate": 0.96, "closing_day": 15, "due_day": 5, "notes": "" }
  ```

### GET /api/credit-cards/{id}
- **Auth:** Yes
- **Description:** Get credit card details

### PUT /api/credit-cards/{id}
- **Auth:** Yes
- **Description:** Update credit card

### DELETE /api/credit-cards/{id}
- **Auth:** Yes
- **Description:** Delete a credit card

### POST /api/credit-cards/{id}/pay
- **Auth:** Yes
- **Description:** Pay credit card (creates transaction, applies interest, updates balances)
- **Request body:** `{ "account_id": "uuid", "amount": 500 }`

---

## Loans

### GET /api/loans?budget_id=<uuid>
- **Auth:** Yes
- **Description:** List loans

### POST /api/loans
- **Auth:** Yes
- **Description:** Create a loan (generates installments automatically)
- **Request body:**
  ```json
  { "budget_id": "uuid", "account_id": "uuid", "type": "personal|hipotecario|automotor|estudiantil|otros", "name": "Car Loan", "total_amount": 300000, "interest_rate": 0.12, "total_installments": 36, "start_date": "2026-01-01", "next_due_date": "2026-02-01", "installment_amount": 9500, "notes": "" }
  ```

### GET /api/loans/{id}
- **Auth:** Yes
- **Description:** Get loan details with installments

### POST /api/loans/{id}/pay-installment
- **Auth:** Yes
- **Description:** Pay the next pending installment (creates transaction, updates balance, advances loan)
- **Request body:** `{ "account_id": "uuid" }`

---

## Payees

### GET /api/payees?budget_id=<uuid>
- **Auth:** Yes
- **Description:** List payees

### POST /api/payees
- **Auth:** Yes
- **Description:** Create a payee

### PUT /api/payees/{id}
- **Auth:** Yes
- **Description:** Rename a payee

### DELETE /api/payees/{id}
- **Auth:** Yes
- **Description:** Delete a payee

### POST /api/payees/{id}/merge
- **Auth:** Yes
- **Description:** Merge this payee into another (updates all references, deletes this one)
- **Request body:** `{ "target_id": "uuid" }`

---

## Goals

### GET /api/goals?budget_id=<uuid>
- **Auth:** Yes
- **Description:** List goals/targets for a budget

### POST /api/goals
- **Auth:** Yes
- **Description:** Create a goal/target (applies to monthly budget automatically)
- **Request body:**
  ```json
  { "category_id": "uuid", "goal_type": "monthly|yearly|target_balance|target_date|true_expense", "amount": 500, "target_date": "2026-12-31|null", "frequency": 12 }
  ```

### PUT /api/goals/{id}
- **Auth:** Yes
- **Description:** Update a goal

### DELETE /api/goals/{id}
- **Auth:** Yes
- **Description:** Delete a goal

---

## Rules

### GET /api/rules?budget_id=<uuid>
- **Auth:** Yes
- **Description:** List automation rules

### POST /api/rules
- **Auth:** Yes
- **Description:** Create a rule (auto-categorize transactions based on conditions)
- **Request body:**
  ```json
  { "budget_id": "uuid", "name": "Supermarket", "condition_field": "payee|amount|notes", "condition_operator": "contains|equals|starts_with|greater_than|less_than", "condition_value": "supermarket", "action_category_id": "uuid|null", "action_payee_id": "uuid|null", "action_notes": "" }
  ```

### PUT /api/rules/{id}
- **Auth:** Yes
- **Description:** Update a rule

### DELETE /api/rules/{id}
- **Auth:** Yes
- **Description:** Delete a rule

---

## Reports

### GET /api/reports/net-worth?budget_id=<uuid>
- **Auth:** Yes
- **Description:** Calculate net worth (assets - liabilities)
- **Response 200:** `{ "assets": 50000, "liabilities": 10000, "net_worth": 40000 }`

### GET /api/reports/cash-flow?budget_id=<uuid>
- **Auth:** Yes
- **Description:** Last 6 months income/expenses summary
- **Response 200:** `[{ "month": "6/2026", "income": 5000, "expenses": 3200, "net": 1800 }, ...]`

### GET /api/reports/budget-vs-reality?budget_id=<uuid>
- **Auth:** Yes
- **Description:** Budget vs actual spending per category for the current month

---

## YNAB Engine

### GET /api/budgets/{id}/ready-to-assign?month=&year=
- **Auth:** Yes
- **Description:** Calculate Ready to Assign = on_budget_balance - total_assigned
- **Response 200:** `{ "ready_to_assign": 2000, "total_on_budget": 5000, "total_assigned": 3000 }`

### POST /api/budgets/{id}/auto-assign
- **Auth:** Yes
- **Description:** Auto-distribute Ready to Assign across underfunded categories (priority: overspends > true_expense > monthly targets)
- **Request body:** `{ "month": 6, "year": 2026 }`
- **Response 200:** `{ "assigned_categories": 3, "total_assigned": 900, "remaining": 100 }`

### POST /api/budgets/{id}/cover
- **Auth:** Yes
- **Description:** Move funds from one category to cover an overspent category
- **Request body:** `{ "from_category": "uuid", "to_category": "uuid", "amount": 50, "month": 6, "year": 2026 }`

### GET /api/budgets/{id}/spotlight?month=&year=
- **Auth:** Yes
- **Description:** Aggregated alerts: uncategorized transactions, overspends, underfunded targets
- **Response 200:** `{ "uncategorized_count": 3, "overspends": [...], "underfunded": [...], "total_alerts": 2 }`

### GET /api/budgets/{id}/cost-to-be-me
- **Auth:** Yes
- **Description:** Total monthly target cost vs expected monthly income (3-month avg)
- **Response 200:** `{ "cost_to_be_me": 3500, "expected_income": 4200, "difference": 700, "is_over_budget": false }`

### GET /api/budgets/{id}/rollover?from_month=&from_year=
- **Auth:** Yes
- **Description:** Monthly rollover balances per category (available = budgeted - spent)
- **Response 200:** `{ "categories": [{ "category_id": "uuid", "name": "Groceries", "rollover_balance": 40 }] }`

### GET /api/categories/{id}/inspector?month=&year=
- **Auth:** Yes
- **Description:** Category detail for inspector panel: balance, assigned, activity, target info, 3-month average
- **Response 200:** Full category detail with target and overspent status

---

## Settings

### GET /api/settings/profile
- **Auth:** Yes
- **Description:** Get user profile

### PUT /api/settings/profile
- **Auth:** Yes
- **Description:** Update user profile (first_name, last_name, email)

### GET /api/settings/budget/{id}
- **Auth:** Yes
- **Description:** Get budget settings with members

### POST /api/settings/budget/{id}/invite
- **Auth:** Yes
- **Description:** Invite a user to the budget by email
- **Request body:** `{ "email": "user@email.com" }`

---

## Sync

### GET /api/sync/logs?budget_id=<uuid>
- **Auth:** Yes
- **Description:** Get sync logs for a budget (last 100)
- **Response 200:** `[{ "id": "uuid", "entity_type": "...", "action": "pending|synced", ... }]`

---

## Response Format

All endpoints return JSON. Errors follow:
```json
{ "error": "description" }
```

HTTP status codes:
- `200` OK
- `201` Created
- `400` Bad Request (invalid input)
- `401` Unauthorized (missing/invalid token)
- `404` Not Found
- `409` Conflict (duplicate)
- `500` Internal Server Error

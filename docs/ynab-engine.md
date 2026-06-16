# YNAB Engine

## Regla 1: Darle trabajo a cada peso

Implemented in `apps/budgets/views.py` — `assign_funds` view.

- User allocates a category budget for a given month/year
- `MonthlyBudget` tracks `budgeted` per category per month
- Available-to-budget: sum of all on-budget account balances minus total budgeted across categories
- Zero-sum enforcement: cannot budget more than available

## Regla 2: Aceptar tus gastos reales

### True Expenses

Implemented in `apps/goals/views.py` — `_apply_goal`.

- Goal type `true_expense`: annual expense split into 12 monthly installments
- `mb.budgeted = amount / frequency` (frequency = months)
- Category goals types: `monthly`, `target_balance`, `target_date`, `true_expense`

### Credit Card Interest

`core/interest.py` — `calcular_interes_diario(saldo, tasa_anual, dias)`.

- Interest accrues daily on unpaid credit card balance after due date
- `tasa_diaria = tasa_anual / 12 / 30`
- Recorded as automatic transaction

## Regla 3: Patear los golpes

### Move Money

`apps/budgets/views.py` — `move_funds` view.

- Transfer budgeted amount from one category to another within the same month
- Both `MonthlyBudget` records updated atomically

### Overspending Rollover

Calculated in budget context processor (`chele/context_processors.py`).

- If a category spent more than budgeted, the overspent amount carries to next month
- Subtracted from available-to-budget in the new month

## Regla 4: Envejecer tu dinero

### Hold for Next Month

`apps/budgets/views.py` — `toggle_hold` view.

- Income assigned to categories marked for next month
- `MonthlyBudget.for_next_month` flag
- Amount held is excluded from current month's available-to-budget

### Auto-hold on Income Categories

Income category groups (`is_income=True`) auto-hold assigned funds for the next month.

## Key Models

| Model | Key Fields | Role |
|-------|-----------|------|
| `CategoryGroup` | `name`, `is_income` | Groups categories, income flag |
| `Category` | `name`, `group`, `budget` | Individual budget line |
| `MonthlyBudget` | `category`, `month`, `year`, `budgeted`, `for_next_month` | Monthly allocation per category |
| `Goal` | `category`, `goal_type`, `amount`, `frequency` | Auto-budget targets |
| `Account` | `balance`, `on_budget` | Funding source |
| `Transaction` | `account`, `category`, `amount`, `date` | Spending/income records |

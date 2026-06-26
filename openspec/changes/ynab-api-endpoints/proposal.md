## Why

The Go API (`chele-api/`) currently covers CRUD for accounts, budgets, transactions, schedules, credit cards, and loans, but lacks the YNAB engine endpoints needed by external clients and future frontends. Without these, the API cannot replicate the full functionality of the Django app: targets, ready-to-assign, auto-assign, cover overspending, spotlight alerts, inspector data, and cost-to-be-me. This change adds those missing endpoints so the Go API is feature-complete with Django.

## What Changes

- **Targets CRUD**: New endpoints to create, read, update, delete targets with all goal types (monthly, yearly, target_balance, target_date, true_expense) plus snooze and refill_up_to.
- **Ready to Assign**: Endpoint that returns the current available-to-assign balance calculated from on-budget accounts minus assigned funds.
- **Auto-assign**: Endpoint that distributes Ready to Assign across underfunded categories.
- **Cover overspending**: Endpoint to transfer funds from one category to cover an overspent category.
- **Spotlight alerts**: Endpoint that aggregates uncategorized transactions, uncovered overspends, and underfunded targets.
- **Inspector data**: Endpoint returning detailed category info (target, averages, activity).
- **Cost to be me**: Endpoint returning total monthly target cost vs expected income.
- **Rollover**: Endpoint returning next month starting balances per category.

## Capabilities

### New Capabilities
- `api-targets`: CRUD for targets with all goal types, snooze, refill_up_to
- `api-ready-to-assign`: Endpoint to calculate available funds
- `api-auto-assign`: Endpoint to distribute funds to underfunded categories
- `api-cover-overspending`: Endpoint to cover overspent categories
- `api-spotlight`: Endpoint returning pending alerts
- `api-inspector`: Endpoint returning category detail for inspector panel
- `api-cost-to-be-me`: Endpoint returning target cost vs income comparison
- `api-rollover`: Endpoint returning next month category starting balances

### Modified Capabilities
- *(none — all endpoints are new)*

## Impact

- **chele-api/internal/models/**: New Target model struct
- **chele-api/internal/handlers/**: New handler files for targets, ready-to-assign, auto-assign, cover, spotlight, inspector, cost-to-be-me, rollover
- **chele-api/internal/service/**: New services for target calculation, overspending detection, rollover computation, spotlight aggregation
- **chele-api/internal/router/**: New routes for all endpoints
- **chele-api/internal/testutil/**: Target table in test schema

## ADDED Requirements

### Requirement: Rollover endpoint
The API SHALL provide an endpoint that calculates starting available balances for the next month for each category. Rollover SHALL be calculated as: current_month_available + budgeted - spent. Negative values (cash overspend) SHALL be carried as negative. Categories with refill_up_to targets SHALL show the underfunded amount accounting for rollover.

#### Scenario: Get next month rollover
- **WHEN** client sends GET `/api/budgets/<id>/rollover?from_month=5&from_year=2026&to_month=6&to_year=2026`
- **THEN** response SHALL be 200 with:
```json
{
  "categories": [
    {"category_id":"<uuid>","name":"Groceries","rollover_balance":40.00},
    {"category_id":"<uuid>","name":"Eating Out","rollover_balance":-15.00}
  ]
}
```

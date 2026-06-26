## ADDED Requirements

### Requirement: Ready to Assign endpoint
The API SHALL provide an endpoint to calculate and return the current Ready to Assign balance. The calculation SHALL be: sum(all on-budget account balances) - sum(all MonthlyBudget.budgeted for current month and year).

#### Scenario: Get Ready to Assign
- **WHEN** client sends GET `/api/budgets/<id>/ready-to-assign?month=6&year=2026`
- **THEN** response SHALL be 200 with `{"ready_to_assign": 2000.00, "total_on_budget": 5000.00, "total_assigned": 3000.00}`

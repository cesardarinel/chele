## ADDED Requirements

### Requirement: Cover overspending endpoint
The API SHALL provide an endpoint to cover an overspent category by transferring funds from another category. The transfer SHALL be atomic.

#### Scenario: Cover overspend
- **WHEN** client sends POST `/api/budgets/<id>/cover` with `{"from_category_id":"<uuid>","to_category_id":"<uuid>","amount":50.00,"month":6,"year":2026}`
- **THEN** response SHALL be 200 with `{"status":"covered"}` and both category balances SHALL be updated

#### Scenario: Cover with insufficient source
- **WHEN** source category has $30 available but client requests $50
- **THEN** response SHALL be 400 with `{"error":"Fondos insuficientes en la categoría origen"}`

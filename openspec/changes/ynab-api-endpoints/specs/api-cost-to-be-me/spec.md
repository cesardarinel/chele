## ADDED Requirements

### Requirement: Cost to be me endpoint
The API SHALL provide an endpoint that calculates and returns the total monthly cost of all active targets (Cost to be me) and compares it against the expected monthly income (average of last 3 months income transactions).

#### Scenario: Get cost to be me
- **WHEN** client sends GET `/api/budgets/<id>/cost-to-be-me`
- **THEN** response SHALL be 200 with:
```json
{
  "cost_to_be_me": 3500.00,
  "expected_income": 4200.00,
  "difference": 700.00,
  "is_over_budget": false
}
```

#### Scenario: Cost exceeds income (reality check)
- **WHEN** cost_to_be_me ($4,000) exceeds expected_income ($3,200)
- **THEN** response SHALL include `"is_over_budget": true` and `"difference": -800.00`

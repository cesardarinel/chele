## ADDED Requirements

### Requirement: Spotlight alerts endpoint
The API SHALL provide an endpoint that aggregates all pending actions for a budget into a single response. The response SHALL include: uncategorized transactions count, uncovered overspends list, underfunded targets list with amounts.

#### Scenario: Get spotlight alerts
- **WHEN** client sends GET `/api/budgets/<id>/spotlight?month=6&year=2026`
- **THEN** response SHALL be 200 with:
```json
{
  "uncategorized_count": 3,
  "overspends": [{"category_id":"<uuid>","category_name":"Groceries","amount":-15.00}],
  "underfunded": [{"category_id":"<uuid>","category_name":"Car Insurance","deficit":300.00}],
  "total_alerts": 2
}
```

#### Scenario: No alerts
- **WHEN** there are no pending actions
- **THEN** response SHALL be 200 with `{"uncategorized_count":0,"overspends":[],"underfunded":[],"total_alerts":0}`

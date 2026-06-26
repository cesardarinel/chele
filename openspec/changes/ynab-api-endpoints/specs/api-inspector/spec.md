## ADDED Requirements

### Requirement: Category inspector endpoint
The API SHALL provide an endpoint returning detailed information about a category for the inspector panel. The response SHALL include: available balance, total assigned this month, total activity, target info (if any), underfunded amount, 3-month average spending, overspent status.

#### Scenario: Get category inspector data
- **WHEN** client sends GET `/api/categories/<id>/inspector?month=6&year=2026`
- **THEN** response SHALL be 200 with:
```json
{
  "category_id": "<uuid>",
  "category_name": "Groceries",
  "available": 150.00,
  "assigned": 200.00,
  "activity": -50.00,
  "target": {"type": "monthly", "amount": 200.00, "refill_up_to": true, "underfunded": 50.00},
  "average_spending_3m": 180.00,
  "overspent": {"is_overspent": false, "amount": 0, "type": ""}
}
```

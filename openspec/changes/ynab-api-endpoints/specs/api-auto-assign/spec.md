## ADDED Requirements

### Requirement: Auto-assign endpoint
The API SHALL provide an endpoint that automatically distributes Ready to Assign funds across all underfunded categories. The distribution order SHALL be: uncovered cash overspends first, then true_expense targets, then monthly targets, then other target types. If Ready to Assign is insufficient, SHALL fund in priority order until exhausted.

#### Scenario: Auto-assign all
- **WHEN** client sends POST `/api/budgets/<id>/auto-assign` with `{"month":6,"year":2026}`
- **THEN** response SHALL be 200 with `{"assigned_categories": 3, "total_assigned": 900.00, "remaining": 100.00}`

#### Scenario: Auto-assign with zero ready
- **WHEN** Ready to Assign is $0 and client sends POST `/api/budgets/<id>/auto-assign`
- **THEN** response SHALL be 200 with `{"assigned_categories": 0, "total_assigned": 0, "remaining": 0}` and no categories SHALL be modified

## ADDED Requirements

### Requirement: Target CRUD
The API SHALL provide full CRUD for targets. Each target SHALL support fields: category_id, goal_type (monthly|yearly|target_balance|target_date|true_expense), amount, target_date (nullable), frequency (int, default 12), refill_up_to (bool, default false), snooze_month (int, nullable), snooze_year (int, nullable), is_completed (bool).

#### Scenario: Create monthly target
- **WHEN** client sends POST `/api/targets` with `{"category_id":"<uuid>","goal_type":"monthly","amount":500}`
- **THEN** response SHALL be 201 with the created target id

#### Scenario: List targets for a category
- **WHEN** client sends GET `/api/targets?category_id=<uuid>`
- **THEN** response SHALL be 200 with a JSON array of targets

#### Scenario: Update target amount
- **WHEN** client sends PUT `/api/targets/<id>` with `{"amount":600}`
- **THEN** response SHALL be 200 and the target amount SHALL be updated

#### Scenario: Delete target
- **WHEN** client sends DELETE `/api/targets/<id>`
- **THEN** response SHALL be 200 and the target SHALL be removed

#### Scenario: Snooze a target
- **WHEN** client sends PUT `/api/targets/<id>` with `{"snooze_month":6,"snooze_year":2026}`
- **THEN** the target SHALL be excluded from underfunded calculations for that month

#### Scenario: Create yearly target
- **WHEN** client sends POST `/api/targets` with `{"category_id":"<uuid>","goal_type":"yearly","amount":1200,"target_date":"2026-12-31"}`
- **THEN** the system SHALL calculate monthly requirement as $100 ($1200/12)

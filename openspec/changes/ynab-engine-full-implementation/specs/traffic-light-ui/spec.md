## ADDED Requirements

### Requirement: Traffic light color system
The system SHALL implement a consistent traffic light color system across all budget views:
- **GREEN**: Category is fully funded (assigned ≥ target requirement or balance ≥ 0 and no issues)
- **YELLOW/AMBER**: Category is underfunded (target exists but not fully assigned) or has a credit card overspend
- **RED**: Category has a cash overspend (negative available balance requiring immediate correction)

These colors SHALL be applied to the category's available balance, indicator dot, and inspector panel.

#### Scenario: Category fully funded shows green
- **WHEN** a category has no target and positive available balance
- **THEN** the available balance SHALL be displayed in GREEN

#### Scenario: Underfunded shows yellow
- **WHEN** a category has a target and assigned amount is less than required
- **THEN** the category SHALL display a YELLOW indicator and the deficit amount

#### Scenario: Cash overspend shows red
- **WHEN** a category's available balance is negative due to cash spending
- **THEN** the available balance SHALL be displayed in RED with a warning icon

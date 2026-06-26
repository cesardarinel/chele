## ADDED Requirements

### Requirement: Snooze target
The system SHALL allow users to snooze (pause) a target for a specific month. When snoozed, the target SHALL NOT generate underfunded warnings or affect auto-assign calculations for that month. The snooze SHALL be limited to the current month and SHALL automatically expire at month end.

#### Scenario: Snooze a target
- **WHEN** user activates snooze on a monthly $500 target for June
- **THEN** the underfunded indicator for June SHALL be hidden, and auto-assign SHALL skip this category

#### Scenario: Snooze auto-expires
- **WHEN** a target was snoozed for June and the calendar advances to July
- **THEN** the target SHALL resume normal underfunded calculation for July

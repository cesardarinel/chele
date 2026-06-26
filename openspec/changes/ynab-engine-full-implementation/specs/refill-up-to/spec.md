## ADDED Requirements

### Requirement: Refill up to mode
The system SHALL support a "Refill up to" target mode. In this mode, the target SHALL calculate the needed amount as: max(0, target_amount - current_category_balance). The current balance includes rollover from previous months. This mode is ideal for expenses with variable balances that should be topped up to a ceiling each month.

#### Scenario: Refill up to with rollover
- **WHEN** user has a "Refill up to $200" target on "Groceries" and the category has $50 remaining from last month
- **THEN** the underfunded amount SHALL be $150 ($200 - $50)

#### Scenario: Refill up to fully funded
- **WHEN** user has a "Refill up to $200" target and the category balance is $200 or more
- **THEN** the underfunded amount SHALL be $0 and the indicator SHALL be green

### Requirement: Set aside another mode
The system SHALL support a "Set aside another" target mode. In this mode, the system SHALL request the full target amount each month regardless of rollover balance. Underfunded = target_amount (always asks for the full amount).

#### Scenario: Set aside another ignores rollover
- **WHEN** user has a "Set aside another $200" target on "Groceries" and the category has $50 from last month
- **THEN** the underfunded amount SHALL be $200 (ignores the $50 rollover)

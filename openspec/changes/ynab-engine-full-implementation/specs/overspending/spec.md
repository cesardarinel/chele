## ADDED Requirements

### Requirement: Overspending detection (cash)
The system SHALL detect when a cash/on-budget category has spent more than its available balance. When detected, the category SHALL display a RED indicator showing the overspent amount. The overspent amount SHALL be carried forward to the next month as a negative available balance that MUST be covered.

#### Scenario: Cash overspending turns red
- **WHEN** user has $25 available in "Groceries" and records a $40 expense
- **THEN** the category SHALL display -$15 in RED

#### Scenario: Cash overspend carries to next month
- **WHEN** a category had -$15 overspent at month end and the next month begins
- **THEN** the available balance for that category SHALL start at -$15

### Requirement: Overspending detection (credit card)
The system SHALL detect when a credit card category has overspent. When detected with a credit card transaction, the system SHALL automatically move the available funds (up to the overspent amount) from the spending category to the credit card's Payment category. The overspent credit card category SHALL display an ORANGE indicator.

#### Scenario: Credit card overspend auto-moves funds
- **WHEN** user has $25 available in "Groceries" and records a $40 expense on Visa
- **THEN** the system SHALL move $25 from "Groceries" available balance to the Visa Payment category

#### Scenario: Credit card overspend shows orange
- **WHEN** a credit card transaction causes overspending
- **THEN** the category SHALL display the overspent amount in ORANGE (not red)

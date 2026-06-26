## ADDED Requirements

### Requirement: Automatic rollover of available balances
At the start of each month, the system SHALL automatically carry forward the available balance (budgeted - spent) of each category to the new month. This becomes the starting available balance for that category. Categories with "Refill up to" targets SHALL use this rollover balance to calculate the underfunded amount.

#### Scenario: Positive rollover
- **WHEN** a category had $100 budgeted and only $60 spent last month
- **THEN** $40 SHALL roll over to the next month as starting available balance

#### Scenario: Negative rollover (cash overspend)
- **WHEN** a category had $100 budgeted and $130 spent last month (cash)
- **THEN** -$30 SHALL roll over as negative starting balance (red indicator)

#### Scenario: Rollover with Refill up to target
- **WHEN** a category has a "Refill up to $200" target and $50 rolled over from last month
- **THEN** the underfunded amount SHALL be $150 ($200 - $50)

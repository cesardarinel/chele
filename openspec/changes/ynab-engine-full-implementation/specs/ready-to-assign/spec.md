## ADDED Requirements

### Requirement: Ready to Assign
The system SHALL display a "Ready to Assign" section at the top of the budget view showing the total funds available for distribution across all categories. The value SHALL be calculated as: sum of all on-budget account balances minus sum of all MonthlyBudget.budgeted for the current month. The system SHALL enforce zero-sum budgeting: users MUST NOT assign more than the Ready to Assign amount.

#### Scenario: Ready to Assign shows correct balance
- **WHEN** user has $5,000 in on-budget accounts and has assigned $3,000 to categories this month
- **THEN** the Ready to Assign value SHALL be $2,000

#### Scenario: Cannot assign more than available
- **WHEN** user attempts to assign $3,000 to a category but only $2,000 is available in Ready to Assign
- **THEN** the system SHALL reject the assignment with an error message "No tienes suficientes fondos disponibles"

#### Scenario: Assigning reduces Ready to Assign
- **WHEN** user assigns $500 from Ready to Assign to a category
- **THEN** the Ready to Assign value SHALL decrease by $500 immediately

#### Scenario: Ready to Assign with new account
- **WHEN** user adds a new on-budget account with balance $3,140
- **THEN** the Ready to Assign value SHALL increase by $3,140

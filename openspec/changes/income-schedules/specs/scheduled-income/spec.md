## ADDED Requirements

### Requirement: Schedule direction
The system SHALL support both income and expense directions for scheduled transactions.

#### Scenario: Create income schedule
- **WHEN** user creates a schedule with direction `income`
- **THEN** the schedule is saved with `direction = 'income'`

#### Scenario: Create expense schedule
- **WHEN** user creates a schedule with direction `expense`
- **THEN** the schedule is saved with `direction = 'expense'`

#### Scenario: Default direction is expense
- **WHEN** user creates a schedule without specifying direction
- **THEN** the schedule is saved with `direction = 'expense'`

### Requirement: Income schedule execution
The system SHALL create positive-amount transactions and increase account balance when processing due income schedules.

#### Scenario: Income creates positive transaction
- **WHEN** a due income schedule is processed
- **THEN** a transaction is created with `amount > 0`

#### Scenario: Income increases account balance
- **WHEN** a due income schedule is processed
- **THEN** the account balance is increased by the scheduled amount

#### Scenario: Expense creates negative transaction
- **WHEN** a due expense schedule is processed
- **THEN** a transaction is created with `amount < 0`

#### Scenario: Expense decreases account balance
- **WHEN** a due expense schedule is processed
- **THEN** the account balance is decreased by the scheduled amount

### Requirement: Schedule direction visibility
The system SHALL display the direction clearly in the schedule list.

#### Scenario: Income shown in green with plus sign
- **WHEN** viewing the schedule list
- **THEN** income schedules display the amount in green with a `+` prefix

#### Scenario: Expense shown in red with minus sign
- **WHEN** viewing the schedule list
- **THEN** expense schedules display the amount in red with a `-` prefix

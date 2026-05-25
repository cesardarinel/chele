## ADDED Requirements

### Requirement: Monthly Savings Goal
The system SHALL allow setting a monthly savings goal on a category.

#### Scenario: Set monthly goal
- **WHEN** a user sets a Monthly Savings Goal of $200 on a category
- **THEN** the system auto-fills the budgeted amount to $200 each month

#### Scenario: Goal progress
- **WHEN** a user views a category with a monthly goal
- **THEN** the system shows progress toward the goal

### Requirement: Target Balance
The system SHALL allow setting a target balance goal on a category.

#### Scenario: Set target balance
- **WHEN** a user sets a Target Balance of $1,000 on a category with current balance $300
- **THEN** the system calculates the needed amount ($700) and suggests or auto-fills the budgeted amount

#### Scenario: Goal reached
- **WHEN** a category's balance reaches or exceeds the target
- **THEN** the system marks the goal as complete and stops suggesting allocations

### Requirement: Target by Date
The system SHALL allow setting a target amount to reach by a specific date.

#### Scenario: Set target by date
- **WHEN** a user sets a Target by Date of $1,200 for December 2026 on a category with $0 balance
- **THEN** the system calculates the monthly contribution ($1,200 / months remaining) and suggests it

### Requirement: True Expense
The system SHALL allow setting up True Expenses (annual/quarterly expenses split into monthly payments).

#### Scenario: Set true expense
- **WHEN** a user sets a True Expense of $1,200/year for car insurance
- **THEN** the system calculates $100/month to assign and auto-fills the budgeted amount each month

#### Scenario: True expense spending
- **WHEN** the annual expense occurs and is recorded against the category
- **THEN** the system deducts from the accumulated balance in that category

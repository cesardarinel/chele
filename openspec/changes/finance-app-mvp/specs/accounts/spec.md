## ADDED Requirements

### Requirement: Account types
The system SHALL support multiple account types: checking, savings, cash, and credit card.

#### Scenario: Create account
- **WHEN** a user creates a new account with a type and opening balance
- **THEN** the system creates the account and records the opening balance as an initial transaction

#### Scenario: Account types display
- **WHEN** viewing the account list
- **THEN** the system shows each account with its type and current balance

### Requirement: Balance calculation
The system SHALL calculate and display the current balance for each account based on all transactions.

#### Scenario: Balance update
- **WHEN** a transaction is added, edited, or deleted
- **THEN** the system recalculates the account balance accordingly

### Requirement: Net worth calculation
The system SHALL calculate net worth as assets (checking, savings, cash) minus liabilities (credit cards).

#### Scenario: Net worth display
- **WHEN** a user views the dashboard
- **THEN** the system displays the current net worth calculated from all accounts

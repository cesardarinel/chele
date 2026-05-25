## ADDED Requirements

### Requirement: Record transaction
The system SHALL allow users to record income and expense transactions.

#### Scenario: Record an expense
- **WHEN** a user enters an expense with amount, date, account, category, and optional payee
- **THEN** the system creates the transaction and updates the account balance and category spending

#### Scenario: Record income
- **WHEN** a user enters an income transaction with amount, date, and account
- **THEN** the system creates the transaction and updates the account balance and available-to-budget

### Requirement: Edit transaction
The system SHALL allow users to edit existing transactions.

#### Scenario: Edit amount
- **WHEN** a user changes the amount of an existing transaction
- **THEN** the system updates the transaction and recalculates affected balances and budget amounts

### Requirement: Delete transaction
The system SHALL allow users to delete transactions.

#### Scenario: Delete a transaction
- **WHEN** a user deletes a transaction
- **THEN** the system removes the transaction and reverses its effect on balances and budget

### Requirement: Transaction list
The system SHALL display transactions for an account or category with pagination.

#### Scenario: View transactions
- **WHEN** a user views an account or category
- **THEN** the system shows a paginated list of transactions sorted by date (newest first)

### Requirement: Transfer between accounts
The system SHALL support transfers between accounts.

#### Scenario: Record a transfer
- **WHEN** a user records a transfer from one account to another
- **THEN** the system creates a transfer pair that debits one account and credits the other

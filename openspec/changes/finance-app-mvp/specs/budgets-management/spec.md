## ADDED Requirements

### Requirement: Create multiple budgets
The system SHALL allow users to create multiple independent budgets.

#### Scenario: Create a new budget
- **WHEN** a user creates a new budget with a name and description
- **THEN** the system creates an empty budget with no accounts, categories, or transactions

#### Scenario: Budget independence
- **WHEN** a user creates two budgets
- **THEN** each budget has its own accounts, categories, transactions, members, payees, rules, and schedules

### Requirement: Switch between budgets
The system SHALL allow switching between budgets from the sidebar.

#### Scenario: Budget selector
- **WHEN** a user clicks the budget selector in the sidebar
- **THEN** the system shows a dropdown with all user's budgets and an option to create a new one

#### Scenario: Switch active budget
- **WHEN** a user selects a different budget from the dropdown
- **THEN** the system reloads all data for the selected budget (accounts, categories, transactions)

### Requirement: Budget members
The system SHALL allow inviting other users to a budget.

#### Scenario: Invite to budget
- **WHEN** a budget owner invites a user by email
- **THEN** the invited user gains access to that budget's accounts, categories, and transactions

#### Scenario: Budget-specific membership
- **WHEN** a user is a member of Budget A but not Budget B
- **THEN** they can only see and edit Budget A's data

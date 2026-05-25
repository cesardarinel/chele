## ADDED Requirements

### Requirement: Credit card as account type
The system SHALL treat credit cards as a distinct account type with special behavior.

#### Scenario: Create credit card account
- **WHEN** a user creates a credit card account
- **THEN** the system marks it as a credit card type with a negative balance (owed amount)

### Requirement: Credit card expense handling
When recording an expense paid with a credit card, the system SHALL move the transaction amount from the spending category to a "Credit Card Payment" category.

#### Scenario: Expense on credit card
- **WHEN** a user records an expense on a credit card and assigns it to a category
- **THEN** the system deducts the amount from that category's available funds and adds it to a "Credit Card Payment" category

### Requirement: Credit card payment
The system SHALL allow users to record payments to the credit card from other accounts.

#### Scenario: Pay credit card bill
- **WHEN** a user records a transfer from a checking account to a credit card
- **THEN** the system reduces the credit card balance (owed amount) and clears the corresponding amount from the "Credit Card Payment" category

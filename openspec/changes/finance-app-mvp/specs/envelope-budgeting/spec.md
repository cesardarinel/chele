## ADDED Requirements

### Requirement: Monthly budget periods
The system SHALL organize budgets by month.

#### Scenario: Default budget month
- **WHEN** a user accesses the budget page
- **THEN** the system shows the current month's budget by default

#### Scenario: Navigate between months
- **WHEN** a user navigates to a different month
- **THEN** the system displays that month's budget allocations and transactions

### Requirement: Budget categories
The system SHALL allow users to create categories for envelope budgeting.

#### Scenario: Create category
- **WHEN** a user creates a category with a name
- **THEN** the system adds the category to the budget

#### Scenario: Delete category
- **WHEN** a user deletes a category
- **THEN** the system removes the category and any remaining funds are returned to the available pool

### Requirement: Fund allocation
The system SHALL allow users to allocate money to categories each month.

#### Scenario: Allocate funds
- **WHEN** a user assigns an amount to a category for a given month
- **THEN** the system deducts from the available-to-budget amount and records the allocation

#### Scenario: Over-allocation prevention
- **WHEN** a user tries to allocate more than the available-to-budget amount
- **THEN** the system shows an error and prevents the allocation

### Requirement: Available-to-budget
The system SHALL calculate the available-to-budget amount as income minus allocations. The goal is to reach zero (Regla 1 YNAB).

#### Scenario: Available calculation
- **WHEN** income transactions are added or allocations change
- **THEN** the system recalculates and displays the available-to-budget amount

#### Scenario: Zero-sum budgeting
- **WHEN** the available-to-budget amount reaches zero
- **THEN** all income has been assigned to categories

### Requirement: Move money between categories (Regla 3 YNAB)
The system SHALL allow moving allocated funds between categories.

#### Scenario: Transfer between categories
- **WHEN** a user transfers an amount from one category to another
- **THEN** the system deducts from the source category and adds to the target category

#### Scenario: Transfer to available-to-budget
- **WHEN** a user transfers from a category back to available-to-budget
- **THEN** the funds are returned to the pool for re-allocation

### Requirement: Overspending rollover (Regla 3 YNAB)
The system SHALL automatically deduct overspending from next month's available-to-budget.

#### Scenario: Overspending rollover
- **WHEN** a category has a negative balance at month end
- **THEN** the system deducts the overspent amount from next month's available-to-budget

### Requirement: Hold for next month (Regla 4 YNAB)
The system SHALL allow users to hold available funds for the next month.

#### Scenario: Hold funds
- **WHEN** a user clicks "A presupuestar" and selects "Reservar para el próximo mes"
- **THEN** the system removes that amount from current month's available-to-budget and adds it to next month's

#### Scenario: Auto-hold
- **WHEN** a user enables auto-hold on an income category
- **THEN** the system automatically holds that category's income for the next month for 12 months

### Requirement: Copy budget from previous month
The system SHALL allow copying budget allocations from the previous month.

#### Scenario: Copy last month
- **WHEN** a user selects "Copiar presupuesto del mes anterior"
- **THEN** the system copies all budgeted amounts from the previous month

#### Scenario: Set to averages
- **WHEN** a user selects promedio 3/6/12 meses
- **THEN** the system calculates and sets budgeted amounts based on average spending for that period

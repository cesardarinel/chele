## ADDED Requirements

### Requirement: Budget vs Reality report
The system SHALL show a comparison between budgeted amounts and actual spending per category for a given month.

#### Scenario: View Budget vs Reality
- **WHEN** a user views the Budget vs Reality report for a month
- **THEN** the system displays each category with budgeted amount, actual spending, and the difference

#### Scenario: Monthly comparison
- **WHEN** a user selects a different month
- **THEN** the system updates the report to show data for that month

### Requirement: Net Worth report
The system SHALL show net worth over time, calculated as assets minus liabilities.

#### Scenario: Net worth graph
- **WHEN** a user views the Net Worth report
- **THEN** the system displays a chart showing net worth over time

#### Scenario: Net worth breakdown
- **WHEN** a user views the Net Worth report
- **THEN** the system shows total assets, total liabilities, and net worth

### Requirement: Cash Flow report
The system SHALL show cash flow over time, tracking income vs expenses from budget accounts.

#### Scenario: Cash flow chart
- **WHEN** a user views the Cash Flow report
- **THEN** the system displays income and expenses over time with separate visualizations

#### Scenario: Cash flow summary
- **WHEN** a user views the Cash Flow report
- **THEN** the system shows total income, total expenses, and net cash flow for the selected period

### Requirement: Category breakdown
The system SHALL show spending breakdown per category with amounts and percentages.

#### Scenario: Category breakdown
- **WHEN** a user views the reports page
- **THEN** the system shows each category with total spent and percentage of total expenses

### Requirement: Income vs Expense summary
The system SHALL show total income, total expenses, and net for a selected period.

#### Scenario: Period summary
- **WHEN** a user selects a date range
- **THEN** the system displays total income, total expenses, and net difference for that period

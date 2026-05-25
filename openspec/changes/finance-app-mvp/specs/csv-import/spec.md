## ADDED Requirements

### Requirement: CSV import (onboarding)
The system SHALL allow users to import transactions from a CSV file during initial setup.

#### Scenario: Import CSV on first use
- **WHEN** a user uploads a CSV file with transactions during onboarding
- **THEN** the system parses the CSV, matches columns, and creates transactions for the selected account

#### Scenario: Column mapping
- **WHEN** the system detects CSV columns
- **THEN** the user can map CSV columns to system fields (date, amount, payee, category, notes)

### Requirement: CSV import (recurring)
The system SHALL allow users to import CSV files for ongoing transaction reconciliation.

#### Scenario: Import CSV for existing account
- **WHEN** a user uploads a CSV file for an existing account
- **THEN** the system parses the file and shows a preview before confirming import

#### Scenario: Duplicate detection
- **WHEN** importing transactions that already exist in the system
- **THEN** the system detects duplicates by matching amount, date, and payee, and skips them

### Requirement: CSV format support
The system SHALL support standard CSV formats from major banks with configurable date and number formats.

#### Scenario: Date format detection
- **WHEN** the CSV contains dates in a supported format (YYYY-MM-DD, DD/MM/YYYY, MM/DD/YYYY)
- **THEN** the system parses them correctly or asks the user to specify the format

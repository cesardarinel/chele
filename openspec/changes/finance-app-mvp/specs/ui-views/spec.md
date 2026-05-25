## ADDED Requirements

### Requirement: Sidebar navigation
The system SHALL display a persistent left sidebar with navigation to all views, matching ActualBudget's layout.

#### Scenario: Sidebar structure
- **WHEN** a user is logged in
- **THEN** the sidebar shows: Presupuesto, Reportes, Programaciones, Más (Beneficiarios, Reglas), Cuentas, Configuración

#### Scenario: Account list in sidebar
- **WHEN** a user views the sidebar
- **THEN** each account is listed under "Cuentas" with its current balance

### Requirement: Budget view
The system SHALL display the budget view as a table with columns: Budgeted, Spent, Balance per category group.

#### Scenario: Budget table
- **WHEN** a user views the budget
- **THEN** the system shows categories grouped by category group, with Budgeted, Spent, and Balance columns

#### Scenario: Month header
- **WHEN** a user views a budget month
- **THEN** the header shows: note icon, chevrons to minimize, 3-dot menu (copiar presupuesto del mes anterior, establecer a cero, promedio 3/6/12 meses)

#### Scenario: Category group actions
- **WHEN** a user hovers over a category group
- **THEN** the system shows options to add category, add note, hide/show, rename, delete group

#### Scenario: Available-to-budget display
- **WHEN** a user views the budget
- **THEN** the system shows "A presupuestar" amount at the top of the budget section

### Requirement: Account Register view
The system SHALL display the account register with a transaction list and action buttons.

#### Scenario: Transaction list
- **WHEN** a user views an account
- **THEN** the system shows transactions with: checkbox, date, beneficiary, category, notes, amount, running balance

#### Scenario: Account actions
- **WHEN** a user views an account
- **THEN** the system shows buttons for Importar, Añadir Nuevo, Filtrar, Buscar

#### Scenario: Bulk actions
- **WHEN** a user selects multiple transactions via checkboxes
- **THEN** the system shows bulk action options (delete, edit, move category, etc.)

### Requirement: Reports Dashboard
The system SHALL display reports as clickable tiles that expand to full screen.

#### Scenario: Report tiles
- **WHEN** a user views the reports page
- **THEN** the system shows report tiles for Net Worth, Cash Flow, and Budget vs Reality

#### Scenario: Full screen report
- **WHEN** a user clicks a report tile
- **THEN** the system displays the report in full screen with detailed data

### Requirement: Spanish language
The entire user interface SHALL be in Spanish, including all labels, messages, errors, and notifications.

#### Scenario: UI in Spanish
- **WHEN** a user accesses any page
- **THEN** all text is displayed in Spanish (menus, buttons, labels, error messages, help text)

#### Scenario: Date and number format
- **WHEN** the system displays dates and numbers
- **THEN** they use Spanish locale format (DD/MM/YYYY, separador decimal coma o punto según configuración regional)

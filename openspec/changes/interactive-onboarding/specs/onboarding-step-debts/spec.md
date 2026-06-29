## ADDED Requirements

### Requirement: Add debts step
The onboarding SHALL guide the user to add any credit cards or loans. Coach marks SHALL highlight the "+ TC" and "+ Préstamo" buttons in the sidebar. A tooltip SHALL explain: "Si tenés tarjetas de crédito o préstamos, agregalos. Chele maneja la deuda automáticamente."

#### Scenario: Coach marks on debt creation buttons
- **WHEN** user is on step 6
- **THEN** the overlay SHALL highlight "+ TC" and "+ Préstamo"
- **THEN** the tooltip SHALL explain debt management

#### Scenario: Debt created
- **WHEN** user creates a credit card or loan
- **THEN** the overlay SHALL detect and advance

#### Scenario: Skip allowed
- **WHEN** user has no debts
- **THEN** an "No tengo deudas →" button SHALL advance to step 7

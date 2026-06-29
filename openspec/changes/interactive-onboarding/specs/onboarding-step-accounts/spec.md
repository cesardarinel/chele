## ADDED Requirements

### Requirement: Create first account step
The onboarding SHALL guide the user to create at least one account. A coach mark SHALL highlight the "+ Cuenta" button in the sidebar. A tooltip SHALL explain: "Agregá tus cuentas bancarias. El saldo inicial aparece en 'Por asignar' — lo vas a distribuir después."

#### Scenario: Coach mark on + Cuenta
- **WHEN** user is on step 2
- **THEN** the overlay SHALL point to the "+ Cuenta" link in the sidebar
- **THEN** the tooltip SHALL show the explanation text
- **THEN** clicking the link SHALL open the account creation form normally

#### Scenario: Account created
- **WHEN** user submits the account creation form
- **THEN** the overlay SHALL detect the new account and show a confirmation
- **THEN** a "Siguiente →" button SHALL advance to step 3

#### Scenario: Skip available
- **WHEN** user has no accounts but wants to continue
- **THEN** a "Ya agregué todo →" button SHALL be available

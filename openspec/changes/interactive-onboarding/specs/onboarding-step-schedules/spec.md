## ADDED Requirements

### Requirement: Schedule recurring items step
The onboarding SHALL guide the user to set up at least one recurring schedule. A coach mark SHALL highlight the "Programaciones" nav link. A tooltip SHALL explain: "Si tenés ingresos o gastos que se repiten cada mes, programalos. Chele los aplica automáticamente."

#### Scenario: Coach mark on schedules nav
- **WHEN** user is on step 5
- **THEN** the overlay SHALL highlight the "Programaciones" link in the navigation
- **THEN** the tooltip SHALL explain schedules

#### Scenario: Schedule created
- **WHEN** user creates a schedule
- **THEN** the overlay SHALL detect and advance

#### Scenario: Skip allowed
- **WHEN** user does not want to create schedules
- **THEN** an "Omitir" link SHALL be available

## ADDED Requirements

### Requirement: Welcome step
The first onboarding step SHALL welcome the user and explain the envelope budgeting method concisely. It SHALL display as a centered card on the overlay with no coach mark.

#### Scenario: Welcome message
- **WHEN** user sees step 1
- **THEN** the overlay SHALL show: "💡 El método de los sobres — Cada peso en tus cuentas debe tener un trabajo. Vas a asignar tu dinero a categorías (sobres)."
- **THEN** a single button "Empezar →" SHALL advance to step 2

## ADDED Requirements

### Requirement: Completion step
The final onboarding step SHALL show a summary of what was accomplished during the onboarding and set `user.onboarding_step = 7`. The overlay SHALL never appear again.

#### Scenario: Completion summary
- **WHEN** user reaches step 7
- **THEN** the overlay SHALL display:
  - Number of accounts created
  - Total amount assigned (RTA = $0)
  - Number of goals configured
  - Number of schedules created
  - "🔍 ¿Qué sigue?" with tips: register daily expenses, review categories monthly, adjust goals
- **THEN** a button "Ir al presupuesto" SHALL set step to 7 and redirect to budget_view

#### Scenario: Never show again
- **WHEN** `user.onboarding_step == 7` on any request
- **THEN** the overlay SHALL NOT render
- **THEN** the middleware SHALL skip all onboarding checks

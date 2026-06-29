## ADDED Requirements

### Requirement: Interactive overlay with coach marks
The onboarding SHALL use a full-screen overlay with a semi-transparent dark background. The overlay SHALL:
- Block all pointer events on background content except the highlighted element
- Display a tooltip/coach mark card pointing to the current step's target element
- Include a progress bar at the bottom showing step N/7
- Support "Siguiente" and "Atrás" navigation buttons
- Animate transitions between steps (fade/slide)
- Be dismissable only when the step's completion condition is met

#### Scenario: Overlay with coach mark
- **WHEN** user is on step 2 (create account)
- **THEN** the overlay SHALL highlight the "+ Cuenta" button with a pulse animation
- **THEN** a tooltip SHALL appear pointing to the button with step explanation
- **THEN** all other UI elements SHALL be non-interactive

#### Scenario: Step completes automatically
- **WHEN** user creates a first account (step 2 condition met)
- **THEN** the overlay SHALL detect the change via polling and advance to step 3

#### Scenario: Progress bar
- **WHEN** onboarding overlay is visible
- **THEN** a fixed bar at the bottom SHALL show: `■■■■□□□□□  3/7  Asignar dinero`

#### Scenario: Mobile adaptation
- **WHEN** screen width < 768px
- **THEN** coach mark tooltips SHALL be bottom sheets (full width, slide up from bottom)
- **THEN** the progress bar SHALL be at the top of the bottom sheet

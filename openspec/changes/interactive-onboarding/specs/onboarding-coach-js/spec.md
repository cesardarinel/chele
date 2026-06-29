## ADDED Requirements

### Requirement: Coach marks JavaScript logic
The onboarding SHALL include JavaScript that handles:
- Overlay show/hide with pointer-events management
- Coach mark positioning (tooltip with arrow pointing to target element)
- Pulse animation on the highlighted element
- Polling every 3 seconds to check step completion conditions via `/onboarding/state`
- Auto-advance when conditions are met
- Mobile adaptation: tooltip → bottom sheet on small screens
- Keyboard navigation: Escape to close (where allowed), Enter to advance

#### Scenario: Polling for step completion
- **WHEN** user is on step 2 (create account)
- **THEN** JS SHALL poll `GET /onboarding/state` every 3 seconds
- **WHEN** response shows `step_completed: true`
- **THEN** JS SHALL show a success animation and enable "Siguiente"

#### Scenario: Element highlighting
- **WHEN** a step has a target element selector
- **THEN** JS SHALL add a highlight class to that element (z-index above overlay, pulse animation)
- **THEN** only that element SHALL receive pointer events

#### Scenario: Mobile detection
- **WHEN** `window.innerWidth < 768`
- **THEN** JS SHALL render tooltips as bottom sheets instead of floating tooltips
- **THEN** the progress bar SHALL move to the top of the bottom sheet

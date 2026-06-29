## ADDED Requirements

### Requirement: OnboardingMiddleware
The system SHALL have middleware that checks `user.onboarding_step` on every request. If `step < 7` and the request is not for an onboarding endpoint, the middleware SHALL inject the onboarding overlay context variables. If the user is at step 0, after registration they SHALL be redirected to `budget_view` with the overlay shown.

#### Scenario: Unobtrusive overlay
- **WHEN** user has `step < 7` and visits any page
- **THEN** the page SHALL render normally with an added overlay on top
- **THEN** the overlay SHALL block interaction with background elements except the highlighted coach mark

#### Scenario: Step 0 after registration
- **WHEN** a user registers and is redirected to `budget_view`
- **THEN** the onboarding overlay SHALL be visible with step 1 (welcome)

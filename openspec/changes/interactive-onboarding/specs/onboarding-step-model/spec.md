## ADDED Requirements

### Requirement: User onboarding_step field
The User model SHALL have an `onboarding_step` IntegerField with default=0. The field tracks the user's progress through the interactive onboarding wizard. Values: 0=not started, 1-6=in progress, 7=completed.

#### Scenario: New user has step 0
- **WHEN** a new user is created via registration
- **THEN** `user.onboarding_step` SHALL be 0

#### Scenario: User completes all steps
- **WHEN** user finishes step 7 of onboarding
- **THEN** `user.onboarding_step` SHALL be 7
- **THEN** the overlay SHALL never appear again for this user

#### Scenario: Existing users unaffected
- **WHEN** an existing user logs in after this migration
- **THEN** `user.onboarding_step` SHALL be set to 7 (assume already onboarded)

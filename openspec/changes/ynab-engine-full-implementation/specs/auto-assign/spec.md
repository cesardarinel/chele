## ADDED Requirements

### Requirement: Auto-assign underfunded categories
The system SHALL provide an "Auto-assign" feature that distributes available Ready to Assign funds across all underfunded categories in a single action. The system SHALL prioritize categories by some order (e.g., earliest due date, user-defined priority). The system SHALL only assign up to the available Ready to Assign amount.

#### Scenario: Auto-assign all underfunded
- **WHEN** user has $1,000 in Ready to Assign and 3 underfunded categories needing $300, $200, and $400 respectively
- **THEN** auto-assign SHALL distribute funds to all 3 categories, setting them to fully funded
- **THEN** Ready to Assign SHALL decrease to $100

#### Scenario: Auto-assign partial (insufficient funds)
- **WHEN** user has $500 in Ready to Assign but underfunded total is $800
- **THEN** auto-assign SHALL fund categories in priority order until Ready to Assign is exhausted

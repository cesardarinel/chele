## ADDED Requirements

### Requirement: Inspector panel
The system SHALL provide an Inspector panel on the right side of the budget view (desktop) that displays detailed information about the selected category. The panel SHALL show: current available balance, total assigned this month, total activity (spending), target details (if any), underfunded amount, average spending (3 months), and quick action buttons (Move Money, Cover, Snooze Target, Edit Target).

#### Scenario: Open inspector for a category
- **WHEN** user clicks on a category in the budget view
- **THEN** the Inspector panel SHALL slide in from the right showing that category's details

#### Scenario: Inspector shows target info
- **WHEN** a category has a $500 monthly target and $200 assigned this month
- **THEN** the Inspector panel SHALL show "Meta: $500/mes · Asignado: $200 · Faltan: $300"

#### Scenario: Inspector shows spending average
- **WHEN** a category has spending history
- **THEN** the Inspector panel SHALL display the average spending over the last 3 months

### Requirement: Inspector on mobile
On mobile screens, the Inspector SHALL be displayed as a bottom sheet modal instead of a side panel. All functionality SHALL be identical.

#### Scenario: Mobile inspector as bottom sheet
- **WHEN** user taps a category on mobile budget view
- **THEN** the category details SHALL appear in a bottom sheet that slides up from the bottom of the screen

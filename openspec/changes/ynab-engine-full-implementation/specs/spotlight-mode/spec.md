## ADDED Requirements

### Requirement: Spotlight mode
The system SHALL provide a centralized notification center ("Spotlight") that aggregates all pending actions and alerts. The Spotlight SHALL display: uncategorized transactions needing review, overspent categories needing coverage (red alert), targets that are significantly underfunded (yellow alert), and imported transactions pending approval. Each alert SHALL include a direct action button (Review, Cover, Assign).

#### Scenario: Spotlight shows uncategorized transactions
- **WHEN** there are 5 uncategorized transactions
- **THEN** the Spotlight SHALL show "5 transacciones sin categorizar" with a "Review" button

#### Scenario: Spotlight shows overspending
- **WHEN** a category has an uncovered overspent balance
- **THEN** the Spotlight SHALL show the category name and amount with a "Cover" button

#### Scenario: Review flow for imported transactions
- **WHEN** user clicks "Review" on uncategorized transactions in Spotlight
- **THEN** the system SHALL enter a review flow where each transaction is shown one by one for categorization
- **WHEN** user categorizes or approves the last transaction
- **THEN** the Spotlight SHALL update to reflect zero pending reviews

### Requirement: Spotlight visibility
The Spotlight SHALL be visible as a collapsible section at the top of the main content area (below the header, above the budget view). It SHALL show a count badge indicating the number of pending items. If there are no pending items, the Spotlight SHALL be hidden.

#### Scenario: Spotlight badge shows count
- **WHEN** there are 3 pending actions
- **THEN** the Spotlight section SHALL display a badge with "3"

#### Scenario: Spotlight hides when empty
- **WHEN** all pending actions have been resolved
- **THEN** the Spotlight section SHALL be hidden

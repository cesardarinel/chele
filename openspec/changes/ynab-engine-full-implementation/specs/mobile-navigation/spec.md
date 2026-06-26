## ADDED Requirements

### Requirement: Bottom navigation bar (mobile)
On mobile screens, the system SHALL display a fixed bottom navigation bar with three tabs: "Plan" (budget view), "Accounts" (account list), and "Spending" (reports/analytics). Each tab SHALL have an icon and label. The active tab SHALL be visually highlighted. Tapping a tab SHALL navigate to the corresponding view without a full page reload (client-side routing via JavaScript or server-side URLs).

#### Scenario: Navigate via bottom bar
- **WHEN** user taps "Accounts" in the bottom navigation bar
- **THEN** the system SHALL display the accounts list view

#### Scenario: Active tab indicator
- **WHEN** user is viewing the budget
- **THEN** the "Plan" tab SHALL be visually highlighted (active state)

### Requirement: Bottom bar replaces mobile header navigation
When the bottom navigation bar is active, the mobile header SHALL be simplified to show only the page title and a hamburger menu for the sidebar.

#### Scenario: Simplified header with bottom bar
- **WHEN** bottom navigation bar is displayed on mobile
- **THEN** the mobile header SHALL only show the current page title and hamburger menu

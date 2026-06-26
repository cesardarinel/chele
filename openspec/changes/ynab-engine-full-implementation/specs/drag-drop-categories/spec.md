## ADDED Requirements

### Requirement: Drag-and-drop category reordering
The system SHALL support drag-and-drop reordering of categories within a group and groups within the budget. When the user drags a category to a new position, the sort_order field SHALL be updated for all affected categories. When a category is dragged to a different group, its group_id SHALL be updated.

#### Scenario: Reorder categories
- **WHEN** user drags "Groceries" above "Eating Out" in the same group
- **THEN** the system SHALL persist the new order via POST request

#### Scenario: Move category between groups
- **WHEN** user drags "Phone" from "Gastos Fijos" group to "Gastos Diarios" group
- **THEN** the category SHALL move to the target group at the dropped position

### Requirement: Multi-select for batch operations
The system SHALL allow selecting multiple categories for batch operations (e.g., batch assign via Auto-assign).

#### Scenario: Multi-select categories
- **WHEN** user Cmd+clicks (Mac) or Ctrl+clicks (Windows) on multiple categories
- **THEN** the selected categories SHALL be highlighted and available for batch actions

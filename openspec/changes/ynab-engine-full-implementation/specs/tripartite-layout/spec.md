## ADDED Requirements

### Requirement: Tripartite desktop layout
The desktop layout SHALL consist of three panels:
1. **Left panel** (16rem): Account list and balances (existing sidebar)
2. **Center panel** (1fr): Budget view with categories expanded to show Assigned, Activity, and Available columns
3. **Right panel** (20rem): Inspector panel (shown when a category is selected, hidden otherwise)
The layout SHALL use CSS Grid: `grid-template-columns: 16rem 1fr 20rem`.

#### Scenario: Three panel layout
- **WHEN** user opens the budget view on a desktop screen (>1024px)
- **THEN** they SHALL see the left sidebar, center budget table, and (if a category is selected) the right inspector panel

#### Scenario: Inspector panel hides on small desktop
- **WHEN** user resizes browser window to <1024px
- **THEN** the Inspector panel SHALL hide automatically (accessible via button or modal)

## ADDED Requirements

### Requirement: Set up goals step
The onboarding SHALL guide the user to create at least one goal/target. A coach mark SHALL highlight any category in the budget view. A tooltip SHALL explain: "Una meta te dice cuánto asignar cada mes para llegar a un objetivo. Hacé click en una categoría y agregá una meta."

#### Scenario: Coach mark on a category
- **WHEN** user is on step 4
- **THEN** the overlay SHALL highlight a category row in the budget table
- **THEN** the tooltip SHALL explain goals

#### Scenario: Goal created
- **WHEN** user creates a goal via the inspector panel
- **THEN** the overlay SHALL detect the new goal and advance
- **THEN** a success message SHALL show: "✅ Meta creada"

#### Scenario: Skip allowed
- **WHEN** user does not want to set goals
- **THEN** an "Omitir — lo haré después" link SHALL be available

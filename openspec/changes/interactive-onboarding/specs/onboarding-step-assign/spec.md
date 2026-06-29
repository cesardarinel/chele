## ADDED Requirements

### Requirement: Assign all money (zero-sum forced)
The onboarding SHALL force the user to assign ALL their Ready to Assign balance to categories. The coach mark SHALL highlight the "Por asignar" section. A tooltip SHALL explain: "Este es tu dinero sin trabajo. Asignaldo a categorías usando los inputs de cada una. El objetivo es que 'Por asignar' llegue a $0."

#### Scenario: RTA highlighted
- **WHEN** user is on step 3
- **THEN** the "Por asignar" section SHALL pulse/highlight
- **THEN** a tooltip SHALL explain the zero-sum concept

#### Scenario: User assigns money
- **WHEN** user types an amount in a category input and submits
- **THEN** the "Por asignar" value SHALL decrease in real-time
- **THEN** the progress bar SHALL update: `$5,000 / $30,000 asignado`

#### Scenario: Cannot advance until RTA = $0
- **WHEN** `ready_to_assign > 0`
- **THEN** the "Siguiente" button SHALL be disabled with text: "Asigná todo tu dinero primero"
- **THEN** the overlay SHALL show the remaining amount: "Falta asignar: $X"

#### Scenario: RTA reaches $0
- **WHEN** `ready_to_assign == 0`
- **THEN** the overlay SHALL show a success animation
- **THEN** "Siguiente →" SHALL be enabled with text "✅ Todo tiene trabajo"
- **THIS STEP IS NOT SKIPPABLE**

#### Scenario: Visual feedback
- **WHEN** user is on step 3
- **THEN** each category input SHALL show a small icon indicating it has "trabajo" (💰) once assigned

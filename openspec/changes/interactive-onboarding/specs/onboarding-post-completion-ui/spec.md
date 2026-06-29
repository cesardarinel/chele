## ADDED Requirements

### Requirement: Post-onboarding persistent UI indicators
After the onboarding is complete, the following UI indicators SHALL remain visible to help users understand their money status at a glance:

#### Scenario: RTA explanation tooltip
- **WHEN** user hovers over "Por asignar" on the budget view
- **THEN** a tooltip SHALL show: "De tu saldo total, esto es lo que aún no tiene un trabajo asignado."

#### Scenario: Category balance breakdown
- **WHEN** user opens the inspector panel for a category
- **THEN** the "Disponible" SHALL show a breakdown: `$200 disponible ($150 del mes pasado + $50 asignado nuevo)`
- **THEN** the target progress bar SHALL be visible: `▰▰▰▰▰▱▱▱▱▱ $500/$1,000`

#### Scenario: Sidebar account indicator
- **WHEN** user views the sidebar account list
- **THEN** accounts with `on_budget=True` SHALL show a small label "💰 En presupuesto"
- **THEN** accounts with `on_budget=False` SHALL show "🏦 Ahorro"

#### Scenario: Zero-sum confirmation
- **WHEN** RTA is $0
- **THEN** the "Por asignar" section SHALL show a subtle ✅ "Todo asignado"

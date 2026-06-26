## ADDED Requirements

### Requirement: Cost to be me calculation
The system SHALL calculate and display the "Cost to be me" — the total monthly amount needed to fully fund all active targets. This SHALL be shown as a summary figure. The system SHALL also calculate and display the "Expected monthly income" (average of last 3 months of income transactions or manual entry).

#### Scenario: View cost to be me
- **WHEN** user has targets totaling $3,500 per month
- **THEN** the system SHALL display "$3,500/mes" as the Cost to be me

### Requirement: Reality Check alert
The system SHALL compare Cost to be me against Expected monthly income. If Cost to be me exceeds Expected income, the system SHALL display a warning alert suggesting the user review their targets.

#### Scenario: Reality Check warning
- **WHEN** Cost to be me ($4,000) exceeds Expected income ($3,200)
- **THEN** the system SHALL display an alert: "Tus metas ($4,000) superan tus ingresos esperados ($3,200). Revisa tus metas para equilibrar tu presupuesto."

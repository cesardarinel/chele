## ADDED Requirements

### Requirement: Full Spanish translation
The entire application SHALL be in Spanish, hardcoded in templates for the MVP.

#### Scenario: All UI text in Spanish
- **WHEN** any page renders
- **THEN** every string (labels, buttons, titles, messages, errors, tooltips, placeholders) is in Spanish

#### Scenario: Number formatting
- **WHEN** monetary amounts are displayed
- **THEN** they use Spanish format with thousands separator (.) and decimal comma (,) or configurable format

#### Scenario: Date formatting
- **WHEN** dates are displayed
- **THEN** they use DD/MM/YYYY format

### Requirement: Spanish terminology consistency
The system SHALL use consistent Spanish terminology across all views.

#### Scenario: Consistent terms
- **WHEN** the same concept appears in different views
- **THEN** the same Spanish term is used (e.g., "Beneficiario" not "Pagador" or "Proveedor" interchangeably)

#### Scenario: Category naming
- **WHEN** creating default categories
- **THEN** they are named in Spanish (e.g., "Comida", "Servicios", "Ahorro", "General")

## ADDED Requirements

### Requirement: Cover overspending
The system SHALL provide a "Cover" action for overspent categories. The Cover action SHALL allow the user to select another category with positive available balance and transfer funds to cover the overspent amount. The transfer SHALL be atomic: both categories update simultaneously.

#### Scenario: Cover overspent from another category
- **WHEN** "Groceries" shows -$15 overspent and user clicks "Cover"
- **THEN** system SHALL prompt to select a source category with sufficient available funds
- **WHEN** user selects "Eating Out" (available: $100) and confirms $15
- **THEN** "Groceries" available SHALL become $0 and "Eating Out" available SHALL become $85

#### Scenario: Cover with insufficient funds
- **WHEN** user tries to cover an overspent category from a source that has insufficient available funds
- **THEN** the system SHALL show an error: "Fondos insuficientes en la categoría seleccionada"

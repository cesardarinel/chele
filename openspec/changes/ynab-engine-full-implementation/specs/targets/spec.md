## ADDED Requirements

### Requirement: Target configuration per category
The system SHALL allow users to set financial targets (metas) on any category. Target types SHALL include: monthly (fixed amount each month), yearly (annual amount prorated monthly), target_balance (accumulate until reaching a balance), target_date (reach a balance by a specific date), true_expense (annual expense split into N monthly installments). Each target SHALL have a visual indicator of its funding status.

#### Scenario: Set monthly target
- **WHEN** user selects "Car Insurance" category and adds a Yearly target of $1,200 due in October
- **THEN** the system SHALL calculate that $200 needs to be assigned each month to stay on track

#### Scenario: Target visual indicator
- **WHEN** a category has a target and the assigned amount is less than the target requires
- **THEN** the category SHALL display a yellow "Underfunded" indicator

#### Scenario: Target fully funded
- **WHEN** the assigned amount for a category equals or exceeds the target requirement
- **THEN** the category SHALL display a green "Fully funded" indicator

### Requirement: Target underfunded calculation
The system SHALL calculate the underfunded amount for each target. For monthly targets: underfunded = target_amount - assigned_this_month. For yearly targets: underfunded = (target_amount / 12) - assigned_this_month. For target_balance: underfunded = target_amount - (current_balance + assigned_this_month - spent_this_month). For target_date: underfunded = remaining_amount / months_remaining - assigned_this_month.

#### Scenario: Yearly target proration
- **WHEN** a $1,200 Yearly target is set in January with due date in December
- **THEN** the system SHALL show $100 as the monthly funding requirement ($1,200 / 12)

#### Scenario: Underfunded shows remaining
- **WHEN** user has a monthly target of $500 and has assigned $200 this month
- **THEN** the underfunded amount for that category SHALL be $300

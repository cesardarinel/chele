## 1. Model: Target

- [ ] 1.1 Create `Target` model with fields: category (FK), goal_type (choices: monthly/yearly/target_balance/target_date/true_expense), amount, target_date (nullable), frequency (int, default 12), refill_up_to (bool, default false), snooze_month (int, nullable), snooze_year (int, nullable), is_completed (bool)
- [ ] 1.2 Create migration for Target model
- [ ] 1.3 Create migration script to convert existing Goal records to Target records (one-shot)
- [ ] 1.4 Register Target in Django admin

## 2. Ready to Assign & Zero-Sum

- [ ] 2.1 Implement `ReadyToAssignService` that calculates: sum(on_budget_balances) - sum(budgeted_current_month)
- [ ] 2.2 Add Ready to Assign display in budget_view template (prominent section at top)
- [ ] 2.3 Add zero-sum validation in `assign_funds` view: reject if assigning more than Ready to Assign
- [ ] 2.4 Update `move_funds` view to validate against Ready to Assign
- [ ] 2.5 Add cache (30s) for Ready to Assign value to prevent recomputation on every request

## 3. Target Engine & Underfunded Calculation

- [ ] 3.1 Implement `TargetService.underfunded(target, category, month, year)` supporting all goal types
- [ ] 3.2 Implement refill_up_to calculation: underfunded = max(0, amount - category_available_balance)
- [ ] 3.3 Implement set_aside_another calculation: underfunded = amount (ignores rollover)
- [ ] 3.4 Implement yearly target proration: monthly = amount / 12
- [ ] 3.5 Implement target_date proration: monthly = (amount - already_assigned) / months_remaining
- [ ] 3.6 Implement true_expense proration: monthly = amount / frequency
- [ ] 3.7 Add `underfunded` property to budget_view context per category
- [ ] 3.8 Add visual indicator (yellow/amber) for underfunded categories in budget view
- [ ] 3.9 Add visual indicator (green) for fully funded categories

## 4. Snooze Target

- [ ] 4.1 Implement snooze logic: set snooze_month/snooze_year on Target
- [ ] 4.2 Exclude snoozed targets from underfunded calculation in TargetService
- [ ] 4.3 Add auto-expire: clear snooze fields at month rollover
- [ ] 4.4 Add "Snooze" button in inspector panel and category dropdown

## 5. Overspending Handling

- [ ] 5.1 Implement `OverspendingService.detect(transaction)` that checks if category is overspent
- [ ] 5.2 Implement RED indicator for cash overspend (negative available balance)
- [ ] 5.3 Implement ORANGE indicator for credit card overspend
- [ ] 5.4 Implement auto-move logic: when TC overspend detected, move available funds to TC Payment category
- [ ] 5.5 Add overspent carry-forward logic: negative balance rolls to next month
- [ ] 5.6 Integrate overspending check into transaction_create, transaction_edit, transaction_delete views

## 6. Cover Overspending

- [ ] 6.1 Implement "Cover" action endpoint POST `/presupuestos/cubrir/` with from_category + amount
- [ ] 6.2 Add Cover button in Spotlight for overspent categories
- [ ] 6.3 Add Cover flow UI: select source category with sufficient available balance
- [ ] 6.4 Implement atomic transfer: reduce source category, zero out overspent category
- [ ] 6.5 Add validation: reject if source category has insufficient funds

## 7. Auto-Assign (Underfunded)

- [ ] 7.1 Implement `AutoAssignService.assign_all(budget, month, year)` that distributes Ready to Assign across underfunded categories
- [ ] 7.2 Define priority order for underfunded categories (by type: cash overspend first, then true_expense, then monthly targets)
- [ ] 7.3 Add "Auto-assign" button in budget view header
- [ ] 7.4 Implement partial assignment: if Ready to Assign < total underfunded, fund in priority order until exhausted

## 8. Inspector Panel

- [ ] 8.1 Create endpoint GET `/presupuestos/categoria/<id>/inspector/` returning JSON with category details
- [ ] 8.2 Create inspector panel template (right side panel on desktop)
- [ ] 8.3 Implement category selection via JavaScript (click on category row → show inspector)
- [ ] 8.4 Display in inspector: available balance, assigned, activity, target info, underfunded, 3-month average
- [ ] 8.5 Add quick action buttons in inspector: Move Money, Cover, Snooze, Edit Target
- [ ] 8.6 Create mobile bottom-sheet version of inspector (modal)

## 9. Spotlight Mode

- [ ] 9.1 Implement `SpotlightService.get_alerts(budget)` that aggregates: uncategorized transactions, uncovered overspends, underfunded targets
- [ ] 9.2 Create Spotlight template (collapsible section below header, above budget view)
- [ ] 9.3 Add badge showing count of pending items
- [ ] 9.4 Hide Spotlight when no pending items
- [ ] 9.5 Implement Review flow for uncategorized transactions (show one by one, categorize, approve)
- [ ] 9.6 Add "Cover" action directly in Spotlight alerts

## 10. Traffic Light UI

- [ ] 10.1 Apply GREEN color to: fully funded categories, positive available balances
- [ ] 10.2 Apply YELLOW/AMBER to: underfunded categories, credit card overspends
- [ ] 10.3 Apply RED to: cash overspends (negative available)
- [ ] 10.4 Add color indicators to sidebar account balances (consistent with category colors)
- [ ] 10.5 Add CSS classes for traffic light colors (reuse existing chele-success/warning/danger)

## 11. Tripartite Desktop Layout

- [ ] 11.1 Restructure budget view layout with CSS Grid: `grid-template-columns: 16rem 1fr`
- [ ] 11.2 Add right panel column (inspector) with `grid-template-columns: 16rem 1fr 20rem` when active
- [ ] 11.3 Add responsive breakpoint: hide inspector panel at <1024px
- [ ] 11.4 Ensure sidebar (left panel) retains its current scroll behavior

## 12. Mobile Bottom Navigation

- [ ] 12.1 Create fixed bottom navigation bar with 3 tabs: Plan, Accounts, Spending
- [ ] 12.2 Add icons for each tab (using existing SVG icons)
- [ ] 12.3 Implement tab switching (client-side show/hide or server-side navigation)
- [ ] 12.4 Add active tab indicator (highlighted state)
- [ ] 12.5 Simplify mobile header when bottom bar is active (title + hamburger only)

## 13. Drag-and-Drop Category Reordering

- [ ] 13.1 Integrate SortableJS library for drag-and-drop
- [ ] 13.2 Create endpoint POST `/presupuestos/categorias/reordenar/` accepting ordered category IDs
- [ ] 13.3 Implement sort_order update for all affected categories on drop
- [ ] 13.4 Implement cross-group drag: update group_id on category when moved to different group
- [ ] 13.5 Enable multi-select for batch operations (Cmd/Ctrl+click)

## 14. Cost to Be Me & Reality Check

- [ ] 14.1 Implement calculation of total monthly target cost (sum of all active target monthly requirements)
- [ ] 14.2 Implement Expected Income calculation (average of last 3 months income transactions)
- [ ] 14.3 Create "Cost to be me" widget in budget view header
- [ ] 14.4 Implement Reality Check alert: show warning if cost > expected income
- [ ] 14.5 Add link to adjust targets from Reality Check alert

## 15. Rollover Logic

- [ ] 15.1 Implement monthly rollover calculation: available = previous_month_available + budgeted - spent
- [ ] 15.2 Integrate rollover into budget_view context per category per month
- [ ] 15.3 Ensure negative rollover (cash overspend) carries as RED indicator
- [ ] 15.4 Ensure positive rollover with Refill up to target reduces underfunded amount

## 16. Go API Endpoints

- [ ] 16.1 Add Target CRUD endpoints in Go API
- [ ] 16.2 Add Ready to Assign endpoint GET `/api/budgets/:id/ready-to-assign`
- [ ] 16.3 Add Auto-assign endpoint POST `/api/budgets/:id/auto-assign`
- [ ] 16.4 Add Cover endpoint POST `/api/budgets/:id/cover`
- [ ] 16.5 Add Spotlight alerts endpoint GET `/api/budgets/:id/spotlight`
- [ ] 16.6 Add Cost to be me endpoint GET `/api/budgets/:id/cost-to-be-me`

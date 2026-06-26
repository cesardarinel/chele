## 1. Model: Target

- [x] 1.1 Create Target model
- [x] 1.2 Create migration for Target model
- [ ] 1.3 Create migration script to convert existing Goal records to Target records (one-shot)
- [x] 1.4 Register Target in Django admin

## 2. Ready to Assign & Zero-Sum

- [x] 2.1 Implement ReadyToAssign calculation
- [x] 2.2 Add Ready to Assign display in budget_view template
- [ ] 2.3 Add zero-sum validation in `assign_funds` view
- [ ] 2.4 Update `move_funds` view to validate against Ready to Assign
- [ ] 2.5 Add cache for Ready to Assign value

## 3. Target Engine & Underfunded Calculation

- [x] 3.1 Implement TargetService.underfunded()
- [x] 3.2 Implement refill_up_to calculation
- [x] 3.3 Implement set_aside_another calculation
- [x] 3.4 Implement yearly target proration
- [x] 3.5 Implement target_date proration
- [x] 3.6 Implement true_expense proration
- [x] 3.7 Add underfunded to budget_view context
- [x] 3.8 Add visual indicator (yellow) for underfunded categories
- [x] 3.9 Add visual indicator (green) for fully funded categories

## 4. Snooze Target

- [x] 4.1 Implement snooze logic (snooze_month/snooze_year)
- [x] 4.2 Exclude snoozed from underfunded calculation
- [ ] 4.3 Add auto-expire: clear snooze fields at month rollover
- [ ] 4.4 Add "Snooze" button in inspector panel

## 5. Overspending Handling

- [ ] 5.1 Implement OverspendingService.detect()
- [x] 5.2 Implement RED indicator for cash overspend
- [ ] 5.3 Implement ORANGE indicator for credit card overspend
- [ ] 5.4 Implement auto-move logic for TC Payment
- [ ] 5.5 Add overspent carry-forward logic
- [ ] 5.6 Integrate overspending check into transaction views

## 6. Cover Overspending

- [x] 6.1 Implement Cover action POST `/presupuestos/cubrir/`
- [ ] 6.2 Add Cover button in Spotlight
- [ ] 6.3 Add Cover flow UI (source category selector)
- [x] 6.4 Implement atomic transfer
- [x] 6.5 Add validation for insufficient funds

## 7. Auto-Assign (Underfunded)

- [x] 7.1 Implement AutoAssign service (Go API)
- [x] 7.2 Define priority order for underfunded
- [ ] 7.3 Add "Auto-assign" button in budget view header
- [x] 7.4 Implement partial assignment

## 8. Inspector Panel

- [x] 8.1 Create inspector endpoint (Go API)
- [ ] 8.2 Create inspector panel template (right side panel)
- [ ] 8.3 Implement category selection via JavaScript
- [ ] 8.4 Display in inspector: balance, assigned, activity, target, avg
- [ ] 8.5 Add quick action buttons
- [ ] 8.6 Create mobile bottom-sheet version

## 9. Spotlight Mode

- [x] 9.1 Implement Spotlight endpoint (Go API)
- [ ] 9.2 Create Spotlight template (collapsible section)
- [ ] 9.3 Add badge with count
- [ ] 9.4 Hide when no pending items
- [ ] 9.5 Implement Review flow for uncategorized transactions
- [ ] 9.6 Add Cover action in Spotlight alerts

## 10. Traffic Light UI

- [x] 10.1 Apply GREEN to funded categories
- [x] 10.2 Apply YELLOW to underfunded
- [x] 10.3 Apply RED to cash overspends
- [ ] 10.4 Add color indicators to sidebar account balances
- [x] 10.5 Add CSS classes for traffic light colors

## 11. Tripartite Desktop Layout

- [ ] 11.1 Restructure with CSS Grid `16rem 1fr`
- [ ] 11.2 Add right panel column `16rem 1fr 20rem`
- [ ] 11.3 Hide inspector panel at <1024px
- [x] 11.4 Ensure sidebar scroll behavior

## 12. Mobile Bottom Navigation

- [ ] 12.1 Create fixed bottom nav bar with 3 tabs
- [ ] 12.2 Add icons for each tab
- [ ] 12.3 Implement tab switching
- [ ] 12.4 Add active tab indicator
- [ ] 12.5 Simplify mobile header

## 13. Drag-and-Drop Category Reordering

- [ ] 13.1 Integrate SortableJS
- [ ] 13.2 Create reorder endpoint
- [ ] 13.3 Implement sort_order update
- [ ] 13.4 Implement cross-group drag
- [ ] 13.5 Enable multi-select

## 14. Cost to Be Me & Reality Check

- [x] 14.1 Implement total monthly target cost calculation
- [x] 14.2 Implement Expected Income calculation
- [ ] 14.3 Create "Cost to be me" widget in budget view header
- [ ] 14.4 Implement Reality Check alert
- [ ] 14.5 Add link to adjust targets

## 15. Rollover Logic

- [ ] 15.1 Implement monthly rollover calculation
- [ ] 15.2 Integrate rollover into budget_view context
- [ ] 15.3 Negative rollover carries as RED
- [x] 15.4 Refill up to reduces underfunded with rollover

## 16. Go API Endpoints

- [x] 16.1 Target CRUD endpoints
- [x] 16.2 Ready to Assign endpoint
- [x] 16.3 Auto-assign endpoint
- [x] 16.4 Cover endpoint
- [x] 16.5 Spotlight alerts endpoint
- [x] 16.6 Cost to be me endpoint

## Context

The Go API at `chele-api/` already has models for accounts, budgets, transactions, schedules, credit cards, loans, payees, goals, and rules. It has working CRUD handlers and services for balance management, interest calculation, schedule processing, and dashboard computation. However, it lacks all YNAB-specific logic: targets, underfunded calculation, auto-assign, cover, overspending detection, spotlight aggregation, rollover, and cost-to-be-me. These need to be added as new handlers and services while maintaining the existing patterns.

## Goals / Non-Goals

**Goals:**
- Target CRUD endpoints with all goal types (monthly, yearly, target_balance, target_date, true_expense)
- Ready to Assign calculation endpoint
- Auto-assign endpoint with priority ordering
- Cover overspending endpoint with atomic transfer
- Spotlight alerts aggregation endpoint
- Category inspector detail endpoint
- Cost to be me calculation endpoint
- Rollover balances endpoint

**Non-Goals:**
- UI changes (these belong to the Django/HTML layer in the `ynab-engine-full-implementation` change)
- Go tests for the new endpoints (to be added separately)
- Migration of existing Django Goal data (handled by the Django-side change)

## Decisions

1. **New Target model in Go**: Add `Target` struct to `internal/models/models.go` with all fields matching the Django Target model. Table: `goals_goal` (reuse existing table since Goals will be migrated to Targets).

2. **Underfunded calculation as service**: `internal/service/target.go` with `CalculateUnderfunded(target, categoryID, month, year)` method supporting all goal types, refill-up-to and snooze.

3. **Ready to Assign as simple query**: SUM on-budget accounts minus SUM MonthlyBudget.budgeted. No caching needed — API calls are expected to be less frequent than page loads.

4. **Auto-assign priority order**: Implement as a service that queries underfunded categories, sorts by priority (cash overspend > true_expense > targets with earliest due date), and assigns until Ready to Assign is exhausted.

5. **Cover as two UPDATE transactions**: Debit source MonthlyBudget, credit target MonthlyBudget, wrapped in a transaction. Same pattern as existing `move_funds`.

6. **Spotlight uses multiple queries**: Three independent queries: uncategorized transactions count, uncovered overspends list (negative available), underfunded targets list. Aggregated in the handler.

7. **Inspector uses single category query**: Loads MonthlyBudget for given category/month/year + Target for that category + Transaction averages. Returns combined JSON.

8. **Rollover calculated per category**: For each category with MonthlyBudget data, calculate `available + budgeted - spent = rollover`. Return as array.

## Risks / Trade-offs

- [Performance] Spotlight endpoint runs 3 queries per call. For budgets with many categories, this could be slow. Mitigation: add month/year filters to limit scope.
- [Data consistency] Cover and Auto-assign modify MonthlyBudget records concurrently. Risk of race conditions if Django modifies same records. Mitigation: low risk since Django operations are serialized per request.
- [Model reuse] Using `goals_goal` table for Targets means the Go API and Django must agree on the schema. The Django migration that renames Goal→Target must update column names that Go expects.

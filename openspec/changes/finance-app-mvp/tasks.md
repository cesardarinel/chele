## 1. Project Setup

- [x] 1.1 Create Django project `chele` with `django-admin startproject`
- [x] 1.2 Configure SQLite with WAL mode in settings
- [x] 1.3 Configure Docker + docker-compose for dev and production environment
- [x] 1.4 Add HTMX, Django templates + Tailwind CSS configuration
- [x] 1.5 Configure environment variables (.env)
- [x] 1.6 Set up base template with sidebar navigation layout (replicando ActualBudget)
- [x] 1.7 Set all UI text in Spanish (hardcodeado en templates)

## 2. Budgets Management

- [x] 2.1 Create Budget model (name, description, owner, created_at) — reemplaza el concepto de "Group"
- [x] 2.2 Create BudgetMembership model (user, budget, role)
- [x] 2.3 Migrate all FK relations from "Group" to "Budget" (accounts, categories, transactions, payees, rules, schedules)
- [x] 2.4 Implement budget creation view + template
- [x] 2.5 Implement budget selector dropdown in sidebar
- [x] 2.6 Implement switch budget logic (reload data for selected budget)
- [x] 2.7 Implement member invitation per budget
- [x] 2.8 Add tests for budgets

## 3. User Authentication

- [x] 3.1 Create User model (extend AbstractUser)
- [x] 3.2 Implement registration view + template
- [x] 3.3 Implement login view + template
- [x] 3.4 Implement logout
- [x] 3.5 Add password validation and hashing
- [x] 3.6 Add session management with configurable timeout
- [x] 3.7 Add tests for auth flow

## 4. Members / Invitations

- [x] 4.1 Implement invitation via email per budget
- [x] 4.2 Implement invitation acceptance flow
- [x] 4.3 Implement budget member list view
- [x] 4.4 Add tests for invitations

## 5. Account Management

- [x] 5.1 Create Account model (type, name, balance, budget FK, on_budget flag)
- [x] 5.2 Implement account CRUD views + templates
- [x] 5.3 Implement balance calculation from transactions
- [x] 5.4 Implement net worth calculation and display
- [x] 5.5 Support account types: checking, savings, cash, credit card
- [x] 5.6 Support on-budget/off-budget toggle
- [x] 5.7 Add tests for accounts

## 6. Envelope Budgeting + YNAB Methodology

- [x] 6.1 Create CategoryGroup model (name, budget FK)
- [x] 6.2 Create Category model (name, group FK, budget FK)
- [x] 6.3 Create MonthlyBudget model (category, month, year, allocated_amount)
- [x] 6.4 Implement month navigation view + template
- [x] 6.5 Implement category group CRUD views + templates
- [x] 6.6 Implement category CRUD views + templates
- [x] 6.7 Implement fund allocation UI + logic (Presupuestado | Gastado | Saldo)
- [x] 6.8 Implement available-to-budget calculation and zero-sum enforcement
- [x] 6.9 Implement move money between categories (Regla 3 YNAB)
- [x] 6.10 Implement overspending rollover to next month (Regla 3 YNAB)
- [x] 6.11 Implement Hold for Next Month (Regla 4 YNAB)
- [x] 6.12 Implement auto-hold on income categories
- [x] 6.13 Implement copy budget from previous month / set to averages
- [x] 6.14 Add tests for budgeting

## 7. Transactions

- [x] 7.1 Create Transaction model (amount, date, payee, account, category, notes, transfer_id, budget FK)
- [x] 7.2 Use UUID as primary key for all models (offline-ready)
- [x] 7.3 Implement transaction CRUD views + templates
- [x] 7.4 Implement paginated transaction list by account/category
- [x] 7.5 Implement transfer between accounts (transfer pair)
- [x] 7.6 Ensure balance and budget updates on CRUD operations
- [x] 7.7 Add tests for transactions

## 8. Credit Cards

- [x] 8.1 Extend Account model to support credit_card type
- [x] 8.2 Create "Pago TC" category automatically for TC accounts
- [x] 8.3 Implement TC expense handling: move funds from category to TC Payment category
- [x] 8.4 Implement TC payment (transfer from another account clears TC Payment category)
- [x] 8.5 Add tests for credit card flow

## 9. Category Goals / Metas

- [x] 9.1 Create Goal model (category, goal_type, amount, target_date, frequency)
- [x] 9.2 Implement Monthly Savings Goal: auto-fill budgeted amount each month
- [x] 9.3 Implement Target Balance: calculate needed amount and auto-fill
- [x] 9.4 Implement Target by Date: calculate monthly contribution
- [x] 9.5 Implement True Expense: split annual cost into monthly payments
- [x] 9.6 Show goal progress in budget view
- [x] 9.7 Mark goals as complete when reached
- [x] 9.8 Add tests for goals

## 10. CSV Import

- [x] 10.1 Create CSV parser with configurable column mapping
- [x] 10.2 Implement CSV upload view + template
- [x] 10.3 Implement column mapping UI (drag/drop or select)
- [x] 10.4 Implement preview before confirmation
- [x] 10.5 Implement duplicate detection (amount + date + payee)
- [x] 10.6 Implement onboarding CSV import flow
- [x] 10.7 Implement recurring CSV import flow
- [x] 10.8 Add tests for CSV import

## 11. UI Layout & Sidebar

- [x] 11.1 Implement persistent sidebar with budget selector dropdown + navigation items
- [x] 11.2 Implement sidebar account list with balances (del presupuesto activo)
- [x] 11.3 Implement "Más" dropdown (Beneficiarios, Reglas)
- [x] 11.4 Implement budget view table layout (Presupuestado | Gastado | Saldo)
- [x] 11.5 Implement month header with note, chevrons, 3-dot menu
- [x] 11.6 Implement Account Register view layout (transaction list, Importar, Añadir, Filtrar, Buscar)
- [x] 11.7 Implement Reports Dashboard tile layout
- [x] 11.8 Ensure 100% of UI text is in Spanish

## 12. Reports

- [x] 12.1 Implement Budget vs Reality report view + template
- [x] 12.2 Implement Net Worth report view + template
- [x] 12.3 Implement Cash Flow report view + template
- [x] 12.4 Implement category breakdown view
- [x] 12.5 Implement Income vs Expense summary view
- [x] 12.6 Add month/period navigation for reports
- [x] 12.7 Implement full-screen report view
- [x] 12.8 Add tests for reports

## 13. Programaciones (Schedules)

- [x] 13.1 Create Schedule model (budget FK, payee, amount, category, frequency, next_date, account)
- [x] 13.2 Implement schedule CRUD views + templates
- [x] 13.3 Implement auto-creation of transactions from schedules
- [x] 13.4 Add tests for schedules

## 14. Beneficiarios (Payees)

- [x] 14.1 Create Payee model (name, budget FK)
- [x] 14.2 Implement payee management view + template
- [x] 14.3 Add rename/merge payee functionality
- [x] 14.4 Add tests for payees

## 15. Reglas (Rules)

- [x] 15.1 Create Rule model (budget FK, conditions, actions, order)
- [x] 15.2 Implement rule CRUD views + templates
- [x] 15.3 Implement rule execution on transaction import
- [x] 15.4 Add tests for rules

## 16. Configuración (Settings)

- [x] 16.1 Implement user settings view (profile, password change)
- [x] 16.2 Implement budget settings view (name, member management)
- [x] 16.3 Add tests for settings

## 17. Sync Engine

- [x] 17.1 Add `updated_at` and `sync_status` fields to all syncable models
- [x] 17.2 Implement sync endpoint: push pending changes
- [x] 17.3 Implement sync endpoint: pull remote changes
- [x] 17.4 Implement last-write-wins conflict resolution
- [x] 17.5 Implement pull-to-refresh trigger in UI
- [x] 17.6 Implement offline queuing (localStorage or IndexedDB)
- [x] 17.7 Add tests for sync engine

## 18. Docker Deployment

- [x] 18.1 Create production Dockerfile
- [x] 18.2 Create docker-compose.yml for production
- [x] 18.3 Configure static files serving (Whitenoise or nginx)
- [x] 18.4 Configure environment variables for production
- [x] 18.5 Test full deployment flow

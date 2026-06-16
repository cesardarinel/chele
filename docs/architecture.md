# Architecture

## Stack

- Python 3.13, Django 6.0, SQLite (WAL), Tailwind CSS (CDN)
- Docker + Gunicorn + Whitenoise on DigitalOcean App Platform

## Project Structure

```
chele/
├── chele/              # Django project config
│   ├── settings.py
│   ├── urls.py
│   ├── wsgi.py
│   └── context_processors.py
├── apps/               # All Django apps
│   ├── accounts/       # Account CRUD + balance
│   ├── budgets/        # Budgets, categories, YNAB logic
│   ├── credit_cards/   # Credit cards with interest
│   ├── csv_import/     # CSV parsing + import
│   ├── goals/          # Category goals
│   ├── loans/          # Loans with installments
│   ├── payees/         # Payee management
│   ├── reports/        # Net worth, cash flow, budget vs reality
│   ├── rules/          # Automation rules
│   ├── schedules/      # Recurring transactions
│   ├── settings_app/   # User + budget settings
│   ├── sync_engine/    # Offline sync (last-write-wins)
│   └── transactions/   # Transaction CRUD
├── core/               # Shared utilities
│   ├── helpers.py
│   ├── interest.py
│   └── templatetags/   # Custom filters (currency, etc.)
├── static/
│   ├── css/chele.css
│   └── img/logo.svg
├── templates/
│   └── base.html       # Sidebar layout, all views extend this
├── conftest.py         # Shared pytest fixtures
└── Dockerfile
```

## Key Decisions

- **UUID primary keys** on all models (offline-ready)
- **Spanish UI** hardcoded in templates (no i18n)
- **`apps/` package** for all Django apps
- **Builtin template filters** via `core.templatetags.chele_filters`
- **Currency format**: `$1,234.56` via `|currency` filter

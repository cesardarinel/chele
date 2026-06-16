# Chele 🐷

**Personal finance app implementing the YNAB (You Need A Budget) methodology.**  
Multi-user, multi-budget, zero-sum envelope budgeting with credit cards, loans, goals, CSV import, and collaboration.

> "Chele" is Argentine slang for money.  
> Named by [@cesardarinel](https://github.com/cesardarinel).

---

## Features

- **Multi-budget architecture** — independent budgets per family, business, personal, etc.
- **YNAB 4 Rules** — envelope budgeting: give every dollar a job, embrace true expenses, roll with the punches, age your money
- **Credit cards** — with daily interest calculation and payment tracking
- **Loans** — installment-based with automatic interest
- **Category goals** — monthly savings, target balance, target by date, true expenses
- **Recurring schedules** — auto-create transactions for bills and income
- **Automation rules** — auto-categorize transactions by payee, amount, or notes
- **CSV import** — column mapping, preview, duplicate detection
- **Reports** — net worth, cash flow, budget vs reality
- **Multi-device sync** — manual pull-to-refresh with last-write-wins
- **Collaboration** — invite members to share budgets
- **Dark mode** — YNAB-inspired color palette (#164E63 primary)
- **100% Spanish UI**

---

## Tech Stack

| Layer | Technology |
|-------|-----------|
| Backend | Python 3.13, Django 6.0 |
| Database | SQLite with WAL mode |
| Frontend | Django Templates, Tailwind CSS (CDN), vanilla JS |
| Deployment | Docker + Gunicorn + Whitenoise |
| Testing | pytest |

---

## Quick Start

### With Docker (recommended)

```bash
cp .env.example .env
# Edit SECRET_KEY in .env
docker compose up -d
```

Open http://localhost:8000. Register a new account and create your first budget.

### Without Docker

```bash
python -m venv .venv && source .venv/bin/activate
pip install -r requirements.txt
cp .env.example .env
# Edit SECRET_KEY in .env
python manage.py migrate
python manage.py collectstatic --noinput
python manage.py runserver
```

### Development mode

```bash
docker compose -f docker-compose.yml -f docker-compose.dev.yml up -d
```

This enables debug mode, auto-reload, and mounts the source code as a volume.

---

## Configuration

All configuration via environment variables (`.env`):

| Variable | Default | Description |
|----------|---------|-------------|
| `SECRET_KEY` | — | Django secret key (required) |
| `DEBUG` | `False` | Enable debug mode |
| `ALLOWED_HOSTS` | `localhost,127.0.0.1` | Comma-separated hosts |
| `PORT` | `8000` | Port for web server |
| `WORKERS` | `3` | Gunicorn worker count |

---

## Project Structure

```
chele/
├── chele/                  # Django project config
│   ├── settings.py
│   ├── urls.py             # All URLs in Spanish (/presupuestos/, /cuentas/, etc.)
│   └── context_processors.py
├── apps/                   # All Django applications
│   ├── accounts/           # Bank accounts & cash management
│   ├── budgets/            # Core: budgets, categories, YNAB logic
│   ├── credit_cards/       # Credit cards with interest
│   ├── csv_import/         # CSV transaction import
│   ├── goals/              # Category goals
│   ├── loans/              # Loans with installments
│   ├── payees/             # Payees/beneficiaries
│   ├── reports/            # Net worth, cash flow, budget vs reality
│   ├── rules/              # Auto-categorization rules
│   ├── schedules/          # Recurring transactions
│   ├── settings_app/       # User and budget settings
│   ├── sync_engine/        # Multi-device sync
│   └── transactions/       # Income/expense transactions
├── core/                   # Shared utilities
│   ├── helpers.py
│   └── interest.py         # Interest calculation (CC + Loans)
├── openspec/               # Design documentation
├── static/
│   ├── css/chele.css
│   └── img/cerdito.svg     # 🐷 mascot
└── templates/              # All Django HTML templates (36 files)
```

Key design decisions documented in `openspec/changes/finance-app-mvp/design.md`.

---

## User Guide

A complete Spanish user guide is available at [`CHILE_GUIA_DE_USO.md`](CHILE_GUIA_DE_USO.md).  
It covers first steps, the 4 YNAB rules, daily usage, multi-budget workflow, and FAQs.

---

## Design Docs

Full OpenSpec documentation lives under `openspec/`:

- [Proposal](openspec/changes/finance-app-mvp/proposal.md) — high-level scope
- [Design](openspec/changes/finance-app-mvp/design.md) — architecture, data model, UI, color system, mascot
- [Tasks](openspec/changes/finance-app-mvp/tasks.md) — 18 task groups, all completed
- [Specs](openspec/changes/finance-app-mvp/specs/) — 13 detailed specs per domain

---

## License

This project is private. All rights reserved.

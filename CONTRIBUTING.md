# Contributing to Chele

Thanks for your interest in contributing!

---

## Development Setup

```bash
git clone <repo-url> chele
cd chele
cp .env.example .env
# Edit SECRET_KEY in .env
docker compose -f docker-compose.yml -f docker-compose.dev.yml up -d
```

The dev setup mounts the source code as a volume so changes reflect instantly.

---

## Project Conventions

### Django apps

- All apps live under `apps/` as Django "contrib-style" apps (each has `urls.py`, `views.py`, `models.py`)
- Business logic lives in models or `core/`, not in views
- URLs are in Spanish (`/presupuestos/`, `/cuentas/`, `/transacciones/`)

### Frontend

- **No build step** — Tailwind CSS loaded via CDN in `templates/base.html`
- All templates extend `base.html`
- Colors use the `chele-*` custom Tailwind color palette (defined in `base.html` tailwind.config)
- Modals follow the pattern: `openModal(id)` / `closeModal(id)` with a hidden overlay
- Animations use CSS classes defined in `static/css/chele.css`

### Color palette

Defined in `templates/base.html` (tailwind.config):

```javascript
chele: {
    primary: '#164E63', 'primary-dark': '#0F3A48',
    success: '#16A34A',
    warning: '#D97706',
    danger: '#DC2626',
    neutral: '#9CA3AF',
    bg: '#0F172A', 'bg-secondary': '#1E293B', 'bg-tertiary': '#334155',
    sidebar: '#0B1121',
    surface: '#1E293B',
    text: '#F1F5F9', 'text-secondary': '#94A3B8', 'text-muted': '#64748B',
    border: '#334155', 'border-light': '#475569',
}
```

Use `chele-*` classes everywhere. Avoid hardcoded hex colors.

### Templates

- Use `bg-chele-surface` instead of `bg-white`
- Use `bg-chele-bg` for page backgrounds
- Use `text-chele-text` / `text-chele-text-secondary` / `text-chele-text-muted` for text hierarchy
- Use `border-chele-border` for borders
- Use `chele-primary` for all primary actions
- Inputs/selects should include `bg-chele-bg text-chele-text`
- Dropdowns/modals should use `bg-chele-surface`

### Tests

```bash
# Install test deps
pip install pytest pytest-django

# Run all tests
pytest

# Run specific test file
pytest auth_tests.py
```

---

## Commit Style

Use conventional commits:

- `feat:` — new feature
- `fix:` — bug fix
- `docs:` — documentation
- `refactor:` — code restructuring
- `style:` — formatting only
- `chore:` — config, dependencies, CI

Keep commits atomic. Each commit should be a logical unit.

---

## Pull Request Process

1. Create a feature branch from `master`
2. Make your changes
3. Run tests (`pytest`)
4. Submit a PR with a clear description
5. Ensure the UI works in both desktop and mobile viewports

# Currency Formatting

## Data Type

All monetary values use `DecimalField(max_digits=15, decimal_places=2)`.

| Model | Field |
|-------|-------|
| `Account` | `balance` |
| `CreditCard` | `balance`, `limit` |
| `Loan` | `total_amount`, `remaining_balance` |
| `Installment` | `amount` |
| `Transaction` | `amount` |

## Template Filter

`core/templatetags/chele_filters.py` registers `currency` as a builtin filter (available in all templates without `{% load %}`).

Usage: `${{ value|currency }}` → `$1,234.56`

- 2 decimal places
- Comma as thousands separator
- Dot as decimal separator
- Registered as a Django builtin via `OPTIONS.builtins` in settings.py

## Interest Rates

Displayed with raw `floatformat:2` (no currency filter, no thousands separator).

Usage: `{{ rate|floatformat:2 }}% anual` → `8.00% anual`

## Form Inputs

Input fields use raw `floatformat:2` (no commas, the value attribute must be parseable by the browser).

```html
<input type="number" step="0.01" value="{{ value|floatformat:2 }}">
```

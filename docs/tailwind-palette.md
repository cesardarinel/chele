# Tailwind Color Palette

YNAB-inspired dark mode. Defined as `chele-*` tokens in `tailwind.config`.

## Backgrounds

| Token | Hex | Usage |
|-------|-----|-------|
| `chele-bg` | `#0F172A` | Page background |
| `chele-bg-secondary` | `#1E293B` | Cards, panels |
| `chele-bg-tertiary` | `#334155` | Inputs, hover states |
| `chele-sidebar` | `#0B1121` | Sidebar background |

## Text

| Token | Hex | Usage |
|-------|-----|-------|
| `chele-text` | `#F1F5F9` | Primary text (white) |
| `chele-text-secondary` | `#94A3B8` | Secondary text |
| `text-gray-400` | `#9CA3AF` | Muted text (replaced `chele-text-muted` for contrast) |
| `text-amber-400` | `#FBBF24` | Links/actions on dark backgrounds |

## Semantic Colors

| Token | Hex | Usage |
|-------|-----|-------|
| `chele-primary` | `#164E63` | Buttons, focus rings |
| `chele-primary-dark` | `#0F3A48` | Button hover |
| `chele-primary-light` | `#1A6B84` | Lighter variant |
| `chele-success` | `#16A34A` | Positive balance |
| `chele-warning` | `#D97706` | Warnings |
| `chele-danger` | `#DC2626` | Negative/overspend |
| `chele-neutral` | `#9CA3AF` | Inactive |

## Borders

| Token | Hex |
|-------|-----|
| `chele-border` | `#334155` |
| `chele-border-light` | `#475569` |

## Accessibility

- Text on dark backgrounds **must** use `text-amber-400` (not `chele-primary #164E63` which is invisible on dark)
- Muted text uses `text-gray-400` (`#9CA3AF`, ~4.5:1 contrast) instead of `#64748B` (~2.2:1)

## Context

**chele** — Aplicación Django + SQLite que implementa la metodología completa de YNAB (4 reglas). Multi-usuario, múltiples presupuestos independientes con switch entre ellos. La app está representada por **Cerdito** (una mascota tipo alcancía) que guía al usuario. Toda la interfaz está 100% en español, responsive (mobile-first). Frontend con Django Templates + HTMX + Tailwind CSS.

---

## Arquitectura Django (recomendación comunidad)

### Estructura de Apps

Todas las apps viven dentro de `apps/` como un paquete Python (recomendación comunidad Django). Cada app representa un dominio de negocio con alta cohesión y bajo acoplamiento.

```
chele/
├── chele/                  # Config del proyecto (settings, urls root, wsgi)
│   ├── settings.py
│   ├── urls.py
│   ├── wsgi.py
│   └── context_processors.py
├── apps/                   # ← Todas las aplicaciones Django
│   ├── __init__.py
│   ├── accounts/           # Cuentas bancarias y efectivo
│   │   ├── models.py       # Account (checking, savings, cash)
│   │   ├── views.py
│   │   ├── urls.py
│   │   └── templates/accounts/
│   ├── budgets/            # Presupuesto, categorías, grupos + YNAB lógica
│   │   ├── models.py       # Budget, BudgetMembership, CategoryGroup, Category, MonthlyBudget
│   │   ├── views.py
│   │   ├── urls.py
│   │   ├── urls_registration.py
│   │   └── templates/budgets/
│   ├── credit_cards/       # Tarjetas de crédito con interés
│   │   ├── models.py       # CreditCard, InterestCharge
│   │   ├── views.py
│   │   ├── urls.py
│   │   └── templates/credit_cards/
│   ├── csv_import/         # Importación CSV
│   │   ├── views.py
│   │   ├── urls.py
│   │   └── templatetags/
│   ├── goals/              # Metas/objetivos por categoría
│   │   ├── models.py       # Goal
│   │   ├── views.py
│   │   └── urls.py
│   ├── loans/              # Préstamos con interés y cuotas
│   │   ├── models.py       # Loan, Installment
│   │   ├── views.py
│   │   ├── urls.py
│   │   └── templates/loans/
│   ├── payees/             # Beneficiarios
│   │   ├── models.py       # Payee
│   │   ├── views.py
│   │   └── urls.py
│   ├── reports/            # Reportes (sin modelos propios)
│   │   ├── views.py
│   │   └── urls.py
│   ├── rules/              # Reglas de automatización
│   │   ├── models.py       # Rule
│   │   ├── views.py
│   │   └── urls.py
│   ├── schedules/          # Programaciones
│   │   ├── models.py       # Schedule
│   │   ├── views.py
│   │   └── urls.py
│   ├── settings_app/       # Configuración (sin modelos propios)
│   │   ├── views.py
│   │   └── urls.py
│   ├── sync_engine/        # Sync multi-dispositivo
│   │   ├── models.py       # SyncLog
│   │   ├── views.py
│   │   └── urls.py
│   └── transactions/       # Transacciones (ingresos/egresos)
│       ├── models.py       # Transaction
│       ├── views.py
│       ├── urls.py
│       └── templates/transactions/
├── core/                   # Utilidades compartidas (sin modelos)
│   ├── interest.py         # Cálculo de intereses (TC y préstamos)
│   └── helpers.py          # Helpers: get_active_budget, current_month_year
├── static/
│   ├── css/chele.css
│   └── img/cerdito.svg
├── templates/
│   ├── base.html           # Base responsive con sidebar colapsable
│   ├── registration/
│   ├── accounts/
│   ├── budgets/
│   ├── credit_cards/
│   ├── csv_import/
│   ├── goals/
│   ├── loans/
│   ├── payees/
│   ├── reports/
│   ├── rules/
│   ├── schedules/
│   └── settings_app/
├── conftest.py             # Fixtures de pytest compartidos
├── auth_tests.py
├── CHILE_GUIA_DE_USO.md
├── Dockerfile
├── docker-compose.yml
├── manage.py
└── requirements.txt
```

### Principios de diseño
- **Apps dentro de `apps/`**: todas las aplicaciones Django como subpaquetes de `apps` (evita contaminar el root del proyecto)
- **`core/` como módulo compartido**: lógica reutilizable que no pertenece a ninguna app en particular (helpers, interés)
- **Fat models, thin views**: lógica de negocio en modelos/services
- **Templates con lógica mínima**: filtros y tags personalizados
- **URLs sin namespaces innecesarios**: cada app define sus URLs, importadas desde `apps.<name>.urls`
- **Tests por app**: cada app tiene su `tests.py` + `conftest.py` global para fixtures compartidos
- **Mobile-first responsive**: Tailwind con breakpoints `sm:` `md:` `lg:`
- **Sin dependencias JS externas**: vanilla JS para interactividad (sidebar, modales)

---

## Goals / Non-Goals

**Goals:**
- Múltiples presupuestos independientes con switch
- Metas/objetivos por categoría (Monthly Savings Goal, Target Balance, Target by Date)
- True Expenses (gastos anuales en cuotas mensuales)
- Las 4 reglas YNAB implementadas
- **Préstamos** como entidad separada con tasa de interés, cuotas, saldo pendiente
- **Tarjetas de crédito** con cálculo de intereses por falta de pago
- **Intereses acumulados diarios** para TC y préstamos
- **Cerdito** como mascota/guía de la aplicación (logo, mensajes, animaciones)
- **Responsive** mobile-first
- Sync manual pull-to-refresh con last-write-wins
- Importación CSV
- Despliegue con Docker

**Non-Goals:**
- API pública para terceros
- WebSockets / tiempo real
- Aplicación móvil nativa
- Migración a Postgres
- API en Go (post-MVP)

---

## Modelo de Datos

```
Usuario
  └── Presupuesto (múltiples, con switch)
        ├── Miembros (usuarios invitados con permisos)
        ├── Grupo de Categorías
        │     └── Categoría
        │           └── Meta (Monthly Savings Goal | Target Balance | Target by Date | True Expense)
        ├── Mes Presupuestal (año/mes)
        │     └── Asignación (categoría, monto presupuestado)
        ├── Cuenta (checking | savings | cash)
        │     ├── on_budget | off_budget
        │     └── Transacción
        ├── Tarjeta de Crédito (credit_card)
        │     ├── límite, tasa_interés, fecha_cierre, fecha_pago
        │     ├── InterésAcumulado (fecha, monto, diario)
        │     └── Transacción
        ├── Préstamo (loan)
        │     ├── tipo (personal, hipotecario, automotor, etc.)
        │     ├── monto_total, tasa_interés, cuotas, fecha_inicio
        │     ├── CuotaProgramada (número, monto, fecha_vencimiento, pagada)
        │     ├── InterésAcumulado (fecha, monto, diario)
        │     └── Transacción
        ├── Beneficiario
        ├── Regla
        └── Programación
```

Cada presupuesto es completamente independiente.

---

## Cuentas vs Tarjetas de Crédito vs Préstamos

### Cuentas (Account)
| Campo | Descripción |
|---|---|
| `type` | `checking` (corriente), `savings` (ahorro), `cash` (efectivo) |
| `on_budget` | Si el saldo está disponible para presupuestar |
| `balance` | Saldo actual (cálculo automático desde transacciones) |

- No tienen interés asociado
- Son la fuente de dinero para presupuestar
- Se pueden transferir entre sí

### Tarjetas de Crédito (CreditCard)
| Campo | Descripción |
|---|---|
| `limit` | Límite de crédito |
| `interest_rate` | Tasa de interés mensual (ej: 0.08 = 8% mensual) |
| `closing_day` | Día de cierre de factura |
| `due_day` | Día de vencimiento |
| `balance` | Saldo actual (negativo = deuda) |
| `on_budget` | Siempre False (son off-budget) |

**Comportamiento:**
- Cuando se hace un gasto con TC, se mueve dinero de la categoría a "Pago TC"
- Si no se paga el total antes del vencimiento, se calcula interés diario:
  - `interés_diario = (saldo_pendiente * tasa_interés) / 30`
  - Se acumula diariamente hasta que se pague
- Los intereses se registran como transacciones automáticas

### Préstamos (Loan)
| Campo | Descripción |
|---|---|
| `type` | `personal`, `hipotecario`, `automotor`, `estudiantil`, `otros` |
| `status` | `active` (activo), `completed` (finalizado) |
| `total_amount` | Monto total del préstamo |
| `interest_rate` | Tasa de interés **anual** |
| `remaining_balance` | Saldo pendiente actual |
| `total_installments` | Total de cuotas |
| `paid_installments` | Cuotas pagadas |
| `installment_amount` | Valor de cada cuota |
| `start_date` | Fecha de inicio |
| `next_due_date` | Próximo vencimiento |
| `account` | Cuenta de pago (FK) — las cuotas pagadas generan transacciones en esta cuenta |

**Comportamiento:**
- Se registra como pasivo (off-budget)
- Cada cuota tiene: número, monto, fecha_vencimiento, pagada
- Si no se paga a tiempo, se calcula interés diario sobre el saldo pendiente
- **Al pagar una cuota**: si tiene `account` asignada, crea una transacción de gasto en esa cuenta y reduce el saldo (visible en el presupuesto)
- **Al saldar** (`remaining_balance <= 0`): cambia `status` a `completed` y pasa a la sección "Historial"
- **Edición**: permite cambiar saldo actual, cuotas pagadas, valor cuota, total cuotas (agrega/elimina registros automáticamente)
- **Eliminación**: solo si no tiene transacciones asociadas

### Intereses (core/interest.py)
```
calcular_interes_diario(saldo, tasa_anual, dias_desde_vencimiento)
  → tasa_diaria = tasa_anual / 12 / 30
  → interés = saldo * tasa_diaria * dias

aplicar_interes(tarjeta_o_prestamo)
  → crea transacción automática por el interés acumulado
  → actualiza el saldo
```

---

## Las 4 Reglas YNAB

### Regla 1: Darle trabajo a cada peso
- Envelope budgeting: asignar cada peso disponible a una categoría
- "Por asignar" debe quedar en cero
- Over-allocation prevenida por el sistema

### Regla 2: Aceptar tus gastos reales
- True Expenses: gastos que ocurren anual/trimestral → divididos en cuotas mensuales
- Metas: Monthly Savings Goal, Target Balance, Target by Date
- Tasas de interés en TC y préstamos son gastos reales que deben presupuestarse

### Regla 3: Patear los golpes
- Mover dinero entre categorías
- Rollover de gastos al próximo mes
- Los intereses por falta de pago entran acá

### Regla 4: Envejecer tu dinero
- Hold for next month
- A Month Ahead
- Auto-hold en categorías de ingreso

---

## Metas / Objetivos

| Tipo de Meta | Descripción | Comportamiento |
|---|---|---|
| Monthly Savings Goal | Ahorrar X por mes | Asigna X automáticamente cada mes |
| Target Balance | Alcanzar saldo X | Asigna lo necesario hasta llegar a X |
| Target by Date | Alcanzar X para una fecha | Calcula cuota mensual: (X - actual) / meses restantes |
| True Expense | Gasto anual X (ej: seguro $1200/año) | Asigna X/12 cada mes |

---

## Decisiones

| Decisión | Opción elegida | Alternativas | Razón |
|---|---|---|---|
| Frontend | Django Templates + HTMX | React/Vue SPA | Mantiene el monolito simple, no requiere API REST compleja ni build step |
| Responsive | Tailwind responsive + sidebar colapsable en mobile | Sidebar fijo | Mobile-first: en < md el sidebar se oculta con menú hamburguesa |
| Sync | Tabla `sync_log` con `updated_at` | CRDT, OT | Last-write-wins es simple y suficiente para MVP |
| CSV Import | Inline en request | Celery async | Para MVP alcanza; si se vuelve lento se migra a Celery |
| IDs | UUIDs | autoincrement integers | UUIDs permiten generar IDs offline (esencial para sync) |
| TC con interés | Modelo separado CreditCard + InterestCharge | Account con tipo credit_card + interés | Separación clara de responsabilidades, cálculo de interés diario |
| Préstamos | App `loans` separada con modelo Loan + Installment + Interest | Cuenta tipo loan | Los préstamos tienen comportamiento distinto (cuotas, amortización, interés) |
| Intereses | `core/interest.py` reutilizable | Lógica duplicada en cada app | Single source of truth para cálculo financiero |
| Mascota | **Cerdito** como logo + mensajes + animaciones | Sin mascota | Identidad de marca, hace la app más amigable |
| UI Layout | Sidebar oscuro YNAB + responsive | Sidebar fijo | YNAB look & feel, adaptable a mobile |
| Estilo visual | YNAB palette + Cerdito branding | Indigo | Consistencia con metodología |
| Terminología | YNAB español: "Por asignar", "Actividad", "Disponible" | ActualBudget | Matching YNAB UX |
| Idioma | 100% español hardcodeado en templates | Django i18n | MVP simple; post-MVP se puede migrar a i18n |
| CSS Framework | Tailwind CSS | Bootstrap, vanilla | Tailwind moderno y fácil de mantener |

---

## UI Views (responsive)

### Desktop (≥ md)
```
┌──────────┬──────────────────────────────────────┐
│ Sidebar  │  Main Content Area                    │
│ #141A26  │                                      │
│          │  [Por asignar: $0]                    │
│ Chele    │  ┌──────────┬──────┬────────┬──────┐ │
│ [Presu ▼]│  │Categoría │Presu.│Activ.  │Disp. │ │
│          │  ├──────────┼──────┼────────┼──────┤ │
│ ● Presup │  │Comida    │ 500  │ -300   │ 200  │ │
│ ● Report │  └──────────┴──────┴────────┴──────┘ │
│ ● Prog.  │                                      │
│ ● Benef. │  [Cerdito te dice: "¡Vas bien!"]     │
│ ● Reglas │                                      │
│ ──────── │                                      │
│ Cuentas  │                                      │
│ [+ Nva]  │                                      │
│ ├ Banco  │                                      │
│ └ Visa   │                                      │
│ ──────── │                                      │
│ ⚙ Config │                                      │
│ 🚪 Salir │                                      │
└──────────┴──────────────────────────────────────┘
```

### Mobile (< md)
```
┌──────────────────────────────────┐
│ ☰ Chele        [Por asignar: $0] │  ← Top bar con hamburguesa
├──────────────────────────────────┤
│  [Presupuesto ▼]   may 2026      │  ← Selector arriba
├──────────────────────────────────┤
│ Categoría  │Presu│Activ│Disp     │
│────────────┼─────┼────┼─────────│
│ Comida     │ 500 │-300│  200     │
│ Transporte │ 200 │-50 │  150     │
│ ...                                │
├──────────────────────────────────┤
│ 🐷 ¡Bien ahí!                    │  ← Cerdito mensaje
└──────────────────────────────────┘
```

Sidebar colapsable con overlay en mobile. Se abre con ☰, se cierra con ✕ o tocando fuera.

### Vistas principales

| Vista | Columnas | Flujo | Mobile |
|---|---|---|---|
| **Presupuesto** | Categoría \| Presupuestado \| Actividad \| Disponible | Asignar → Gastar → Mover → Reservar | Tabla scroll horizontal en mobile |
| **Registro de Cuenta** | Fecha \| Beneficiario \| Categoría \| Monto | Importar → Categorizar → Conciliar | Cards en vez de tabla en mobile |
| **TC** | Límite, tasa, fecha cierre/pago, saldo, intereses | Gastar → Acumula interés si no paga → Pagar | Resumen + detalle |
| **Préstamos** | Tipo, total, cuotas, tasa, saldo, próximo vencimiento | Solicitar → Pagar cuota → Interés si atrasa | Resumen + detalle |
| **Reportes** | Patrimonio Neto, Cash Flow, Budget vs Reality | Click → pantalla completa | Scroll vertical |
| **Programaciones** | Lista de transacciones recurrentes | Crear → Auto-ejecutar | Cards |
| **Beneficiarios** | Lista con fusión | Gestionar | Lista simple |
| **Reglas** | Condiciones → Acciones | Automatizar | Cards |
| **Configuración** | Perfil + presupuesto + miembros | Administrar | Stack vertical |

---

## Paleta de Colores — Chele Cerdito

### Nueva paleta unificada

```
CHILE COLOR SYSTEM v2
═══════════════════════════════════════════════════

Primary (Cerdito Gold):     #F59E0B
Primary Dark:               #D97706
Primary Light:              #FDE68A

Secondary (Cerdito Pink):   #EC4899
Secondary Dark:             #DB2777

Dark (Sidebar):             #1C1917  (warm dark, reemplaza #141A26)
Dark Text:                  #F5F5F4
Dark Hover:                 rgba(255,255,255,0.08)

Background:                 #FFFBEB  (warm cream, reemplaza #F7F8FA)
Surface (cards):            #FFFFFF
Border:                     #E7E5E4

Success:                    #10B981  (emerald)
Danger:                     #EF4444  (red)
Warning:                    #F59E0B  (amber)

Text Primary:               #1C1917
Text Secondary:             #78716C
Text Muted:                 #A8A29E

Gradients:
  Sidebar: linear-gradient(180deg, #1C1917 0%, #292524 100%)
  Hero:    linear-gradient(135deg, #F59E0B 0%, #EC4899 100%)
```

### Mapeo de roles

| Rol | Color | Hex | Uso |
|---|---|---|---|
| Sidebar bg | Warm Dark | `#1C1917` | Barra lateral, fondo oscuro |
| Sidebar hover | White 8% | `rgba(255,255,255,0.08)` | Hover items sidebar |
| Primary | Gold | `#F59E0B` | Botones principales, links, acentos |
| Primary hover | Dark Gold | `#D97706` | Hover de botones primarios |
| Secondary | Pink | `#EC4899` | Acentos secundarios, badges |
| Background | Warm Cream | `#FFFBEB` | Fondo general de la app |
| Success | Emerald | `#10B981` | Montos positivos, disponible |
| Danger | Red | `#EF4444` | Montos negativos, deudas, intereses |
| Warning | Amber | `#F59E0B` | Alertas, Cerdito mensajes |

---

## Cerdito — Mascota de la App

### Sistema de iconos

```
Logo oficial: SVG personalizado (no emoji)
═══════════════════════════════════════════════════

Cerdito Feliz    😊 🐷 → SVG con sonrisa, mejillas rosadas
  Uso: login, onboarding, saldo positivo, empty states
  Animación: bounce suave

Cerdito Preocupado 😟 🐷 → SVG con ceño fruncido
  Uso: alertas, deudas, gasto excesivo
  Animación: wobble/temblor

Cerdito Ahorrando  🪙 🐷 → SVG con moneda en la mano
  Uso: metas de ahorro, true expenses
  Animación: coin drop

Cerdito Dormido    😴 🐷 → SVG con burbuja Zzz
  Uso: empty states, "sin actividad"
  Animación: respiración suave
```

### Presencia en la app

| Lugar | Icono | Mensaje |
|---|---|---|
| Login | Cerdito Feliz + logo | "Tomá el control de tu dinero" |
| Sidebar | Cerdito Feliz + "Chele" | Logo siempre visible |
| Empty states | Cerdito Dormido | "🐷 No hay transacciones todavía" |
| Por asignar > 0 | Cerdito Feliz | "🐷 ¡Tenés dinero por asignar!" |
| Deuda/Interés | Cerdito Preocupado | "🐷 Cuidado con los intereses" |
| Meta cumplida | Cerdito Ahorrando | "🐷 ¡Meta alcanzada!" |
| 404 | Cerdito Preocupado | "🐷 Esta página no existe" |
| Error | Cerdito Preocupado | "🐷 Algo salió mal" |
| Loading/Sync | Cerdito (spinning) | "🐷 Sincronizando..." |

---

## Animaciones

```
ANIMATION SYSTEM
═══════════════════════════════════════════════════

Transiciones de página:   fade-in + slide-up (300ms)
Sidebar mobile:           slide-left/right (300ms ease)
Modal:                    fade-in + scale (200ms)
Cerdito:                  bounce (1s infinite)
Números/balances:         count-up (500ms)
Hover cards:              shadow + translateY (-2px)
Hover botones:            scale (1.02) + shadow
Loading:                  skeleton screens + spin coin

CSS Classes:
  animate-bounce-slow     → 2s bounce (Cerdito)
  animate-fade-in         → 300ms fade
  animate-slide-up        → 300ms slide + fade
  animate-count-up        → number animation
  animate-coin-drop       → coin falling animation
  animate-pulse-glow      → glowing pulse for alerts
  animate-shake           → shake for errors
  animate-bounce          → Tailwind default
```

| Elemento | Animación | Duración | Trigger |
|---|---|---|---|
| Cerdito logo | Bounce | 2s | Page load |
| Cerdito mensajes | Slide up + fade | 300ms | Aparece en DOM |
| Sidebar mobile | Translate X | 300ms | Click hamburguesa |
| Modal | Scale + fade | 200ms | Open/close |
| Números | Count up | 500ms | Page load / update |
| Cards hover | Translate Y | 150ms | Hover |
| Botones | Scale | 100ms | Hover |
| Loading | Pulse glow | 1.5s | Mientras carga |
| Sync | Spin | 1s | Durante sync |
| Error | Shake | 300ms | En error |

---

## Cálculo de Intereses

```
interés_diario = saldo_pendiente * (tasa_anual / 12 / 30) * días_atraso

Ejemplo TC:
  Saldo: $10,000 | Tasa: 96% anual (8% mensual) | Días atraso: 15
  interés = 10000 * (0.96/12/30) * 15 = $400

Ejemplo Préstamo:
  Saldo: $50,000 | Tasa: 36% anual (3% mensual) | Días atraso: 10
  interés = 50000 * (0.36/12/30) * 10 = $500
```

Los intereses se registran como transacciones automáticas con categoría "Intereses" y beneficiario según corresponda (banco/entidad).

---

## Risks / Trade-offs

- [SQLite concurrencia] Django + SQLite con WAL mode soporta lecturas concurrentes bien, pero escrituras tienen lock → Mitigación: sync manual reduce presión
- [Sync conflictos] Last-write-wins puede perder datos → Aceptado para MVP
- [Crecimiento] SQLite tiene límites prácticos (~100GB) → Mitigación: monitorear, migrar a Postgres si es necesario
- [Intereses acumulados] El cálculo diario requiere un cron/task programado → Mitigación: celery-beat o cron de Django
- [Responsive] La tabla de presupuesto tiene muchas columnas en mobile → Mitigación: scroll horizontal + versión cards en pantallas muy chicas

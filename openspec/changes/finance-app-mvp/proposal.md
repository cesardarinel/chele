## Why

Crear **chele**, una aplicación de finanzas personales multi-usuario que implementa la metodología completa de YNAB (You Need A Budget) — la misma en la que se basa ActualBudget. Construida con Django + SQLite como stack inicial, 100% en español, permitiendo crear múltiples presupuestos independientes (familia, empresa, personal, etc.) con switch entre ellos, todo colaborativo con sync manual multi-dispositivo.

## What Changes

- Nuevo proyecto Django monolítico llamado **chele** con autenticación por email+contraseña
- Interfaz 100% en español replicando las vistas y formularios de ActualBudget/YNAB
- **Múltiples presupuestos independientes** (como "archivos" separados en ActualBudget): crear N presupuestos con switch entre ellos
- Cada presupuesto tiene sus propias: cuentas, categorías, transacciones, miembros, payees, reglas, programaciones
- Metas/objetivos por categoría (Monthly Savings Goal, Target Balance, Target by Date)
- True Expenses: gastos anuales divididos en cuotas mensuales
- Implementación completa de las 4 reglas YNAB
- Envelope budgeting mensual: asignar dinero a categorías, disponible a presupuestar
- Importación CSV de transacciones (onboarding y recurrente)
- Tratamiento especial para tarjetas de crédito
- Sync manual pull-to-refresh con last-write-wins
- Reportes: Budget vs Reality, Net Worth, Cash Flow
- Despliegue vía Docker

## Capabilities

### New Capabilities

- `user-auth`: Registro e inicio de sesión con email y contraseña, manejo de sesiones
- `budgets-management`: Creación de múltiples presupuestos independientes (familia, empresa, personal, etc.), switch entre ellos, invitación de miembros por presupuesto
- `envelope-budgeting`: Presupuesto mensual basado en sobres con las 4 reglas YNAB: categorías, asignación de fondos, seguimiento gasto vs presupuestado, rollover, hold for next month
- `category-goals`: Metas/objetivos por categoría: Monthly Savings Goal, Target Balance, Target by Date, True Expenses (gasto anual en cuotas)
- `transactions`: Registro manual de transacciones (ingresos y egresos), vinculación a cuentas y categorías
- `credit-cards`: Manejo especial de tarjetas de crédito siguiendo metodología YNAB
- `csv-import`: Importación de transacciones desde archivos CSV, tanto inicial como recurrente
- `sync-engine`: Sincronización manual pull-to-refresh entre dispositivos con resolución last-write-wins
- `accounts`: Gestión de cuentas (efectivo, bancarias, tarjetas) con saldos, on/off budget
- `reports`: Reportes Budget vs Reality, Net Worth, Cash Flow
- `ui-views`: Vistas y formularios replicando la UI de ActualBudget/YNAB
- `i18n-es`: Interfaz 100% en español

### Modified Capabilities

- Ninguna (proyecto nuevo)

## Impact

- Nuevo repositorio Django + SQLite
- Modelo de datos: Usuario → Budget (múltiple) → Cuentas/Categorías/Transacciones/Miembros
- Categorías con soporte de metas (objetivos mensuales, saldo objetivo, fecha objetivo)
- Lógica de True Expenses (gasto anual distribuido en cuotas mensuales)
- Switch de presupuesto en la barra lateral (selector de presupuesto activo)
- Docker y docker-compose para desarrollo y producción

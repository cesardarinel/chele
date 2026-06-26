## Why

Chele implementa el método de sobres (envelope budgeting) pero carece de las características clave que hacen a YNAB una herramienta de planificación proactiva en lugar de un simple rastreador de gastos. Los usuarios no pueden establecer metas automáticas (targets), no tienen visibilidad de "Ready to Assign", no hay código semafórico (verde/amarillo/rojo) para el estado del presupuesto, y no existe un flujo de revisión de transacciones. Sin estas funcionalidades, el presupuesto sigue siendo reactivo. Este cambio implementa el motor YNAB completo para transformar Chele en una herramienta de ingeniería financiera personal proactiva.

## What Changes

- **Ready to Assign**: Nueva sección superior que muestra fondos disponibles para distribuir. Cada asignación reduce este saldo. Zero-sum enforcement: no se puede asignar más de lo disponible.
- **Targets (Metas)**: Sistema de objetivos por categoría: mensual, anual, target_balance, target_date, true_expense. Con indicadores visuales verde (cumplido) / amarillo (subfinanciado) / rojo (sobrepasado).
- **Refill up to vs Set aside another**: Dos modos de meta. "Refill up to" completa hasta un límite (considerando remanentes con rollover). "Set aside another" pide el monto completo cada mes.
- **Snooze Target**: Posibilidad de pausar una meta por un mes específico.
- **Overspending handling**: Distinción entre sobregasto en efectivo (rojo, requiere corrección inmediata) y sobregasto en tarjeta de crédito (naranja, auto-mueve fondos disponibles a la categoría de pago).
- **Cover overspending**: Botón/flujo para cubrir sobregastos desde otra categoría.
- **Auto-assign (Underfunded)**: Asignación automática de fondos a categorías subfinanciadas con un clic.
- **Inspector Panel**: Panel derecho en escritorio con detalle de la categoría seleccionada: metas, promedios de gasto, acciones rápidas.
- **Spotlight Mode**: Centro de notificaciones con flujo de revisión de transacciones importadas.
- **Movile navigation**: Barra inferior con tres pestañas (Plan, Accounts, Spending).
- **Drag-and-drop**: Reordenamiento de categorías y grupos por arrastre.
- **Cost to be me / Reality Check**: Suma de todos los targets mensuales vs ingresos esperados, con alerta si los targets superan los ingresos.
- **Código semafórico**: Verde (OK), Amarillo (subfinanciado/TC), Rojo (sobregasto efectivo) en toda la UI.
- **Layout tripartito escritorio**: Panel izquierdo (cuentas), central (presupuesto con columnas Assigned/Activity/Available), derecho (Inspector).
- **Rollover automático**: Los remanentes no asignados se transfieren al siguiente mes.

## Capabilities

### New Capabilities
- `ready-to-assign`: Sección superior con saldo disponible para distribuir y zero-sum enforcement.
- `targets`: Sistema de metas por categoría (monthly, yearly, target_balance, target_date, true_expense) con indicadores visuales.
- `refill-up-to`: Modo de meta que completa hasta un límite considerando rollover de remanentes.
- `snooze-target`: Pausar temporalmente una meta por un mes.
- `overspending`: Manejo de sobregastos con distinción efectivo (rojo) vs tarjeta de crédito (naranja) y auto-movimiento a payment.
- `cover-overspending`: Flujo para cubrir sobregastos desde otra categoría.
- `auto-assign`: Asignación automática de fondos a categorías subfinanciadas.
- `inspector-panel`: Panel derecho con detalle de categoría (metas, promedios, acciones).
- `spotlight-mode`: Centro de notificaciones con flujo de revisión de transacciones.
- `mobile-navigation`: Barra inferior con pestañas Plan/Accounts/Spending.
- `drag-drop-categories`: Reordenamiento de categorías y grupos por arrastre.
- `cost-to-be-me`: Suma de targets mensuales vs ingresos, con alerta de desequilibrio.
- `traffic-light-ui`: Código semafórico (verde/amarillo/rojo) en toda la interfaz.
- `tripartite-layout`: Layout escritorio con tres paneles (cuentas, presupuesto, inspector).
- `rollover`: Transferencia automática de remanentes al siguiente mes.

### Modified Capabilities
- `budget-view`: Las columnas pasan de "Presupuestado/Actividad/Disponible" a "Assigned/Activity/Available" con código semafórico.
- `account-detail`: Integración con Spotlight para revisión de transacciones importadas.
- `monthly-budget`: El modelo MonthlyBudget soporta targets, rollover y cálculo de underfunded.

## Impact

- **apps/budgets/**: Nuevos modelos (Target, TargetType), modificaciones a MonthlyBudget para soportar rollover y underfunded. Nuevas vistas para auto-assign, cover, snooze.
- **apps/transactions/**: Flujo de revisión (Spotlight), manejo de sobregasto en tarjetas de crédito.
- **templates/**: Nuevo layout tripartito, barra inferior móvil, inspector panel, tarjetas con código semafórico.
- **static/**: CSS para drag-and-drop, slider, semáforo, animaciones.
- **API Go**: Endpoints para targets, auto-assign, cover, spotlight review.

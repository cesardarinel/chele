## Why

Los usuarios nuevos llegan a Chele y ven una pantalla vacía con $0 en todos lados. No saben qué hacer, no entienden el método de sobres (YNAB), y abandonan. El concepto de "darle trabajo a cada peso" es contra-intuitivo al principio. Sin un onboarding interactivo que los guíe paso a paso —crear cuentas, asignar dinero, poner metas, programar ingresos— el usuario nunca experimenta el "click" mental de entender cómo funciona el presupuesto. Este cambio agrega un onboarding forzado paso a paso con coach marks, tooltips interactivos y zero-sum enforcement visual.

## What Changes

- **Onboarding step field**: Agregar `onboarding_step` (IntegerField, default=0) al modelo User. Controla el progreso del onboarding (0=no iniciado, 7=completado).
- **OnboardingMiddleware**: Middleware que detecta si `user.onboarding_step < 7` y si el usuario está en una ruta que no sea del onboarding. Si está incompleto, redirige a la vista actual pero con el overlay del onboarding superpuesto.
- **Overlay interactivo**: Capa semi-transparente que bloquea la interacción con el fondo. Solo el elemento destacado (coach mark) es clickeable. Tooltip con texto explicativo y botones "Siguiente"/"Atrás".
- **Barra de progreso**: Barra inferior flotante mostrando paso actual / total.
- **7 pasos**: Bienvenida → Cuentas → Asignar dinero → Metas → Programaciones → Deudas → Completado.
- **Paso 3 obligatorio**: Ready to Assign debe llegar a $0 para avanzar. El usuario asigna desde la UI real, el onboarding detecta el cambio y avanza automáticamente.
- **Coach marks**: Tooltips con flecha que apuntan a elementos específicos de la UI (botón "+ Cuenta", "Por asignar", categorías, nav items).
- **Indicadores visuales post-onboarding**: Tooltips en "Por asignar" que explican si el dinero tiene trabajo o no. Barras de progreso en metas. Desglose de rollover vs nuevo en disponible.

## Capabilities

### New Capabilities
- `onboarding-step-model`: Campo `onboarding_step` en User + migración
- `onboarding-middleware`: Middleware que enforcea el flujo de onboarding
- `onboarding-overlay`: Overlay interactivo con coach marks y tooltips
- `onboarding-step-welcome`: Paso 1 — bienvenida y explicación del método
- `onboarding-step-accounts`: Paso 2 — crear primera cuenta
- `onboarding-step-assign`: Paso 3 — asignar todo el dinero (zero-sum forzado)
- `onboarding-step-goals`: Paso 4 — configurar metas/targets
- `onboarding-step-schedules`: Paso 5 — programar ingresos/gastos recurrentes
- `onboarding-step-debts`: Paso 6 — agregar deudas (TC y préstamos)
- `onboarding-step-done`: Paso 7 — resumen y finalización
- `onboarding-coach-js`: JavaScript para coach marks, detección de acciones, polling
- `onboarding-post-completion-ui`: Indicadores visuales que quedan después del onboarding (explicación de RTA, progreso de metas, desglose disponible)

## Impact

- **chele/settings.py**: Agregar `OnboardingMiddleware` a MIDDLEWARE
- **chele/urls.py**: Incluir `onboarding.urls`
- **chele/context_processors.py**: Inyectar `onboarding_step`, `onboarding_active`
- **templates/base.html**: Incluir overlay si `onboarding_active`
- **apps/onboarding/***: Nueva app con middleware, views, templates
- **static/js/onboarding.js**: Coach marks, polling, overlay logic
- **static/css/onboarding.css**: Estilos del overlay, tooltips, animaciones
- **User model**: +1 field (`onboarding_step`)

## 1. Modelo y migración

- [ ] 1.1 Agregar `onboarding_step` IntegerField(default=0) al modelo User via migración de datos
- [ ] 1.2 Crear migración que setea `onboarding_step=7` para todos los usuarios existentes
- [ ] 1.3 Agregar `onboarding_step` al admin de User

## 2. App de onboarding

- [ ] 2.1 Crear `apps/onboarding/` con `__init__.py`
- [ ] 2.2 Crear `apps/onboarding/middleware.py`: OnboardingMiddleware que inyecta `onboarding_step` y `onboarding_active` en el contexto
- [ ] 2.3 Crear `apps/onboarding/views.py`: vista `onboarding_state` (GET /onboarding/state → JSON con step, completed, rta)
- [ ] 2.4 Crear `apps/onboarding/urls.py`: ruta `onboarding/state/`
- [ ] 2.5 Registrar middleware en `chele/settings.py`
- [ ] 2.6 Incluir urls en `chele/urls.py`
- [ ] 2.7 Actualizar `chele/context_processors.py` para incluir `onboarding_active` y `onboarding_step`

## 3. Templates del overlay

- [ ] 3.1 Crear `templates/onboarding/overlay.html` — contenedor principal del overlay + progreso
- [ ] 3.2 Crear `templates/onboarding/step_1_welcome.html` — bienvenida
- [ ] 3.3 Crear `templates/onboarding/step_2_accounts.html` — crear cuenta
- [ ] 3.4 Crear `templates/onboarding/step_3_assign.html` — asignar dinero
- [ ] 3.5 Crear `templates/onboarding/step_4_goals.html` — metas
- [ ] 3.6 Crear `templates/onboarding/step_5_schedules.html` — programaciones
- [ ] 3.7 Crear `templates/onboarding/step_6_debts.html` — deudas
- [ ] 3.8 Crear `templates/onboarding/step_7_done.html` — completado
- [ ] 3.9 Integrar overlay en `base.html` antes de `</body>`

## 4. JavaScript del onboarding

- [ ] 4.1 Crear `static/js/onboarding.js` con:
  - Overlay show/hide y gestión de pointer-events
  - Posicionamiento de coach marks (tooltip con flecha)
  - Pulse animation en elemento destacado
  - Polling cada 3s a `/onboarding/state`
  - Auto-avance cuando condición se cumple
  - Submit event detection para verificación inmediata
  - Mobile adaptation (tooltip → bottom sheet)
  - Keyboard navigation

## 5. CSS del onboarding

- [ ] 5.1 Crear `static/css/onboarding.css` con:
  - Estilos del overlay (fondo semi-transparente, z-index)
  - Tooltip con flecha (posición absoluta, arrow via pseudo-element)
  - Pulse animation para elemento destacado
  - Bottom sheet styles para mobile
  - Progress bar styles
  - Transiciones (fade-in, slide-up)

## 6. Paso 1: Bienvenida

- [ ] 6.1 Implementar template step_1_welcome.html
- [ ] 6.2 Tooltip centrado con logo y explicación del método YNAB
- [ ] 6.3 Botón "Empezar →" que avanza step a 2

## 7. Paso 2: Crear cuenta

- [ ] 7.1 Implementar template step_2_accounts.html
- [ ] 7.2 Coach mark apuntando a "+ Cuenta" en sidebar
- [ ] 7.3 Tooltip explicando saldo inicial y RTA
- [ ] 7.4 Polling detecta cuando `Account.objects.count() >= 1`
- [ ] 7.5 Botón "Ya agregué todo →" para saltar

## 8. Paso 3: Asignar dinero (zero-sum forzado)

- [ ] 8.1 Implementar template step_3_assign.html
- [ ] 8.2 Coach mark apuntando a sección "Por asignar"
- [ ] 8.3 Tooltip explicando zero-sum
- [ ] 8.4 Polling detecta cuando `ready_to_assign == 0`
- [ ] 8.5 Botón "Siguiente" deshabilitado mientras RTA > 0
- [ ] 8.6 Mostrar progreso: "$X / $Y asignado"
- [ ] 8.7 Animación de éxito cuando RTA = $0
- [ ] 8.8 NO skippable

## 9. Paso 4: Metas

- [ ] 9.1 Implementar template step_4_goals.html
- [ ] 9.2 Coach mark apuntando a una categoría
- [ ] 9.3 Polling detecta cuando `Goal.objects.count() >= 1`
- [ ] 9.4 Botón "Omitir" para saltar

## 10. Paso 5: Programaciones

- [ ] 10.1 Implementar template step_5_schedules.html
- [ ] 10.2 Coach mark apuntando a nav "Programaciones"
- [ ] 10.3 Polling detecta cuando `Schedule.objects.count() >= 1`
- [ ] 10.4 Botón "Omitir" para saltar

## 11. Paso 6: Deudas

- [ ] 11.1 Implementar template step_6_debts.html
- [ ] 11.2 Coach marks apuntando a "+ TC" y "+ Préstamo"
- [ ] 11.3 Polling detecta cuando `CreditCard.objects.count() >= 1 or Loan.objects.count() >= 1`
- [ ] 11.4 Botón "No tengo deudas →" para saltar

## 12. Paso 7: Completado

- [ ] 12.1 Implementar template step_7_done.html
- [ ] 12.2 Mostrar resumen de lo realizado en pasos anteriores
- [ ] 12.3 Tips post-onboarding
- [ ] 12.4 Botón "Ir al presupuesto" → setea step=7 y redirige

## 13. Post-onboarding UI indicators

- [ ] 13.1 Agregar tooltip en "Por asignar" que explica qué significa
- [ ] 13.2 Agregar breakdown de disponible en inspector panel: "($X del mes pasado + $Y nuevo)"
- [ ] 13.3 Agregar labels en sidebar: "💰 En presupuesto" / "🏦 Ahorro"
- [ ] 13.4 Mostrar "✅ Todo asignado" cuando RTA = $0
- [ ] 13.5 Barra de progreso de meta en inspector panel

## 14. Tests

- [ ] 14.1 Test: onboarding_step se crea como 0 para usuarios nuevos
- [ ] 14.2 Test: usuarios existentes migran con step=7
- [ ] 14.3 Test: middleware inyecta contexto correcto
- [ ] 14.4 Test: endpoint /onboarding/state devuelve JSON correcto
- [ ] 14.5 Test: overlay no se muestra cuando step=7

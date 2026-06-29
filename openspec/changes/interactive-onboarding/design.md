## Context

Chele actualmente tiene un funnel de registro que termina en una pantalla de presupuesto completamente en cero (sin cuentas, sin dinero asignado, sin guía). Los usuarios nuevos no entienden el método YNAB de "darle trabajo a cada peso" porque la UI no lo explica ni lo fuerza. No existe un indicador visual de qué dinero está planificado vs qué está "sin trabajo". No hay onboarding, no hay coach marks, no hay zero-sum enforcement visible.

## Goals / Non-Goals

**Goals:**
- Agregar `onboarding_step` al modelo User (0-7)
- Middleware que inyecta el overlay cuando step < 7
- Overlay interactivo semi-transparente con coach marks y tooltips
- 7 pasos: Bienvenida → Cuentas → Asignar (forzado) → Metas → Programaciones → Deudas → Completado
- Paso 3 (asignar) es obligatorio: RTA debe llegar a $0
- Detección automática de acciones del usuario via polling
- Indicadores post-onboarding persistentes (tooltips en RTA, breakdown de disponible, progreso de metas)
- Mobile responsive: tooltips → bottom sheets

**Non-Goals:**
- No se modifica el flujo de registro existente (solo se agrega overlay después)
- No se eliminan funcionalidades existentes
- No se agregan librerías externas (solo vanilla JS + Tailwind)

## Decisions

1. **Campo en User vs sesión**: Se usa `User.onboarding_step` (DB) en vez de sesión para que persista entre dispositivos y sesiones. Los usuarios existentes se migran con `step=7` (ya onboardeados).

2. **Middleware vs decorador**: Se usa middleware en vez de decorador por view para cubrir todas las rutas sin modificar cada vista. El middleware solo inyecta variables de contexto, no redirige (excepto si step=0 después de registro).

3. **Detección por polling vs webhook**: Se usa polling cada 3s a `GET /onboarding/state` en vez de webhooks/eventos porque es más simple de implementar y el intervalo de 3s es aceptable para UX. El endpoint devuelve `{step, step_completed, ready_to_assign}`.

4. **Overlay vs wizard separado**: Se usa overlay sobre la UI real en vez de un wizard de páginas separadas porque: (a) el usuario aprende interactuando con la interfaz real, (b) no hay que mantener rutas/views separadas, (c) la transición post-onboarding es natural (el overlay desaparece y la UI está igual).

5. **Coach marks con selectores CSS**: Cada paso define un `target_selector` (ej: `a[href="/cuentas/crear/"]`) que JS usa para posicionar el tooltip. Si el selector no existe en la página actual, se muestra un tooltip centrado.

6. **Polling endpoint on `/onboarding/state`**: Endpoint que recibe GET y devuelve JSON con el estado del paso actual. El JS decide si avanzar basado en la respuesta.

7. **Post-onboarding indicators**: Se implementan como tooltips HTML nativos (title attribute mejorados con CSS) y modificaciones menores al template (breakdown de disponible en inspector, labels en sidebar).

## Risks / Trade-offs

- [UX] El polling cada 3s puede sentir lento si el usuario asigna dinero y el paso no avanza inmediatamente. Mitigación: el JS también detecta el evento `submit` en los forms de asignación y fuerza una verificación inmediata.
- [Rendimiento] El polling es a un endpoint liviano (solo queries COUNT). Riesgo bajo.
- [Complejidad JS] El overlay con pointer-events management puede tener bugs en edge cases. Mitigación: pruebas manuales en Chrome/Safari/Firefox.
- [Mobile] Los coach marks con posicionamiento absoluto pueden no funcionar bien en pantallas pequeñas. Mitigación: en mobile, los tooltips siempre son bottom sheets (sin posicionamiento relativo a elementos).

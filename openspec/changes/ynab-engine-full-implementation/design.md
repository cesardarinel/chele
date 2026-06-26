## Context

Chele actualmente implementa un sistema de presupuesto básico con asignación manual por categoría y cálculo de "available to budget". No existe concepto de targets/metas, código semafórico, ready-to-assign, auto-assign, manejo de sobregastos diferenciado (efectivo vs TC), ni panel Inspector. La UI es funcional pero no sigue el estándar YNAB de "darle trabajo a cada peso". La base de datos ya tiene modelos Account, Transaction, MonthlyBudget, Category, CategoryGroup, Schedule y Goal. El Goal model actual es limitado (no soporta refill-up-to, snooze, ni cómputo de underfunded). No existe modelo Target.

## Goals / Non-Goals

**Goals:**
- Implementar Ready to Assign como sección superior del presupuesto con zero-sum enforcement
- Sistema de Targets (metas) con tipos monthly, yearly, target_balance, target_date, true_expense
- Modo Refill up to (completar hasta límite con rollover) y Set aside another (monto fijo cada mes)
- Snooze Target: pausar meta por un mes específico
- Código semafórico: Verde (OK), Amarillo (subfinanciado/TC), Rojo (sobregasto efectivo)
- Overspending handling: cash red, credit card orange, auto-move to payment
- Cover overspending: flujo para cubrir desde otra categoría
- Auto-assign (Underfunded): asignar a todas las categorías subfinanciadas con un clic
- Inspector Panel: panel derecho con detalle de categoría seleccionada
- Spotlight mode: centro de notificaciones con flujo Review para transacciones
- Layout tripartito en escritorio (cuentas | presupuesto | inspector)
- Barra inferior de navegación en móvil (Plan | Accounts | Spending)
- Drag-and-drop para reordenar categorías/grupos
- Cost to be me / Reality Check: suma targets mensuales vs ingresos esperados
- Rollover automático de remanentes al siguiente mes

**Non-Goals:**
- Importación automática bancaria (linked accounts) — se mantiene solo importación CSV manual
- Multi-moneda — Chele opera en una sola moneda (DOP/USD según presupuesto)
- Reglas de automatización complejas (if-this-then-that) más allá de las existentes
- Reportes nuevos fuera del alcance YNAB core (los reportes existentes se mantienen)

## Decisions

1. **Nuevo modelo Target en vez de extender Goal**: Goal existe pero está acoplado a MonthlyBudget. Target necesita su propia lógica (snooze, refill-up-to, underfunded calculation). Se crea modelo `Target` con campos: `category`, `goal_type`, `amount`, `target_date`, `frequency`, `refill_up_to` (bool), `snooze_month`, `snooze_year`. Migración: los Goals existentes se convierten a Targets con un script one-shot.

2. **Cálculo de underfunded en backend vía servicio**: El cálculo de cuánto falta para cumplir cada target se hace en Python (service layer), no en SQL ni en template. Esto permite lógica compleja (refill-up-to con rollover, targets anuales prorrateados). Se expone vía API y contexto de template.

3. **Ready to Assign como fórmula en memoria**: Ready to Assign = (saldo cuentas on_budget) - (suma de budgeted de todas las categorías del mes actual). No se persiste, se calcula en cada request. Zero-sum se valida al asignar: si el POST de assign_funds excede Ready to Assign, se rechaza.

4. **Overspending handling vía signals o service hook**: Al crear/editar/eliminar una transacción, se verifica si la categoría queda en overspent. Si es cash → marca roja y requiere cover. Si es TC → mueve automáticamente el overspent a la categoría de pago de esa TC. Se implementa como función llamada desde los views de transacciones (no signals, para mantener consistencia con el patrón actual).

5. **Inspector Panel como componente HTML/CSS con htmx o fetch**: El inspector se carga al hacer clic en una categoría. Se usa un endpoint `/presupuestos/categoria/<id>/inspector/` que devuelve HTML parcial (o JSON para API). En móvil se muestra como modal bottom-sheet.

6. **Spotlight como sección fija en base.html**: El spotlight es un área colapsable en la parte superior del main, justo debajo del header. Muestra alertas agrupadas: transacciones sin categorizar, sobregastos sin cubrir, targets no financiados. Cada alerta tiene acción (Review, Cover, Assign).

7. **Layout tripartito con CSS Grid**: En desktop `md:grid-cols-[16rem_1fr_20rem]`. Panel izquierdo (sidebar actual) se mantiene como está. Panel central (budget view) se expande. Panel derecho (inspector) es nuevo, oculto hasta que se selecciona una categoría.

8. **Barra inferior móvil con tres tabs**: Se reemplaza el header actual en móvil por una barra inferior fija con iconos + labels. Los tabs son: Plan (budget view), Accounts (account list), Spending (reports dashboard). La navegación es client-side via JavaScript (mostrar/ocultar secciones) o server-side via URLs.

9. **Drag-and-drop con SortableJS**: Librería externa mínima para reordenamiento de categorías/grupos. Se integra con endpoints POST para persistir el nuevo orden.

10. **Refill up to vs Set aside another**: Ambos son tipos de Target. `refill_up_to=True` significa que el target es "completar hasta el monto". El cálculo de underfunded = max(0, amount - rollover_balance). `refill_up_to=False` (Set aside another) = underfunded = amount (pide el monto completo cada mes).

## Risks / Trade-offs

- [Rendimiento] Ready to Assign se calcula en cada request con queries agregadas. Para presupuestos con muchas categorías y transacciones, puede lentificarse. Mitigación: cache de 30 segundos para el valor de Ready to Assign.
- [Complejidad] La migración de Goal a Target requiere script one-shot. Riesgo de pérdida de datos si no se prueba exhaustivamente. Mitigación: backup de DB antes de migrar, prueba en staging.
- [UX] El layout tripartito puede sentirse apretado en tablets. Mitigación: el panel inspector se oculta automáticamente en pantallas <1024px, accesible vía modal.
- [Deuda técnica] La lógica de sobregasto en TC (auto-move a payment) es frágil si hay múltiples TCs. Mitigación: cada transacción de TC verifica su categoría de pago asociada; si no existe, se crea automáticamente al registrar la TC.
- [Mobile] La barra inferior con 3 tabs requiere reestructurar la navegación móvil actual. Compatibilidad hacia atrás asegurada manteniendo el header actual como fallback.

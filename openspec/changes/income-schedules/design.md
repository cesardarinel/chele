## Context

El modelo `Schedule` actualmente solo soporta gastos: `process_due_schedules` siempre crea transacciones con `amount=-abs(amount)` y resta del saldo de la cuenta. No hay campo de dirección, por lo que no es posible programar ingresos recurrentes como sueldos o pagos de clientes.

La transacción (`Transaction`) ya soporta montos positivos (ingresos) y negativos (gastos). Solo hace falta agregar la dirección a `Schedule` para que al ejecutarse genere el signo correcto.

## Goals / Non-Goals

**Goals:**
- Agregar campo `direction` al modelo `Schedule` con opciones `income` y `expense`
- `process_due_schedules` debe crear transacciones con monto positivo y aumentar el saldo cuando `direction=income`
- Formulario debe tener toggle ingreso/gasto visible
- Lista debe mostrar color y signo según dirección
- Backward compatible: schedules existentes sin dirección se tratan como `expense`

**Non-Goals:**
- No cambiar el modelo `Transaction` (ya soporta ambos signos)
- No agregar transferencias entre cuentas propias (eso usa `transfer_id`)
- No cambiar la frecuencia o lógica de fechas

## Decisions

| Decisión | Opción | Razón |
|----------|--------|-------|
| Nombre del campo | `direction` con choices `income`/`expense` | Inglés consistente con el resto del modelo (`frequency`, `is_active`) |
| Default | `expense` | Backward compatible con registros existentes |
| Signo en ingreso | `Transaction.amount = abs(amount)` y `account.balance += amount` | Consistente con cómo se crean ingresos manualmente en `transaction_create` |
| Visual en lista | `+$` verde para ingreso, `-$` rojo para gasto | Misma convención que YNAB y el registro de cuenta |

## Risks / Trade-offs

- Schedules existentes sin `direction` se migran con default `expense` → comportamiento idéntico al anterior
- Si el usuario cambia una programación de gasto a ingreso existente, el próximo ciclo usará la nueva dirección sin afectar ejecuciones pasadas

## Migration Plan

1. `python manage.py makemigrations` genera migración con default `expense`
2. `python manage.py migrate` aplica
3. Schedules existentes heredan `expense` automáticamente

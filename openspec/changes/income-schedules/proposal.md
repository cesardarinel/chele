## Why

Las programaciones solo soportan gastos (egresos). En YNAB y Actual Budget se pueden programar tanto ingresos (sueldo, pagos de clientes) como gastos. Sin esto, el usuario no puede planificar sus ingresos recurrentes.

## What Changes

- Agregar campo `direction` al modelo `Schedule` con valores `income` (ingreso) y `expense` (gasto)
- Actualizar `process_due_schedules` para que los ingresos creen transacciones con monto positivo y aumenten el saldo de la cuenta
- Agregar toggle ingreso/gasto en el formulario de creación y edición
- Mostrar color verde con `+` para ingresos y rojo con `-` para gastos en la lista
- Crear migración de base de datos
- Agregar tests para programaciones de ingreso

## Capabilities

### New Capabilities
- `scheduled-income`: Capacidad de programar transacciones de ingreso recurrentes (sueldo, pagos de clientes, etc.) con dirección ingreso/gasto, que al ejecutarse crean transacciones con el signo correcto y actualizan el saldo de la cuenta apropiadamente.

### Modified Capabilities
<!-- No existing specs to modify -->

## Impact

- `apps/schedules/models.py` — nuevo campo `direction` en Schedule
- `apps/schedules/views.py` — lógica de `process_due_schedules` usa la dirección para determinar signo del monto; formularios guardan dirección
- `templates/schedules/schedule_form.html` — toggle ingreso/gasto en el formulario
- `templates/schedules/schedule_list.html` — color y signo según dirección
- `apps/schedules/tests.py` — test de creación y ejecución de ingreso programado
- Migración de base de datos

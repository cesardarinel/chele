# Guía de Uso — Chele

Chele es una aplicación de finanzas personales basada en la metodología YNAB (You Need A Budget). Te permite crear múltiples presupuestos independientes (familia, empresa, personal, etc.) y gestionar tus finanzas de forma colaborativa.

---

## Índice

1. [Conceptos fundamentales](#1-conceptos-fundamentales)
2. [Primeros pasos](#2-primeros-pasos)
3. [Las 4 reglas YNAB en Chele](#3-las-4-reglas-ynab-en-chele)
4. [Uso diario](#4-uso-diario)
5. [Múltiples presupuestos](#5-múltiples-presupuestos)
6. [Trabajo en equipo](#6-trabajo-en-equipo)
7. [Preguntas frecuentes](#7-preguntas-frecuentes)

---

## 1. Conceptos fundamentales

```
Chele organiza tus finanzas así:

Usuario
  └── Presupuesto (familia, empresa, personal, etc.)
        ├── Cuentas  (bancarias, efectivo, tarjetas)
        ├── Categorías (comida, servicios, ahorro, etc.)
        │     └── Metas (objetivos de ahorro)
        ├── Transacciones (ingresos y gastos)
        ├── Beneficiarios (a quién le pagas)
        └── Programaciones (gastos recurrentes)
```

### Presupuesto vs Cuenta
| Concepto | Qué es |
|---|---|
| **Presupuesto** | Un plan financiero completo. Ej: "Familia", "Negocio", "Ahorro viaje". Cada presupuesto es independiente. |
| **Cuenta** | Un lugar donde guardás dinero. Ej: "Banco Nación", "Billetera", "Visa". Las cuentas viven dentro de un presupuesto. |

### Categorías
Son los "sobres" donde asignás tu dinero. Se agrupan en grupos:

| Grupo | Categorías ejemplo |
|---|---|
| Gastos Fijos | Alquiler, Internet, Seguro |
| Gastos Diarios | Comida, Transporte, Salidas |
| Ahorro | Fondo de emergencia, Vacaciones |
| True Expenses | Seguro auto (anual), Impuesto municipal |

---

## 2. Primeros pasos

### 2.1 Crear tu primer presupuesto

```
1. Registrate con email y contraseña
2. Elegí "Crear nuevo presupuesto"
3. Ponele nombre: "Familia", "Mi presupuesto", etc.
4. ¡Listo! Ya tenés tu primer presupuesto vacío
```

### 2.2 Agregar cuentas

Agregá todas las cuentas que tengas:

```
Ejemplo:
  ┌─────────────────────────────────────────┐
  │  Nueva Cuenta                           │
  │                                         │
  │  Nombre:    Banco Nación                │
  │  Tipo:      ● Corriente  ○ Ahorro       │
  │             ○ Efectivo  ○ Tarjeta       │
  │  Saldo:     $ 150,000                   │
  │  Incluir:   ● En el presupuesto         │
  │             ○ Solo seguimiento          │
  │                                         │
  │  [Crear Cuenta]                         │
  └─────────────────────────────────────────┘
```

**¿En el presupuesto o solo seguimiento?**
- **En el presupuesto**: el dinero de esta cuenta se puede asignar a categorías. Usá esto para cuentas del día a día.
- **Solo seguimiento**: la cuenta aparece en tu patrimonio neto pero no afecta el presupuesto. Usá esto para inversiones, hipotecas, etc.

### 2.3 Configurar categorías

Chele crea categorías por defecto, pero podés editarlas:

```
Grupos de ejemplo:
  ├── Gastos Fijos
  │     ├── Alquiler / Hipoteca
  │     ├── Servicios (luz, agua, internet)
  │     └── Suscripciones
  ├── Gastos Diarios
  │     ├── Comida
  │     ├── Transporte
  │     └── Salidas / Ocio
  ├── Ahorro
  │     ├── Fondo de emergencia
  │     └── Vacaciones
  └── True Expenses
        └── Seguro auto (ver metas)
```

### 2.4 Asignar tu primer presupuesto

```
1. Andá a la vista "Presupuesto"
2. Fijate en "A presupuestar" (arriba) — es tu dinero disponible
3. Asigná montos a cada categoría hasta que "A presupuestar" llegue a $0
4. ¡Eso es todo! Cada peso tiene un trabajo
```

### 2.5 Importar transacciones

Si venís usando otro sistema o tenés extractos bancarios:

```
1. Andá a la cuenta que querés importar
2. Click en "Importar"
3. Seleccioná tu archivo CSV
4. Mapeá las columnas (fecha, monto, beneficiario)
5. Revisá el preview
6. Confirmá la importación
```

También podés cargar transacciones manualmente con "Añadir nueva".

---

## 3. Las 4 reglas YNAB en Chele

### Regla 1: Darle trabajo a cada peso

Cada peso que ingresa debe asignarse a una categoría. Cuando "A presupuestar" llega a $0, todo tu dinero tiene un propósito.

```
Ingresaste $100,000 este mes:
  ├── Alquiler:        $30,000
  ├── Comida:          $20,000
  ├── Transporte:      $10,000
  ├── Ahorro:          $15,000
  ├── Salidas:         $10,000
  └── Seguro auto:     $15,000  ← True Expense
                        ───────
  A presupuestar:        $0  ✅
```

### Regla 2: Aceptar tus gastos reales

Usá **metas** para planificar gastos que sabés que van a llegar:

| Meta | Ejemplo | Comportamiento |
|---|---|---|
| Monthly Savings Goal | Ahorrar $10,000/mes | Asigna $10,000 automáticamente cada mes |
| Target Balance | Fondo de emergencia $500,000 | Asigna hasta llegar a $500,000 |
| Target by Date | Viaje $300,000 para Dic 2026 | Calcula cuota mensual hasta la fecha |
| True Expense | Seguro auto $180,000/año | Asigna $15,000/mes automáticamente |

**¿Cómo configurar una meta?**
```
1. En la vista Presupuesto, hover sobre una categoría
2. Click en "Configurar meta"
3. Elegí el tipo de meta y completá los datos
4. Chele calculará y asignará automáticamente
```

### Regla 3: Patear los golpes

Los imprevistos pasan. Cuando te pasás en una categoría:

```
Situación: Gastaste $25,000 en Comida pero solo presupuestaste $20,000

Opción 1: Mové dinero de otra categoría
  → Click en Saldo de "Salidas", transferí $5,000 a "Comida"

Opción 2: Dejá que se descuente el próximo mes
  → "A presupuestar" del próximo mes tendrá $5,000 menos
```

### Regla 4: Envejecer tu dinero

El objetivo es vivir con el dinero del **mes pasado**. Para eso usá "Reservar para el próximo mes":

```
1. Cuando te llegue un ingreso, no lo asignes todavía
2. Click en "A presupuestar"
3. Elegí "Reservar para el próximo mes"
4. El dinero estará disponible el mes que viene
5. Ideal: que todo el presupuesto del mes se cubra con dinero del mes anterior

Podés activar "Reserva automática" en categorías de ingreso
para que siempre se retenga automáticamente.
```

---

## 4. Uso diario

### Cada vez que hacés un gasto

```
1. Andá a la Cuenta correspondiente
2. Click en "Añadir nuevo"
3. Completá:
     ● Fecha:     hoy
     ● Beneficiario: "Supermercado ABC"
     ● Categoría: "Comida"
     ● Nota:      (opcional)
     ● Monto:     $12,500
4. Guardá
5. Chele actualiza el saldo de la cuenta y el gasto de la categoría automáticamente
```

### Cada vez que recibís un ingreso

```
1. Andá a la Cuenta donde recibiste el dinero
2. Añadí una nueva transacción con monto positivo
3. Categorizalo como "Ingreso" (o la categoría específica)
4. El monto aparece en "A presupuestar"
5. Asignalo a las categorías que necesites
```

### Fin de mes

```
1. Revisá si alguna categoría quedó negativa → mové dinero para cubrirla
2. Lo que sobra en las categorías se acumula para el próximo mes
3. Si tenés ingreso nuevo, decidí: ¿lo asigno ahora o lo reservo para el mes que viene?
4. Creá el presupuesto del próximo mes (podés copiar el anterior)
```

### Sync entre dispositivos

```
1. En cualquier dispositivo, click en el botón de sync (↻)
2. Chele envía tus cambios al servidor y descarga los cambios de otros
3. Si dos personas editaron lo mismo, gana el último que guardó
4. Podés trabajar sin conexión: los cambios quedan pendientes hasta el próximo sync
```

---

## 5. Múltiples presupuestos

Podés crear tantos presupuestos como quieras. Cada uno es completamente independiente.

```
Ejemplo de uso real:
  ┌─────────────────────────────────────────────┐
  │  Presupuesto activo: [Familia ▼]            │
  │                                             │
  │  ┌─────────────────────────────────────────┐│
  │  │  + Nuevo Presupuesto                     ││
  │  ├─────────────────────────────────────────┤│
  │  │  ● Familia           ← activo           ││
  │  │  ○ Emprendimiento                       ││
  │  │  ○ Ahorro Viaje Europa                  ││
  │  │  ○ Proyecto Personal                    ││
  │  └─────────────────────────────────────────┘│
  └─────────────────────────────────────────────┘
```

**Características:**
- Cada presupuesto tiene sus propias cuentas, categorías y transacciones
- Podés compartir cada presupuesto con diferentes personas
- El switch es instantáneo desde la barra lateral
- Los reportes, programaciones y reglas son por presupuesto

---

## 6. Trabajo en equipo

Podés compartir un presupuesto con tu familia, socios, etc.

```
1. Andá a Configuración del presupuesto
2. "Invitar miembro"
3. Ingresá el email de la persona
4. La persona recibe una invitación
5. Al aceptarla, puede ver y editar el presupuesto

Todos los miembros pueden:
  ✓ Ver y editar transacciones
  ✓ Asignar y mover dinero entre categorías
  ✓ Importar y exportar datos

Nota: Solo el dueño puede eliminar el presupuesto o invitar/remover miembros.
```

---

## 7. Preguntas frecuentes

### ¿Cuál es la diferencia entre "En el presupuesto" y "Solo seguimiento"?
Las cuentas **en el presupuesto** aportan dinero a "A presupuestar". Las de **solo seguimiento** solo se ven en el patrimonio neto. Ej: tu cuenta sueldo va "en el presupuesto", tu fondo de inversión va "solo seguimiento".

### ¿Qué pasa si gasto más de lo presupuestado?
No pasa nada grave. Chele te muestra la categoría en rojo. Podés mover dinero de otra categoría para cubrirlo, o dejarlo y se descuenta del "A presupuestar" del próximo mes.

### ¿Qué son las True Expenses?
Son gastos que pagás una vez al año (seguro, patente, impuestos) pero que conviene ahorrar mes a mes. Configurás una meta True Expense y Chele divide el total en 12 cuotas mensuales automáticas.

### ¿Puedo tener un presupuesto para mi negocio y otro personal?
Sí. Creá un presupuesto "Negocio" y otro "Personal". Son completamente independientes. Cada uno tiene sus cuentas, categorías, y podés compartir solo el del negocio con tu socio y solo el personal con tu familia.

### ¿Cómo sé si estoy mejorando?
Mirá el reporte de **Patrimonio Neto** (tus activos menos tus deudas) a lo largo del tiempo. Si la línea va para arriba, vas bien. También el reporte de **Budget vs Reality** te muestra si estás cumpliendo tus presupuestos.

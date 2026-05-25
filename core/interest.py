from datetime import date, timedelta


def calcular_interes_diario(saldo, tasa_anual, dias_atraso):
    if saldo <= 0 or tasa_anual <= 0 or dias_atraso <= 0:
        return 0.0
    tasa_diaria = float(tasa_anual) / 12.0 / 30.0
    return round(float(saldo) * tasa_diaria * int(dias_atraso), 2)


def generar_intereses(saldo, tasa_anual, desde, hasta=None):
    if hasta is None:
        hasta = date.today()
    dias = (hasta - desde).days
    return calcular_interes_diario(saldo, tasa_anual, dias)


def _get_saldo(objeto):
    if hasattr(objeto, 'remaining_balance'):
        return float(objeto.remaining_balance)
    if hasattr(objeto, 'balance'):
        return float(objeto.balance)
    return 0


def _set_saldo(objeto, delta):
    if hasattr(objeto, 'remaining_balance'):
        objeto.remaining_balance = float(objeto.remaining_balance) + delta
    elif hasattr(objeto, 'balance'):
        objeto.balance = float(objeto.balance) - delta
    objeto.save()


def aplicar_interes(objeto_con_saldo_y_tasa, desde):
    from apps.transactions.models import Transaction
    saldo = abs(_get_saldo(objeto_con_saldo_y_tasa))
    interes = generar_intereses(
        saldo,
        objeto_con_saldo_y_tasa.interest_rate,
        desde,
    )
    if interes > 0:
        Transaction.objects.create(
            budget=objeto_con_saldo_y_tasa.budget,
            account=None,
            date=date.today(),
            amount=-interes,
            notes=f'Interés {objeto_con_saldo_y_tasa}',
        )
        _set_saldo(objeto_con_saldo_y_tasa, interes)
    return interes

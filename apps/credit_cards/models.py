import uuid
from datetime import date
from django.db import models


class CreditCard(models.Model):
    id = models.UUIDField(primary_key=True, default=uuid.uuid4, editable=False)
    budget = models.ForeignKey('budgets.Budget', on_delete=models.CASCADE, related_name='credit_cards', verbose_name='presupuesto')
    name = models.CharField(max_length=200, verbose_name='nombre')
    limit = models.DecimalField(max_digits=15, decimal_places=2, verbose_name='límite')
    balance = models.DecimalField(max_digits=15, decimal_places=2, default=0, verbose_name='saldo actual', help_text='Negativo = deuda')
    interest_rate = models.DecimalField(max_digits=6, decimal_places=4, verbose_name='tasa anual', help_text='Ej: 0.96 = 96% anual (8% mensual)')
    closing_day = models.IntegerField(verbose_name='día de cierre')
    due_day = models.IntegerField(verbose_name='día de vencimiento')
    last_interest_date = models.DateField(null=True, blank=True, verbose_name='último interés calculado')
    notes = models.TextField(blank=True, verbose_name='notas')
    created_at = models.DateTimeField(auto_now_add=True)
    updated_at = models.DateTimeField(auto_now=True)

    class Meta:
        verbose_name = 'tarjeta de crédito'
        verbose_name_plural = 'tarjetas de crédito'

    def __str__(self):
        return self.name

    def calcular_interes(self, desde=None):
        from core.interest import generar_intereses
        if desde is None:
            desde = self.last_interest_date or date.today()
        hoy = date.today()
        if self.balance >= 0 or hoy <= desde:
            return 0
        return generar_intereses(abs(float(self.balance)), self.interest_rate, desde, hoy)


class InterestCharge(models.Model):
    id = models.UUIDField(primary_key=True, default=uuid.uuid4, editable=False)
    credit_card = models.ForeignKey(CreditCard, on_delete=models.CASCADE, related_name='interest_charges', verbose_name='tarjeta')
    amount = models.DecimalField(max_digits=15, decimal_places=2, verbose_name='monto')
    date = models.DateField(auto_now_add=True, verbose_name='fecha')
    days_overdue = models.IntegerField(verbose_name='días de atraso')
    applied = models.BooleanField(default=False, verbose_name='aplicado')

    class Meta:
        verbose_name = 'cargo de interés'
        verbose_name_plural = 'cargos de intereses'

    def __str__(self):
        return f'Interés {self.amount} - {self.credit_card.name}'

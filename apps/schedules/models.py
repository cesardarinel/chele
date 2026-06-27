import uuid
from django.db import models


class Schedule(models.Model):
    FREQUENCY_CHOICES = [
        ('weekly', 'Semanal'),
        ('biweekly', 'Quincenal'),
        ('monthly', 'Mensual'),
        ('quarterly', 'Trimestral'),
        ('yearly', 'Anual'),
    ]
    DIRECTION_CHOICES = [
        ('expense', 'Gasto'),
        ('income', 'Ingreso'),
    ]
    id = models.UUIDField(primary_key=True, default=uuid.uuid4, editable=False)
    budget = models.ForeignKey('budgets.Budget', on_delete=models.CASCADE, related_name='schedules', verbose_name='presupuesto')
    payee = models.ForeignKey('payees.Payee', on_delete=models.SET_NULL, null=True, blank=True, verbose_name='beneficiario')
    category = models.ForeignKey('budgets.Category', on_delete=models.SET_NULL, null=True, blank=True, verbose_name='categoría')
    account = models.ForeignKey('accounts.Account', on_delete=models.CASCADE, verbose_name='cuenta')
    amount = models.DecimalField(max_digits=15, decimal_places=2, verbose_name='monto')
    frequency = models.CharField(max_length=20, choices=FREQUENCY_CHOICES, verbose_name='frecuencia')
    next_date = models.DateField(verbose_name='próxima fecha')
    notes = models.TextField(blank=True, verbose_name='notas')
    skip_weekends = models.BooleanField(default=False, verbose_name='saltar fines de semana (pasar al lunes)')
    apply_before_weekend = models.BooleanField(default=False, verbose_name='aplicar antes del finde (pasar al viernes)')
    is_active = models.BooleanField(default=True, verbose_name='activo')
    direction = models.CharField(max_length=10, choices=DIRECTION_CHOICES, default='expense', verbose_name='dirección')
    created_at = models.DateTimeField(auto_now_add=True)
    updated_at = models.DateTimeField(auto_now=True)

    class Meta:
        verbose_name = 'programación'
        verbose_name_plural = 'programaciones'

    def __str__(self):
        return f'{self.payee.name if self.payee else "—"} - {self.amount}'

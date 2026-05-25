import uuid
from django.db import models


class Transaction(models.Model):
    id = models.UUIDField(primary_key=True, default=uuid.uuid4, editable=False)
    budget = models.ForeignKey('budgets.Budget', on_delete=models.CASCADE, related_name='transactions', verbose_name='presupuesto')
    account = models.ForeignKey('accounts.Account', on_delete=models.CASCADE, null=True, blank=True, related_name='transactions', verbose_name='cuenta')
    date = models.DateField(verbose_name='fecha')
    amount = models.DecimalField(max_digits=15, decimal_places=2, verbose_name='monto')
    payee = models.ForeignKey('payees.Payee', on_delete=models.SET_NULL, null=True, blank=True, related_name='transactions', verbose_name='beneficiario')
    category = models.ForeignKey('budgets.Category', on_delete=models.SET_NULL, null=True, blank=True, related_name='transactions', verbose_name='categoría')
    notes = models.TextField(blank=True, verbose_name='notas')
    transfer_id = models.UUIDField(null=True, blank=True, verbose_name='id de transferencia')
    reconciled = models.BooleanField(default=False, verbose_name='conciliado')
    cleared = models.BooleanField(default=True, verbose_name='confirmado')
    created_at = models.DateTimeField(auto_now_add=True)
    updated_at = models.DateTimeField(auto_now=True)

    class Meta:
        verbose_name = 'transacción'
        verbose_name_plural = 'transacciones'
        ordering = ['-date', '-created_at']

    def __str__(self):
        return f'{self.date} - {self.amount}'

import uuid
from django.db import models


class Account(models.Model):
    id = models.UUIDField(primary_key=True, default=uuid.uuid4, editable=False)
    budget = models.ForeignKey('budgets.Budget', on_delete=models.CASCADE, related_name='accounts', verbose_name='presupuesto')
    name = models.CharField(max_length=200, verbose_name='nombre')
    on_budget = models.BooleanField(default=True, verbose_name='en el presupuesto')
    balance = models.DecimalField(max_digits=15, decimal_places=2, default=0, verbose_name='saldo')
    notes = models.TextField(blank=True, verbose_name='notas')
    created_at = models.DateTimeField(auto_now_add=True)
    updated_at = models.DateTimeField(auto_now=True)

    class Meta:
        verbose_name = 'cuenta'
        verbose_name_plural = 'cuentas'
        ordering = ['name']

    def __str__(self):
        return self.name

import uuid
from django.db import models


class Rule(models.Model):
    CONDITION_FIELDS = [
        ('payee', 'Beneficiario'),
        ('amount', 'Monto'),
        ('notes', 'Notas'),
    ]
    CONDITION_OPERATORS = [
        ('contains', 'contiene'),
        ('equals', 'es igual a'),
        ('starts_with', 'empieza con'),
        ('greater_than', 'mayor que'),
        ('less_than', 'menor que'),
    ]
    id = models.UUIDField(primary_key=True, default=uuid.uuid4, editable=False)
    budget = models.ForeignKey('budgets.Budget', on_delete=models.CASCADE, related_name='rules', verbose_name='presupuesto')
    name = models.CharField(max_length=200, verbose_name='nombre')
    condition_field = models.CharField(max_length=20, choices=CONDITION_FIELDS, verbose_name='campo')
    condition_operator = models.CharField(max_length=20, choices=CONDITION_OPERATORS, verbose_name='operador')
    condition_value = models.CharField(max_length=500, verbose_name='valor')
    action_category = models.ForeignKey('budgets.Category', on_delete=models.SET_NULL, null=True, blank=True, verbose_name='asignar categoría')
    action_payee = models.ForeignKey('payees.Payee', on_delete=models.SET_NULL, null=True, blank=True, verbose_name='asignar beneficiario')
    action_notes = models.CharField(max_length=500, blank=True, verbose_name='asignar notas')
    sort_order = models.IntegerField(default=0, verbose_name='orden')
    is_active = models.BooleanField(default=True, verbose_name='activa')
    created_at = models.DateTimeField(auto_now_add=True)
    updated_at = models.DateTimeField(auto_now=True)

    class Meta:
        verbose_name = 'regla'
        verbose_name_plural = 'reglas'
        ordering = ['sort_order']

    def __str__(self):
        return self.name

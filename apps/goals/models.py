import uuid
from django.db import models
from django.core.validators import MinValueValidator


class Goal(models.Model):
    GOAL_TYPES = [
        ('monthly', 'Ahorro mensual'),
        ('target_balance', 'Saldo objetivo'),
        ('target_date', 'Objetivo con fecha'),
        ('true_expense', 'Gasto anual'),
    ]
    id = models.UUIDField(primary_key=True, default=uuid.uuid4, editable=False)
    category = models.ForeignKey('budgets.Category', on_delete=models.CASCADE, related_name='goals', verbose_name='categoría')
    goal_type = models.CharField(max_length=20, choices=GOAL_TYPES, verbose_name='tipo de meta')
    amount = models.DecimalField(max_digits=15, decimal_places=2, validators=[MinValueValidator(0)], verbose_name='monto')
    target_date = models.DateField(null=True, blank=True, verbose_name='fecha objetivo')
    frequency = models.IntegerField(default=12, verbose_name='frecuencia (meses)')
    is_completed = models.BooleanField(default=False, verbose_name='completada')
    created_at = models.DateTimeField(auto_now_add=True)
    updated_at = models.DateTimeField(auto_now=True)

    class Meta:
        verbose_name = 'meta'
        verbose_name_plural = 'metas'

    def __str__(self):
        return f'{self.get_goal_type_display()} - {self.category.name}'

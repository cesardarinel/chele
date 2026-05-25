import uuid
from django.db import models


class Payee(models.Model):
    id = models.UUIDField(primary_key=True, default=uuid.uuid4, editable=False)
    budget = models.ForeignKey('budgets.Budget', on_delete=models.CASCADE, related_name='payees', verbose_name='presupuesto')
    name = models.CharField(max_length=200, verbose_name='nombre')
    merge_to = models.ForeignKey('self', on_delete=models.SET_NULL, null=True, blank=True, verbose_name='fusionar con')
    created_at = models.DateTimeField(auto_now_add=True)
    updated_at = models.DateTimeField(auto_now=True)

    class Meta:
        verbose_name = 'beneficiario'
        verbose_name_plural = 'beneficiarios'
        ordering = ['name']

    def __str__(self):
        return self.name

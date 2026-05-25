import uuid
from django.db import models


class SyncLog(models.Model):
    id = models.UUIDField(primary_key=True, default=uuid.uuid4, editable=False)
    budget = models.ForeignKey('budgets.Budget', on_delete=models.CASCADE, related_name='sync_logs', verbose_name='presupuesto')
    user = models.ForeignKey('auth.User', on_delete=models.CASCADE, verbose_name='usuario')
    entity_type = models.CharField(max_length=50, verbose_name='tipo de entidad')
    entity_id = models.UUIDField(verbose_name='id de entidad')
    action = models.CharField(max_length=20, verbose_name='acción')
    payload = models.JSONField(default=dict, verbose_name='datos')
    created_at = models.DateTimeField(auto_now_add=True)

    class Meta:
        verbose_name = 'registro de sync'
        verbose_name_plural = 'registros de sync'

    def __str__(self):
        return f'{self.entity_type} - {self.action}'

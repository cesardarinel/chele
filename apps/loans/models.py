import uuid
from datetime import date
from django.db import models


class Loan(models.Model):
    LOAN_TYPES = [
        ('personal', 'Personal'),
        ('hipotecario', 'Hipotecario'),
        ('automotor', 'Automotor'),
        ('estudiantil', 'Estudiantil'),
        ('otros', 'Otros'),
    ]
    STATUS_CHOICES = [
        ('active', 'Activo'),
        ('completed', 'Finalizado'),
    ]
    id = models.UUIDField(primary_key=True, default=uuid.uuid4, editable=False)
    budget = models.ForeignKey('budgets.Budget', on_delete=models.CASCADE, related_name='loans', verbose_name='presupuesto')
    type = models.CharField(max_length=20, choices=LOAN_TYPES, verbose_name='tipo')
    name = models.CharField(max_length=200, verbose_name='nombre')
    status = models.CharField(max_length=20, choices=STATUS_CHOICES, default='active', verbose_name='estado')
    total_amount = models.DecimalField(max_digits=15, decimal_places=2, verbose_name='monto total')
    interest_rate = models.DecimalField(max_digits=6, decimal_places=4, verbose_name='tasa anual', help_text='Ej: 0.36 = 36% anual')
    remaining_balance = models.DecimalField(max_digits=15, decimal_places=2, verbose_name='saldo pendiente')
    total_installments = models.IntegerField(verbose_name='total cuotas')
    paid_installments = models.IntegerField(default=0, verbose_name='cuotas pagadas')
    start_date = models.DateField(verbose_name='fecha inicio')
    next_due_date = models.DateField(verbose_name='próximo vencimiento')
    installment_amount = models.DecimalField(max_digits=15, decimal_places=2, verbose_name='valor cuota')
    notes = models.TextField(blank=True, verbose_name='notas')
    created_at = models.DateTimeField(auto_now_add=True)
    updated_at = models.DateTimeField(auto_now=True)

    class Meta:
        verbose_name = 'préstamo'
        verbose_name_plural = 'préstamos'

    def __str__(self):
        return f'{self.get_type_display()} - {self.name}'

    def calcular_interes(self, desde=None):
        from core.interest import generar_intereses
        if desde is None:
            desde = self.next_due_date
        hoy = date.today()
        if hoy <= desde:
            return 0
        return generar_intereses(self.remaining_balance, self.interest_rate, desde, hoy)


class Installment(models.Model):
    id = models.UUIDField(primary_key=True, default=uuid.uuid4, editable=False)
    loan = models.ForeignKey(Loan, on_delete=models.CASCADE, related_name='installments', verbose_name='préstamo')
    number = models.IntegerField(verbose_name='número')
    amount = models.DecimalField(max_digits=15, decimal_places=2, verbose_name='monto')
    due_date = models.DateField(verbose_name='fecha vencimiento')
    paid = models.BooleanField(default=False, verbose_name='pagada')
    paid_date = models.DateField(null=True, blank=True, verbose_name='fecha pago')
    notes = models.TextField(blank=True, verbose_name='notas')

    class Meta:
        verbose_name = 'cuota'
        verbose_name_plural = 'cuotas'
        ordering = ['number']

    def __str__(self):
        return f'Cuota {self.number} - {self.loan.name}'

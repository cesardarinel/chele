import uuid
from django.db import models
from django.conf import settings


class Budget(models.Model):
    id = models.UUIDField(primary_key=True, default=uuid.uuid4, editable=False)
    name = models.CharField(max_length=200, verbose_name='nombre')
    description = models.TextField(blank=True, verbose_name='descripción')
    owner = models.ForeignKey(settings.AUTH_USER_MODEL, on_delete=models.CASCADE, related_name='owned_budgets', verbose_name='dueño')
    members = models.ManyToManyField(settings.AUTH_USER_MODEL, through='BudgetMembership', related_name='budgets', verbose_name='miembros')
    created_at = models.DateTimeField(auto_now_add=True)
    updated_at = models.DateTimeField(auto_now=True)

    class Meta:
        verbose_name = 'presupuesto'
        verbose_name_plural = 'presupuestos'

    def __str__(self):
        return self.name


class BudgetMembership(models.Model):
    ROLE_CHOICES = [
        ('owner', 'Dueño'),
        ('editor', 'Editor'),
    ]
    id = models.UUIDField(primary_key=True, default=uuid.uuid4, editable=False)
    user = models.ForeignKey(settings.AUTH_USER_MODEL, on_delete=models.CASCADE)
    budget = models.ForeignKey(Budget, on_delete=models.CASCADE)
    role = models.CharField(max_length=20, choices=ROLE_CHOICES, default='editor', verbose_name='rol')
    invited_at = models.DateTimeField(auto_now_add=True)
    accepted_at = models.DateTimeField(null=True, blank=True)

    class Meta:
        unique_together = ('user', 'budget')
        verbose_name = 'membresía'
        verbose_name_plural = 'membresías'

    def __str__(self):
        return f'{self.user.email} en {self.budget.name}'


class CategoryGroup(models.Model):
    id = models.UUIDField(primary_key=True, default=uuid.uuid4, editable=False)
    budget = models.ForeignKey(Budget, on_delete=models.CASCADE, related_name='category_groups', verbose_name='presupuesto')
    name = models.CharField(max_length=200, verbose_name='nombre')
    sort_order = models.IntegerField(default=0, verbose_name='orden')
    is_income = models.BooleanField(default=False, verbose_name='es ingreso')
    created_at = models.DateTimeField(auto_now_add=True)
    updated_at = models.DateTimeField(auto_now=True)

    class Meta:
        verbose_name = 'grupo de categorías'
        verbose_name_plural = 'grupos de categorías'
        ordering = ['sort_order', 'name']

    def __str__(self):
        return self.name


class Category(models.Model):
    id = models.UUIDField(primary_key=True, default=uuid.uuid4, editable=False)
    budget = models.ForeignKey(Budget, on_delete=models.CASCADE, related_name='categories', verbose_name='presupuesto')
    group = models.ForeignKey(CategoryGroup, on_delete=models.CASCADE, related_name='categories', verbose_name='grupo')
    name = models.CharField(max_length=200, verbose_name='nombre')
    sort_order = models.IntegerField(default=0, verbose_name='orden')
    is_hidden = models.BooleanField(default=False, verbose_name='oculto')
    notes = models.TextField(blank=True, verbose_name='notas')
    created_at = models.DateTimeField(auto_now_add=True)
    updated_at = models.DateTimeField(auto_now=True)

    class Meta:
        verbose_name = 'categoría'
        verbose_name_plural = 'categorías'
        ordering = ['sort_order', 'name']

    def __str__(self):
        return self.name


class MonthlyBudget(models.Model):
    id = models.UUIDField(primary_key=True, default=uuid.uuid4, editable=False)
    category = models.ForeignKey(Category, on_delete=models.CASCADE, related_name='monthly_budgets', verbose_name='categoría')
    month = models.IntegerField(verbose_name='mes')
    year = models.IntegerField(verbose_name='año')
    budgeted = models.DecimalField(max_digits=15, decimal_places=2, default=0, verbose_name='presupuestado')
    created_at = models.DateTimeField(auto_now_add=True)
    updated_at = models.DateTimeField(auto_now=True)

    class Meta:
        unique_together = ('category', 'month', 'year')
        verbose_name = 'presupuesto mensual'
        verbose_name_plural = 'presupuestos mensuales'

    def __str__(self):
        return f'{self.category.name} - {self.month}/{self.year}'

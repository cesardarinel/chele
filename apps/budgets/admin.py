from django.contrib import admin
from .models import Budget, BudgetMembership, CategoryGroup, Category, MonthlyBudget

admin.site.register(Budget)
admin.site.register(BudgetMembership)
admin.site.register(CategoryGroup)
admin.site.register(Category)
admin.site.register(MonthlyBudget)

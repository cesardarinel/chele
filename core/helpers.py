from django.shortcuts import get_object_or_404
from apps.budgets.models import Budget


def get_active_budget(request):
    budget_id = request.session.get('active_budget_id')
    if not budget_id:
        first = Budget.objects.filter(members=request.user).first()
        if first:
            budget_id = str(first.id)
            request.session['active_budget_id'] = budget_id
    return get_object_or_404(Budget, id=budget_id, members=request.user)


def get_budget_id(request):
    budget_id = request.session.get('active_budget_id')
    if not budget_id:
        first = Budget.objects.filter(members=request.user).first()
        if first:
            budget_id = str(first.id)
            request.session['active_budget_id'] = budget_id
    return budget_id


def current_month_year(request):
    from datetime import date
    today = date.today()
    month = int(request.GET.get('mes', today.month))
    year = int(request.GET.get('anio', today.year))
    return month, year

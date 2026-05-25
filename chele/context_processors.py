from apps.budgets.models import Budget


def active_budget(request):
    if not request.user.is_authenticated:
        return {}
    budget_id = request.session.get('active_budget_id')
    budgets = Budget.objects.filter(members=request.user)
    active = budgets.filter(id=budget_id).first() if budget_id else budgets.first()
    if not budget_id and active:
        request.session['active_budget_id'] = str(active.id)
    return {
        'user_budgets': budgets,
        'active_budget': active,
    }

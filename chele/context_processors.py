from datetime import date
from django.db.models import Sum
from apps.budgets.models import Budget, Category, MonthlyBudget
from apps.transactions.models import Transaction
from apps.goals.services import TargetService
from apps.onboarding.models import OnboardingProfile


def active_budget(request):
    if not request.user.is_authenticated:
        return {}
    budget_id = request.session.get('active_budget_id')
    budgets = Budget.objects.filter(members=request.user)
    active = budgets.filter(id=budget_id).first() if budget_id else budgets.first()
    if not budget_id and active:
        request.session['active_budget_id'] = str(active.id)

    ctx = {'user_budgets': budgets, 'active_budget': active}

    if active:
        accounts = active.accounts.all()
        on_budget_qs = accounts.filter(on_budget=True)
        off_budget_qs = accounts.filter(on_budget=False)
        ctx['sidebar_on_budget'] = on_budget_qs
        ctx['sidebar_off_budget'] = off_budget_qs
        ctx['sidebar_total_on_budget'] = onBudget = on_budget_qs.aggregate(Sum('balance'))['balance__sum'] or 0
        ctx['sidebar_total_off_budget'] = offBudget = off_budget_qs.aggregate(Sum('balance'))['balance__sum'] or 0
        ctx['sidebar_grand_total'] = onBudget + offBudget

        # Spotlight alerts (global)
        today = date.today()
        ctx['uncategorized_count'] = Transaction.objects.filter(
            budget=active, category__isnull=True
        ).count()

        # Overspent categories (cash)
        overspends = []
        for cat in Category.objects.filter(budget=active, is_hidden=False):
            budgeted = MonthlyBudget.objects.filter(category=cat, month=today.month, year=today.year).aggregate(Sum('budgeted'))['budgeted__sum'] or 0
            spent = abs(float(Transaction.objects.filter(category=cat, date__month=today.month, date__year=today.year).aggregate(Sum('amount'))['amount__sum'] or 0))
            avail = float(budgeted) - spent
            if avail < -0.01:
                overspends.append({'id': str(cat.id), 'name': cat.name, 'amount': abs(avail)})
        ctx['uncovered_overspends'] = overspends

        ts = TargetService(active, today.month, today.year)
        underfunded = ts.list_underfunded()
        ctx['underfunded_count'] = len(underfunded)

    # Onboarding state
    onboarding_step = 7
    onboarding_active = False
    if request.user.is_authenticated:
        try:
            profile = request.user.onboarding
            onboarding_step = profile.step
            onboarding_active = profile.step < 7
        except OnboardingProfile.DoesNotExist:
            pass
    ctx['onboarding_step'] = onboarding_step
    ctx['onboarding_active'] = onboarding_active

    return ctx

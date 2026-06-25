from django.db.models import Sum
from apps.budgets.models import Budget
from apps.credit_cards.models import CreditCard
from apps.loans.models import Loan


def active_budget(request):
    if not request.user.is_authenticated:
        return {}
    budget_id = request.session.get('active_budget_id')
    budgets = Budget.objects.filter(members=request.user)
    active = budgets.filter(id=budget_id).first() if budget_id else budgets.first()
    if not budget_id and active:
        request.session['active_budget_id'] = str(active.id)

    ctx = {
        'user_budgets': budgets,
        'active_budget': active,
    }

    if active:
        accounts = active.accounts.all()
        on_budget_qs = accounts.filter(on_budget=True)
        off_budget_qs = accounts.filter(on_budget=False)
        ctx['sidebar_total_on_budget'] = on_budget_qs.aggregate(Sum('balance'))['balance__sum'] or 0
        ctx['sidebar_total_off_budget'] = off_budget_qs.aggregate(Sum('balance'))['balance__sum'] or 0
        ctx['sidebar_grand_total'] = ctx['sidebar_total_on_budget'] + ctx['sidebar_total_off_budget']

        ccs = CreditCard.objects.filter(budget=active)
        ctx['sidebar_cc_debt'] = cc_debt = sum(cc.balance for cc in ccs if cc.balance < 0)

        loans = Loan.objects.filter(budget=active, status='active')
        ctx['sidebar_loan_debt'] = loan_debt = loans.aggregate(Sum('remaining_balance'))['remaining_balance__sum'] or 0
        ctx['sidebar_total_debt'] = abs(cc_debt) + loan_debt

    return ctx

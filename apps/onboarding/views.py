from datetime import date
from django.db.models import Sum
from django.http import JsonResponse
from apps.accounts.models import Account
from apps.budgets.models import Budget, MonthlyBudget
from apps.goals.models import Goal
from apps.schedules.models import Schedule
from apps.credit_cards.models import CreditCard
from apps.loans.models import Loan
from .models import OnboardingProfile


def get_or_create_profile(user):
    try:
        return user.onboarding
    except OnboardingProfile.DoesNotExist:
        return OnboardingProfile.objects.create(user=user, step=0)


def onboarding_state(request):
    if not request.user.is_authenticated:
        return JsonResponse({'step': 0, 'step_completed': False, 'ready_to_assign': 0})

    profile = get_or_create_profile(request.user)
    budget_id = request.session.get('active_budget_id')
    ready_to_assign = 0

    if budget_id:
        try:
            budget = Budget.objects.get(id=budget_id, members=request.user)
            total = float(Account.objects.filter(
                budget=budget, on_budget=True
            ).aggregate(Sum('balance'))['balance__sum'] or 0)
            assigned = float(MonthlyBudget.objects.filter(
                category__budget=budget, month=date.today().month, year=date.today().year
            ).aggregate(Sum('budgeted'))['budgeted__sum'] or 0)
            ready_to_assign = round(max(0, total - assigned), 2)
        except Budget.DoesNotExist:
            pass

    return JsonResponse({
        'step': profile.step,
        'step_completed': _check_step_completed(profile.step, budget_id, request.user),
        'ready_to_assign': ready_to_assign,
    })


def advance_step(request):
    if not request.user.is_authenticated:
        return JsonResponse({'error': 'unauthorized'}, status=401)

    profile = get_or_create_profile(request.user)
    profile.step = min(profile.step + 1, 7)
    profile.save()
    return JsonResponse({'step': profile.step})


def _check_step_completed(step, budget_id, user):
    if not budget_id:
        return False if step > 1 else True
    try:
        budget = Budget.objects.get(id=budget_id, members=user)
    except Budget.DoesNotExist:
        return False

    if step >= 7:
        return True
    elif step == 1:
        return True
    elif step == 2:
        return Account.objects.filter(budget=budget).count() >= 1
    elif step == 3:
        total = float(Account.objects.filter(budget=budget, on_budget=True).aggregate(Sum('balance'))['balance__sum'] or 0)
        assigned = float(MonthlyBudget.objects.filter(
            category__budget=budget, month=date.today().month, year=date.today().year
        ).aggregate(Sum('budgeted'))['budgeted__sum'] or 0)
        return total - assigned <= 0.01
    elif step == 4:
        return Goal.objects.filter(category__budget=budget).count() >= 1
    elif step == 5:
        return Schedule.objects.filter(budget=budget).count() >= 1
    elif step == 6:
        return CreditCard.objects.filter(budget=budget).count() >= 1 or Loan.objects.filter(budget=budget).count() >= 1
    return False

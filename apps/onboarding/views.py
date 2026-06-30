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


CONDITIONS = [
    {'id': 'welcome', 'priority': 0, 'title': 'Bienvenido',
     'html': '<div class="text-center"><div class="text-3xl mb-3">💰</div><h3 class="text-base font-bold mb-2">Bienvenido a Chele</h3><p class="text-sm text-chele-text-secondary">Cada peso debe tener un trabajo. Vamos a organizar tu dinero paso a paso.</p></div>',
     'check': lambda b, u: False,  # only shown if step=0
     'check_completed': lambda b, u: True},
    {'id': 'accounts', 'priority': 1, 'title': 'Tus cuentas',
     'target': 'a[href*="/cuentas/crear/"]',
     'html': '<p><strong>Agregá tus cuentas bancarias.</strong> El saldo inicial lo vas a distribuir después.</p>',
     'check': lambda b, u: Account.objects.filter(budget=b).count() == 0,
     'check_completed': lambda b, u: Account.objects.filter(budget=b).count() >= 1},
    {'id': 'assign', 'priority': 2, 'title': 'Asignar dinero',
     'target': '[class*="Por asignar"], [class*="ready"]',
     'html': '<p><strong>Tenés dinero sin trabajo.</strong> Asignalo a categorías hasta que "Por asignar" llegue a $0.</p>',
     'check': lambda b, u: _get_rta(b) > 0,
     'check_completed': lambda b, u: _get_rta(b) <= 0.01},
    {'id': 'goals', 'priority': 3, 'title': 'Metas de ahorro',
     'target': '.category-item, tr.cursor-pointer',
     'html': '<p><strong>¿Querés ahorrar para algo?</strong> Hacé click en una categoría y poné una meta.</p>',
     'check': lambda b, u: Goal.objects.filter(category__budget=b).count() == 0,
     'check_completed': lambda b, u: Goal.objects.filter(category__budget=b).count() >= 1},
    {'id': 'schedules', 'priority': 4, 'title': 'Gastos recurrentes',
     'target': 'a[href*="/programaciones/"]',
     'html': '<p><strong>¿Tenés gastos o ingresos que se repiten?</strong> Programalos y se aplican solos.</p>',
     'check': lambda b, u: Schedule.objects.filter(budget=b).count() == 0,
     'check_completed': lambda b, u: Schedule.objects.filter(budget=b).count() >= 1},
]


def _get_rta(budget):
    total = float(Account.objects.filter(budget=budget, on_budget=True).aggregate(
        Sum('balance')
    )['balance__sum'] or 0)
    assigned = float(MonthlyBudget.objects.filter(
        category__budget=budget, month=date.today().month, year=date.today().year
    ).aggregate(Sum('budgeted'))['budgeted__sum'] or 0)
    return round(max(0, total - assigned), 2)


def get_or_create_profile(user):
    try:
        return user.onboarding
    except OnboardingProfile.DoesNotExist:
        return OnboardingProfile.objects.create(user=user, step=0)


def _get_active_condition(budget, user, step, dismissed=None):
    """Return the first condition that needs attention, or None."""
    if dismissed is None:
        dismissed = []
    if step == 0:
        return 'welcome'
    for c in CONDITIONS:
        if c['id'] in dismissed:
            continue
        if c['check'](budget, user):
            return c['id']
    return None


def onboarding_state(request):
    if not request.user.is_authenticated:
        return JsonResponse({'active': False, 'condition': None, 'ready_to_assign': 0})

    profile = get_or_create_profile(request.user)
    budget_id = request.session.get('active_budget_id')
    rta = 0
    condition_id = None
    condition_completed = False

    if budget_id:
        try:
            budget = Budget.objects.get(id=budget_id, members=request.user)
            rta = _get_rta(budget)
            dismissed = request.session.get('dismissed_conditions', [])
            condition_id = _get_active_condition(budget, request.user, profile.step, dismissed)
            if condition_id:
                for c in CONDITIONS:
                    if c['id'] == condition_id:
                        condition_completed = c['check_completed'](budget, request.user)
                        break
        except Budget.DoesNotExist:
            pass

    return JsonResponse({
        'step': profile.step,
        'active': condition_id is not None and profile.step >= 0,
        'condition': condition_id,
        'condition_completed': condition_completed,
        'ready_to_assign': rta,
    })


def advance_step(request):
    if not request.user.is_authenticated:
        return JsonResponse({'error': 'unauthorized'}, status=401)

    profile = get_or_create_profile(request.user)
    profile.step = min(profile.step + 1, 7)
    profile.save()
    return JsonResponse({'step': profile.step})


def dismiss_condition(request):
    if not request.user.is_authenticated:
        return JsonResponse({'error': 'unauthorized'}, status=401)
    profile = get_or_create_profile(request.user)
    if profile.step == 0:
        profile.step = 1
        profile.save()
    import json
    body = json.loads(request.body) if request.body else {}
    condition_id = body.get('condition')
    if condition_id:
        dismissed = request.session.get('dismissed_conditions', [])
        if condition_id not in dismissed:
            dismissed.append(condition_id)
        request.session['dismissed_conditions'] = dismissed
    return JsonResponse({'status': 'dismissed'})

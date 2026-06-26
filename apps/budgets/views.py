import calendar
import json
from datetime import date, datetime, timedelta
from django.shortcuts import render, redirect, get_object_or_404
from django.urls import reverse
from django.contrib.auth.decorators import login_required
from django.contrib import messages
from django.db.models import Sum
from .models import Budget, BudgetMembership, CategoryGroup, Category, MonthlyBudget
from apps.accounts.models import Account
from apps.transactions.models import Transaction
from apps.schedules.models import Schedule
from apps.schedules.views import process_due_schedules


@login_required
def budget_view(request):
    budget_id = request.session.get('active_budget_id')
    if not budget_id:
        first = Budget.objects.filter(members=request.user).first()
        if first:
            budget_id = str(first.id)
            request.session['active_budget_id'] = budget_id
    try:
        budget = Budget.objects.get(id=budget_id, members=request.user)
    except (Budget.DoesNotExist, ValueError):
        return redirect('budget_create')
    process_due_schedules(budget_id, budget)
    today = date.today()
    month = int(request.GET.get('mes', today.month))
    year = int(request.GET.get('anio', today.year))

    groups = CategoryGroup.objects.filter(budget=budget).prefetch_related('categories')
    accounts = Account.objects.filter(budget=budget)
    total_balance = float(accounts.filter(on_budget=True).aggregate(Sum('balance'))['balance__sum'] or 0)

    range_val = int(request.GET.get('rango', 0))
    visible_months = []
    for i in range(-range_val, range_val + 1):
        m = month + i
        y = year
        while m > 12:
            m -= 12
            y += 1
        while m < 1:
            m += 12
            y -= 1
        visible_months.append({'month': m, 'year': y, 'active': m == month and y == year})

    for i, vm in enumerate(visible_months):
        month_income = float(Transaction.objects.filter(
            budget=budget, date__month=vm['month'], date__year=vm['year'],
            amount__gt=0,
        ).aggregate(Sum('amount'))['amount__sum'] or 0)

        # Include scheduled income for this month (so user can distribute future income)
        scheduled_income = float(Schedule.objects.filter(
            budget=budget, is_active=True, direction='income',
            next_date__month=vm['month'], next_date__year=vm['year'],
        ).aggregate(Sum('amount'))['amount__sum'] or 0)
        month_income += scheduled_income

        month_expenses = float(Transaction.objects.filter(
            budget=budget, date__month=vm['month'], date__year=vm['year'],
            amount__lt=0,
        ).aggregate(Sum('amount'))['amount__sum'] or 0)

        total_budgeted = float(MonthlyBudget.objects.filter(
            category__budget=budget, month=vm['month'], year=vm['year'],
        ).aggregate(Sum('budgeted'))['budgeted__sum'] or 0)

        vm['total_budgeted'] = round(total_budgeted, 2)
        vm['month_spent'] = round(abs(month_expenses), 2)
        vm['month_income'] = round(month_income, 2)

        # Calculate overspent from previous month
        # If last month's spending > last month's budgeted, the difference is overspent
        vm['overspent_prev'] = 0
        if i == 0:
            # First visible month: check the month before it
            pm = vm['month'] - 1
            py = vm['year']
            if pm < 1:
                pm = 12
                py -= 1
            prev_exp = float(Transaction.objects.filter(
                budget=budget, date__month=pm, date__year=py, amount__lt=0,
            ).aggregate(Sum('amount'))['amount__sum'] or 0)
            prev_bud = float(MonthlyBudget.objects.filter(
                category__budget=budget, month=pm, year=py,
            ).aggregate(Sum('budgeted'))['budgeted__sum'] or 0)
            vm['overspent_prev'] = round(max(0, abs(prev_exp) - prev_bud), 2)

        # Calculate leftover from previous month
        # For the first month ever, this starts at 0 and we use total balance as base
        vm['carried_over'] = 0
        if i > 0:
            prev = visible_months[i - 1]
            vm['carried_over'] = round(
                prev['month_income'] - prev['total_budgeted'] - prev['month_spent'],
                2
            )

        # Available funds = carried_over + this month's income
        vm['available_funds'] = round(vm['carried_over'] + month_income, 2)

        # For next month = what's left after this month
        vm['for_next_month'] = round(
            vm['available_funds'] - vm['total_budgeted'] - vm['overspent_prev'],
            2
        )

    # To Budget = Income this month + Remanente - Budgeted - Overspent
    active = [vm for vm in visible_months if vm['active']]
    if active:
        a = active[0]
        available_to_budget = round(
            a['month_income'] + a['carried_over'] - a['total_budgeted'] - a['overspent_prev'],
            2
        )
    else:
        available_to_budget = 0

    # For the first-ever month: use total balance as the starting point
    # (before any carry-over logic kicks in)
    if available_to_budget == 0 and total_balance > 0:
        available_to_budget = round(total_balance, 2)

    # Auto-expire snooze for past months
    from apps.goals.models import Goal as GoalModel
    for g in GoalModel.objects.filter(snooze_month__isnull=False, snooze_year__isnull=False):
        if g.snooze_year < year or (g.snooze_year == year and g.snooze_month < month):
            g.snooze_month = None
            g.snooze_year = None
            g.save()

    # Ready to Assign
    ready_to_assign = round(max(0, total_balance - sum(
        float(MonthlyBudget.objects.filter(
            category__budget=budget, month=month, year=year
        ).aggregate(Sum('budgeted'))['budgeted__sum'] or 0)
        for m in [month] for y in [year]
    )), 2)

    # Underfunded categories (Targets)
    from apps.goals.services import TargetService
    ts = TargetService(budget, month, year)
    underfunded_categories = ts.list_underfunded()
    total_underfunded = round(sum(u['deficit'] for u in underfunded_categories), 2)

    # Cost to be me
    cost_to_be_me = 0
    for g in GoalModel.objects.filter(category__budget=budget, is_completed=False).select_related('category'):
        cost_to_be_me += ts.calculate_underfunded(g) if g.goal_type in ('monthly', 'yearly', 'true_expense') else 0
    if cost_to_be_me == 0:
        cost_to_be_me = total_underfunded

    # Expected Income (average last 3 months)
    from django.db.models import Avg
    expected_income = float(Transaction.objects.filter(
        budget=budget, amount__gt=0,
        date__gte=date(today.year, today.month, 1) - timedelta(days=90)
    ).aggregate(avg=Avg('amount'))['avg'] or 0)

    is_over_budget = cost_to_be_me > expected_income and expected_income > 0

    # Rollover per category (previous month)
    prev_month = month - 1
    prev_year = year
    if prev_month < 1:
        prev_month = 12
        prev_year -= 1
    rollover_categories = []
    from apps.budgets.models import Category as BCategory
    for cat in BCategory.objects.filter(budget=budget, is_hidden=False):
        pts = TargetService(budget, prev_month, prev_year)
        avail = pts.get_category_available(str(cat.id))
        if avail > 0.01:
            rollover_categories.append({'category_id': str(cat.id), 'name': cat.name, 'balance': round(avail, 2)})

    month_names = {1:'Ene',2:'Feb',3:'Mar',4:'Abr',5:'May',6:'Jun',
                   7:'Jul',8:'Ago',9:'Sep',10:'Oct',11:'Nov',12:'Dic'}
    for vm in visible_months:
        vm['label'] = f"{month_names[vm['month']]} {vm['year']}"

    budget_data = []
    for group in groups:
        categories_data = []
        for cat in group.categories.filter(is_hidden=False):
            cat_months = []
            for vm in visible_months:
                mb, _ = MonthlyBudget.objects.get_or_create(
                    category=cat, month=vm['month'], year=vm['year']
                )
                spent = Transaction.objects.filter(
                    budget=budget, category=cat,
                    date__month=vm['month'], date__year=vm['year'],
                ).aggregate(Sum('amount'))['amount__sum'] or 0

                # Include rollover from previous month
                prev_month = vm['month'] - 1
                prev_year = vm['year']
                if prev_month < 1:
                    prev_month = 12
                    prev_year -= 1
                prev_pts = TargetService(budget, prev_month, prev_year)
                rollover = prev_pts.get_category_available(str(cat.id))
                if rollover < 0:
                    rollover = 0

                balance = float(mb.budgeted) + rollover - float(abs(spent) if spent < 0 else spent)

                cat_months.append({
                    'budgeted': mb.budgeted,
                    'spent': abs(spent) if spent < 0 else spent,
                    'balance': balance,
                    'rollover': round(rollover, 2),
                    'month': vm['month'],
                    'year': vm['year'],
                })

            has_activity = Transaction.objects.filter(
                budget=budget, category=cat
            ).exists()

            # Check underfunded/overspent status
            balance = float(cat_months[0]['balance']) if cat_months else 0
            is_overspent = balance < 0
            underfunded_amt = 0
            for uc in underfunded_categories:
                if uc['category_id'] == str(cat.id):
                    underfunded_amt = uc['deficit']
                    break

            categories_data.append({
                'category': cat,
                'months': cat_months,
                'is_income': group.is_income,
                'can_delete': not has_activity,
                'balance': balance,
                'is_overspent': is_overspent,
                'underfunded': underfunded_amt > 0,
                'underfunded_amount': underfunded_amt,
            })
        budget_data.append({'group': group, 'categories': categories_data})

    return render(request, 'budgets/budget_view.html', {
        'active_budget': budget,
        'groups': budget_data,
        'months': visible_months,
        'available_to_budget': available_to_budget,
        'total_balance': total_balance,
        'current_month': month,
        'current_year': year,
        'ready_to_assign': ready_to_assign,
        'underfunded_categories': underfunded_categories,
        'total_underfunded': total_underfunded,
        'cost_to_be_me': round(cost_to_be_me, 2),
        'expected_income': round(expected_income, 2),
        'is_over_budget': is_over_budget,
        'rollover_categories': rollover_categories,
    })


@login_required
def budget_create(request):
    if request.method == 'POST':
        name = request.POST.get('name')
        budget = Budget.objects.create(name=name, owner=request.user)
        BudgetMembership.objects.create(user=request.user, budget=budget, role='owner', accepted_at=datetime.now())
        request.session['active_budget_id'] = str(budget.id)

        default_groups = [
            ('Ingresos', True),
            ('Gastos Fijos', False),
            ('Gastos Diarios', False),
            ('Ahorro', False),
        ]
        for i, (gname, is_inc) in enumerate(default_groups):
            group = CategoryGroup.objects.create(budget=budget, name=gname, sort_order=i, is_income=is_inc)

        default_categories = {
            'Ingresos': ['Sueldo', 'Freelance', 'Inversiones', 'Otros ingresos'],
            'Gastos Fijos': ['Alquiler', 'Servicios', 'Suscripciones', 'True Expenses'],
            'Gastos Diarios': ['Comida', 'Transporte', 'Salidas', 'Salud'],
            'Ahorro': ['Fondo de emergencia', 'Vacaciones'],
        }
        for gname, cats in default_categories.items():
            group = CategoryGroup.objects.get(budget=budget, name=gname)
            for j, cname in enumerate(cats):
                Category.objects.create(budget=budget, group=group, name=cname, sort_order=j)

        messages.success(request, f'Presupuesto "{name}" creado correctamente.')
        return redirect('budget_view')
    return render(request, 'budgets/budget_create.html')


@login_required
def switch_budget(request):
    if request.method == 'POST':
        budget_id = request.POST.get('budget_id')
        if Budget.objects.filter(id=budget_id, members=request.user).exists():
            request.session['active_budget_id'] = budget_id
    return redirect('budget_view')


@login_required
def category_inspector(request, id):
    budget_id = request.session.get('active_budget_id')
    cat = get_object_or_404(Category, id=id, budget_id=budget_id)
    from datetime import date
    today = date.today()
    month = int(request.GET.get('mes', today.month))
    year = int(request.GET.get('anio', today.year))

    mb = MonthlyBudget.objects.filter(category=cat, month=month, year=year).first()
    budgeted = float(mb.budgeted) if mb else 0

    spent = abs(float(Transaction.objects.filter(
        category=cat, date__month=month, date__year=year
    ).aggregate(Sum('amount'))['amount__sum'] or 0))

    available = budgeted - spent

    from apps.goals.services import TargetService
    ts = TargetService(cat.budget, month, year)
    target_info = {}
    goal = cat.goals.filter(is_completed=False).first()
    if goal:
        underfunded = ts.calculate_underfunded(goal)
        target_info = {
            'type': goal.goal_type,
            'amount': float(goal.amount),
            'refill_up_to': goal.refill_up_to,
            'underfunded': round(underfunded, 2),
        }

    # 3-month average
    from django.db.models import Avg
    ninety_days_ago = date(today.year, today.month, 1) - timedelta(days=90)
    avg_3m = Transaction.objects.filter(
        category=cat,
        date__gte=ninety_days_ago
    ).aggregate(avg=Avg('amount'))['avg'] or 0

    return HttpResponse(json.dumps({
        'category_id': str(cat.id),
        'category_name': cat.name,
        'available': round(available, 2),
        'assigned': round(budgeted, 2),
        'activity': round(spent, 2),
        'target': target_info,
        'average_spending_3m': round(abs(float(avg_3m)), 2),
        'overspent': available < -0.01,
    }), content_type='application/json')


@login_required
def assign_funds(request):
    if request.method == 'POST':
        budget_id = request.session.get('active_budget_id')
        category_id = request.POST.get('category_id')
        amount = float(request.POST.get('amount', 0))
        month = int(request.POST.get('month'))
        year = int(request.POST.get('year'))

        # Zero-sum validation
        total_on_budget = float(Account.objects.filter(budget_id=budget_id, on_budget=True).aggregate(Sum('balance'))['balance__sum'] or 0)
        if total_on_budget > 0:
            total_assigned = float(MonthlyBudget.objects.filter(category__budget_id=budget_id, month=month, year=year).exclude(category_id=category_id).aggregate(Sum('budgeted'))['budgeted__sum'] or 0)
            if total_assigned + amount > total_on_budget:
                messages.error(request, f'No tienes suficientes fondos disponibles. Ready to Assign: ${total_on_budget - total_assigned:.2f}')
                return redirect('budget_view')

        category = get_object_or_404(Category, id=category_id, budget_id=budget_id)
        mb, _ = MonthlyBudget.objects.get_or_create(category=category, month=month, year=year)
        mb.budgeted = amount
        mb.save()
        messages.success(request, 'Fondos asignados correctamente.')
    return redirect('budget_view')


@login_required
def reorder_categories(request):
    if request.method == 'POST':
        from django.http import JsonResponse
        import json as j
        data = j.loads(request.body)
        order = data.get('order', [])
        for item in order:
            cid = item.get('id')
            sort = item.get('sort_order', 0)
            gid = item.get('group_id')
            Category.objects.filter(id=cid).update(sort_order=sort)
            if gid:
                Category.objects.filter(id=cid).update(group_id=gid)
        return JsonResponse({'status': 'ok'})
    return redirect('budget_view')


@login_required
def move_funds(request):
    if request.method == 'POST':
        budget_id = request.session.get('active_budget_id')
        from_category_id = request.POST.get('from_category')
        to_category_id = request.POST.get('to_category')
        amount = float(request.POST.get('amount'))
        month = int(request.POST.get('month'))
        year = int(request.POST.get('year'))

        from_mb, _ = MonthlyBudget.objects.get_or_create(
            category_id=from_category_id, month=month, year=year
        )
        to_mb, _ = MonthlyBudget.objects.get_or_create(
            category_id=to_category_id, month=month, year=year
        )
        from_mb.budgeted = float(from_mb.budgeted) - amount
        to_mb.budgeted = float(to_mb.budgeted) + amount
        from_mb.save()
        to_mb.save()
        messages.success(request, 'Fondos transferidos.')
    return redirect('budget_view')


@login_required
def auto_assign(request):
    if request.method == 'POST':
        budget_id = request.session.get('active_budget_id')
        budget = get_object_or_404(Budget, id=budget_id)
        month = int(request.POST.get('month', 0))
        year = int(request.POST.get('year', 0))
        if not month or not year:
            messages.error(request, 'Mes y año requeridos.')
            return redirect('budget_view')

        from apps.goals.services import TargetService
        ts = TargetService(budget, month, year)
        underfunded = ts.list_underfunded()

        total_on_budget = float(Account.objects.filter(budget=budget, on_budget=True).aggregate(Sum('balance'))['balance__sum'] or 0)
        total_assigned = float(MonthlyBudget.objects.filter(category__budget=budget, month=month, year=year).aggregate(Sum('budgeted'))['budgeted__sum'] or 0)
        available = max(0, total_on_budget - total_assigned)

        assigned_count = 0
        total_assigned_amt = 0.0
        for u in underfunded:
            if available <= 0:
                break
            amt = min(u['deficit'], available)
            mb, _ = MonthlyBudget.objects.get_or_create(category_id=u['category_id'], month=month, year=year)
            mb.budgeted = float(mb.budgeted) + amt
            mb.save()
            assigned_count += 1
            total_assigned_amt += amt
            available -= amt

        messages.success(request, f'Fondos asignados: {assigned_count} categorías, ${total_assigned_amt:.2f}')
    return redirect('budget_view')


@login_required
def cover_overspending(request):
    if request.method == 'POST':
        budget_id = request.session.get('active_budget_id')
        from_category_id = request.POST.get('from_category')
        to_category_id = request.POST.get('to_category')
        amount = float(request.POST.get('amount', 0))
        month = int(request.POST.get('month'))
        year = int(request.POST.get('year'))

        # Validate source has enough available
        budget = get_object_or_404(Budget, id=budget_id)
        from apps.goals.services import TargetService
        ts = TargetService(budget, month, year)
        from_avail = None
        try:
            from_avail = ts.get_category_available(from_category_id)
        except Exception:
            from_avail = 0

        if from_avail < amount:
            messages.error(request, 'Fondos insuficientes en la categoría origen.')
            return redirect('budget_view')

        from_mb, _ = MonthlyBudget.objects.get_or_create(
            category_id=from_category_id, month=month, year=year
        )
        to_mb, _ = MonthlyBudget.objects.get_or_create(
            category_id=to_category_id, month=month, year=year
        )
        from_mb.budgeted = float(from_mb.budgeted) - amount
        to_mb.budgeted = float(to_mb.budgeted) + amount
        from_mb.save()
        to_mb.save()
        messages.success(request, 'Sobregasto cubierto correctamente.')
    return redirect('budget_view')


@login_required
def hold_for_next_month(request):
    if request.method == 'POST':
        budget_id = request.session.get('active_budget_id')
        amount = float(request.POST.get('amount', 0))
        request.session[f'held_{budget_id}'] = request.session.get(f'held_{budget_id}', 0) + amount
        messages.success(request, 'Fondos reservados para el próximo mes.')
    return redirect('budget_view')


@login_required
def copy_budget(request):
    if request.method == 'POST':
        budget_id = request.session.get('active_budget_id')
        from_month = int(request.POST.get('from_month'))
        from_year = int(request.POST.get('from_year'))
        to_month = int(request.POST.get('to_month'))
        to_year = int(request.POST.get('to_year'))

        from_budgets = MonthlyBudget.objects.filter(
            category__budget_id=budget_id,
            month=from_month,
            year=from_year,
        )
        for fb in from_budgets:
            tb, _ = MonthlyBudget.objects.get_or_create(
                category=fb.category, month=to_month, year=to_year
            )
            tb.budgeted = fb.budgeted
            tb.save()
        messages.success(request, 'Presupuesto copiado.')
    return redirect('budget_view')


@login_required
def category_create(request):
    if request.method == 'POST':
        budget_id = request.session.get('active_budget_id')
        group_id = request.POST.get('group_id')
        name = request.POST.get('name')
        Category.objects.create(
            budget_id=budget_id,
            group_id=group_id,
            name=name,
        )
        messages.success(request, 'Categoría creada.')
    return redirect(reverse('budget_view') + '?editar=1')


@login_required
def category_edit(request, id):
    cat = get_object_or_404(Category, id=id)
    if request.method == 'POST':
        cat.name = request.POST.get('name')
        cat.save()
        messages.success(request, 'Categoría actualizada.')
    return redirect(reverse('budget_view') + '?editar=1')


@login_required
def category_delete(request, id):
    cat = get_object_or_404(Category, id=id)
    if request.method == 'POST':
        has_txns = Transaction.objects.filter(category=cat).exists()
        if has_txns:
            cat.is_hidden = True
            cat.save()
            messages.success(request, f'Categoría "{cat.name}" oculta.')
        else:
            cat.delete()
            messages.success(request, f'Categoría "{cat.name}" eliminada.')
    return redirect(reverse('budget_view') + '?editar=1')


@login_required
@login_required
def category_group_create(request):
    if request.method == 'POST':
        budget_id = request.session.get('active_budget_id')
        name = request.POST.get('name')
        CategoryGroup.objects.create(budget_id=budget_id, name=name)
        messages.success(request, 'Grupo creado.')
    return redirect(reverse('budget_view') + '?editar=1')


@login_required
def category_group_edit(request, id):
    group = get_object_or_404(CategoryGroup, id=id)
    if request.method == 'POST':
        group.name = request.POST.get('name', group.name)
        group.save()
        messages.success(request, f'Grupo renombrado a "{group.name}".')
    return redirect(reverse('budget_view') + '?editar=1')


@login_required
def category_group_delete(request, id):
    group = get_object_or_404(CategoryGroup, id=id)
    if request.method == 'POST':
        if group.is_income:
            messages.error(request, 'No se puede eliminar el grupo de Ingresos.')
            return redirect(reverse('budget_view') + '?editar=1')
        has_cats = Category.objects.filter(group=group).exists()
        has_txns = Transaction.objects.filter(category__group=group).exists()
        if has_txns:
            messages.error(request, f'No se puede eliminar "{group.name}" porque tiene categorías con transacciones.')
        elif has_cats:
            Category.objects.filter(group=group).delete()
            group.delete()
            messages.success(request, f'Grupo "{group.name}" y sus categorías eliminados.')
        else:
            group.delete()
            messages.success(request, f'Grupo "{group.name}" eliminado.')
    return redirect(reverse('budget_view') + '?editar=1')


def register(request):
    if request.method == 'POST':
        from django.contrib.auth.models import User
        email = request.POST.get('email')
        password = request.POST.get('password')
        name = request.POST.get('name')
        if User.objects.filter(email=email).exists():
            messages.error(request, 'El email ya está registrado.')
            return render(request, 'registration/register.html')
        user = User.objects.create_user(username=email, email=email, password=password)
        user.first_name = name
        user.save()
        from django.contrib.auth import login
        login(request, user)
        return redirect('budget_create')
    return render(request, 'registration/register.html')

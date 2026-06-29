from datetime import date, timedelta
from django.shortcuts import render, redirect, get_object_or_404
from django.contrib.auth.decorators import login_required
from django.contrib import messages
from .models import Schedule
from apps.transactions.models import Transaction
from apps.accounts.models import Account


def _skip_weekend(d):
    if d.weekday() == 5:
        return d + timedelta(days=2)
    if d.weekday() == 6:
        return d + timedelta(days=1)
    return d


def _apply_before_weekend(d):
    if d.weekday() == 5:
        return d - timedelta(days=1)
    if d.weekday() == 6:
        return d - timedelta(days=2)
    return d


def _adjust_weekend(date_val, skip_weekends, apply_before):
    if apply_before:
        return _apply_before_weekend(date_val)
    if skip_weekends:
        return _skip_weekend(date_val)
    return date_val


def _advance_date(schedule):
    next_date = schedule.next_date
    freqs = {
        'weekly': timedelta(days=7),
        'biweekly': timedelta(days=14),
        'monthly': timedelta(days=31),
        'quarterly': timedelta(days=92),
        'yearly': timedelta(days=365),
    }
    delta = freqs.get(schedule.frequency, timedelta(days=31))
    new_date = next_date + delta
    new_date = _adjust_weekend(new_date, schedule.skip_weekends, schedule.apply_before_weekend)
    return new_date


def process_due_schedules(budget_id, budget):
    today = date.today()
    due = Schedule.objects.filter(budget_id=budget_id, is_active=True, next_date__lte=today)
    count = 0
    for s in due.select_related('account'):
        actual_date = _adjust_weekend(s.next_date, s.skip_weekends, s.apply_before_weekend)
        raw = float(s.amount)
        account = s.account
        if s.direction == 'income':
            amount = abs(raw)
            account.balance = float(account.balance) + amount
        else:
            amount = -abs(raw)
            account.balance = float(account.balance) - abs(raw)
        Transaction.objects.create(
            budget_id=budget_id, account=account,
            date=actual_date, amount=amount,
            payee=s.payee, category=s.category,
            notes=s.notes or '',
        )
        account.save()
        s.next_date = _advance_date(s)
        s.save()
        count += 1
    return count


@login_required
def schedules_list(request):
    budget_id = request.session.get('active_budget_id')
    processed = process_due_schedules(budget_id, None)
    schedules = Schedule.objects.filter(budget_id=budget_id).select_related('payee', 'category', 'account')
    today = date.today()
    for s in schedules:
        s.due_soon = s.next_date and s.next_date <= today + timedelta(days=3)
    return render(request, 'schedules/schedule_list.html', {'schedules': schedules})


@login_required
def schedule_create(request):
    if request.method == 'POST':
        budget_id = request.session.get('active_budget_id')
        next_date_str = request.POST.get('next_date')
        next_date = date.fromisoformat(next_date_str) if next_date_str else date.today()
        skip = request.POST.get('skip_weekends') == 'on'
        before = request.POST.get('apply_before_weekend') == 'on'
        if skip or before:
            next_date = _adjust_weekend(next_date, skip, before)
        schedule = Schedule.objects.create(
            budget_id=budget_id,
            payee_id=request.POST.get('payee_id') or None,
            category_id=request.POST.get('category_id') or None,
            account_id=request.POST.get('account_id'),
            amount=request.POST.get('amount'),
            frequency=request.POST.get('frequency'),
            next_date=next_date,
            notes=request.POST.get('notes', ''),
            skip_weekends=skip,
            apply_before_weekend=before,
            direction=request.POST.get('direction', 'expense'),
        )
        messages.success(request, 'Programación creada.')
        return redirect('schedules_list')
    return render(request, 'schedules/schedule_form.html')


@login_required
def schedule_edit(request, id):
    schedule = get_object_or_404(Schedule, id=id)
    if request.method == 'POST':
        schedule.amount = request.POST.get('amount')
        schedule.frequency = request.POST.get('frequency')
        schedule.next_date = request.POST.get('next_date')
        schedule.notes = request.POST.get('notes', '')
        schedule.skip_weekends = request.POST.get('skip_weekends') == 'on'
        schedule.apply_before_weekend = request.POST.get('apply_before_weekend') == 'on'
        schedule.direction = request.POST.get('direction', 'expense')
        schedule.save()
        messages.success(request, 'Programación actualizada.')
        return redirect('schedules_list')
    return render(request, 'schedules/schedule_form.html', {'schedule': schedule})


@login_required
def schedule_apply_now(request, id):
    schedule = get_object_or_404(Schedule, id=id)
    if request.method == 'POST':
        actual_date = _adjust_weekend(date.today(), schedule.skip_weekends, schedule.apply_before_weekend)
        raw = float(schedule.amount)
        account = schedule.account
        if schedule.direction == 'income':
            amount = abs(raw)
            account.balance = float(account.balance) + amount
        else:
            amount = -abs(raw)
            account.balance = float(account.balance) - abs(raw)
        Transaction.objects.create(
            budget_id=schedule.budget_id, account=account,
            date=actual_date, amount=amount,
            payee=schedule.payee, category=schedule.category,
            notes=schedule.notes or '',
        )
        account.save()
        schedule.next_date = _advance_date(schedule)
        schedule.save()
        messages.success(request, f'Programación "{schedule.payee.name if schedule.payee else "—"}" aplicada.')
    return redirect(request.META.get('HTTP_REFERER', 'schedules_list'))


@login_required
def schedule_delete(request, id):
    schedule = get_object_or_404(Schedule, id=id)
    if request.method == 'POST':
        schedule.delete()
        messages.success(request, 'Programación eliminada.')
    return redirect('schedules_list')

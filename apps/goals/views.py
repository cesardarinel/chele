from datetime import date
from django.shortcuts import render, redirect, get_object_or_404
from django.contrib.auth.decorators import login_required
from django.contrib import messages
from django.db.models import Sum
from dateutil.relativedelta import relativedelta
from .models import Goal
from apps.budgets.models import MonthlyBudget


@login_required
def goal_list(request):
    budget_id = request.session.get('active_budget_id')
    goals = Goal.objects.filter(category__budget_id=budget_id).select_related('category')
    return render(request, 'goals/goal_list.html', {'goals': goals})


@login_required
def goal_create(request):
    if request.method == 'POST':
        budget_id = request.session.get('active_budget_id')
        goal = Goal.objects.create(
            category_id=request.POST.get('category_id'),
            goal_type=request.POST.get('goal_type'),
            amount=request.POST.get('amount'),
            target_date=request.POST.get('target_date') or None,
            frequency=request.POST.get('frequency', 12),
        )
        _apply_goal(goal)
        messages.success(request, 'Meta creada.')
        return redirect('budget_view')
    return render(request, 'goals/goal_form.html')


def _apply_goal(goal):
    today = date.today()
    mb, _ = MonthlyBudget.objects.get_or_create(
        category=goal.category, month=today.month, year=today.year
    )
    if goal.goal_type == 'monthly':
        mb.budgeted = goal.amount
    elif goal.goal_type == 'target_balance':
        spent = goal.category.transactions.aggregate(Sum('amount'))['amount__sum'] or 0
        needed = float(goal.amount) - abs(float(spent))
        if needed > 0:
            mb.budgeted = needed
    elif goal.goal_type == 'target_date' and goal.target_date:
        delta = relativedelta(goal.target_date, today)
        months_remaining = delta.years * 12 + delta.months
        if months_remaining > 0:
            mb.budgeted = float(goal.amount) / months_remaining
    elif goal.goal_type == 'true_expense':
        mb.budgeted = float(goal.amount) / goal.frequency
    mb.save()


@login_required
def goal_edit(request, id):
    goal = get_object_or_404(Goal, id=id)
    if request.method == 'POST':
        goal.amount = request.POST.get('amount', goal.amount)
        goal.target_date = request.POST.get('target_date') or None
        goal.save()
        _apply_goal(goal)
        messages.success(request, 'Meta actualizada.')
        return redirect('goal_list')
    return render(request, 'goals/goal_form.html', {'goal': goal})


@login_required
def goal_delete(request, id):
    goal = get_object_or_404(Goal, id=id)
    if request.method == 'POST':
        goal.delete()
        messages.success(request, 'Meta eliminada.')
    return redirect('goal_list')

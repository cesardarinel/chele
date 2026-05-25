from datetime import datetime
from django.shortcuts import redirect
from django.contrib.auth.decorators import login_required
from django.contrib import messages
from django.http import JsonResponse
from .models import SyncLog
from apps.budgets.models import Budget, Category, MonthlyBudget
from apps.accounts.models import Account
from apps.transactions.models import Transaction
from apps.payees.models import Payee


@login_required
def sync_now(request):
    budget_id = request.session.get('active_budget_id')
    push_count = sync_push(request, budget_id)
    pull_data = sync_pull_data(request, budget_id)
    messages.success(request, f'Sincronizado. {push_count} cambios enviados.')
    return redirect(request.META.get('HTTP_REFERER', 'budget_view'))


def sync_push(request, budget_id=None):
    if not budget_id and request.method == 'POST':
        budget_id = request.POST.get('budget_id')
    budget_id = budget_id or request.session.get('active_budget_id')

    pending = SyncLog.objects.filter(budget_id=budget_id, action='pending')
    count = pending.count()

    for log in pending:
        SyncLog.objects.create(
            budget_id=budget_id,
            user=request.user,
            entity_type=log.entity_type,
            entity_id=log.entity_id,
            action='synced',
            payload=log.payload,
        )
        log.delete()

    if request.method == 'POST':
        return JsonResponse({'synced': count})
    return 0


def sync_pull_data(request, budget_id):
    synced = SyncLog.objects.filter(budget_id=budget_id, action='synced').order_by('created_at')
    data = []
    for log in synced:
        data.append({
            'entity_type': log.entity_type,
            'entity_id': str(log.entity_id),
            'payload': log.payload,
        })
    return data


@login_required
def sync_pull(request):
    budget_id = request.session.get('active_budget_id')
    data = sync_pull_data(request, budget_id)
    return JsonResponse({'changes': data})

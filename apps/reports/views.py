from apps.credit_cards.models import CreditCard
from datetime import date, timedelta
from django.shortcuts import render
from django.contrib.auth.decorators import login_required
from django.db.models import Sum, Q
from apps.accounts.models import Account
from apps.budgets.models import Budget, CategoryGroup, Category, MonthlyBudget
from apps.transactions.models import Transaction


@login_required
def reports_dashboard(request):
    return render(request, 'reports/dashboard.html')


@login_required
def net_worth_report(request):
    budget_id = request.session.get('active_budget_id')
    accounts = Account.objects.filter(budget_id=budget_id)
    assets = accounts.aggregate(Sum('balance'))['balance__sum'] or 0
    liabilities = CreditCard.objects.filter(budget_id=budget_id).aggregate(Sum('balance'))['balance__sum'] or 0
    net_worth = assets - abs(liabilities)
    return render(request, 'reports/net_worth.html', {
        'assets': assets,
        'liabilities': abs(liabilities),
        'net_worth': net_worth,
    })


@login_required
def cash_flow_report(request):
    budget_id = request.session.get('active_budget_id')
    today = date.today()
    months_data = []
    for i in range(5, -1, -1):
        m = today.month - i
        y = today.year
        while m < 1:
            m += 12
            y -= 1
        income = Transaction.objects.filter(
            budget_id=budget_id,
            date__month=m, date__year=y,
            amount__gt=0,
        ).aggregate(Sum('amount'))['amount__sum'] or 0
        expenses = Transaction.objects.filter(
            budget_id=budget_id,
            date__month=m, date__year=y,
            amount__lt=0,
        ).aggregate(Sum('amount'))['amount__sum'] or 0
        months_data.append({
            'month': f'{m}/{y}',
            'income': income,
            'expenses': abs(expenses),
            'net': income - abs(expenses),
        })
    return render(request, 'reports/cash_flow.html', {'months': months_data})


@login_required
def budget_vs_reality_report(request):
    budget_id = request.session.get('active_budget_id')
    today = date.today()
    month = int(request.GET.get('mes', today.month))
    year = int(request.GET.get('anio', today.year))
    groups = CategoryGroup.objects.filter(budget_id=budget_id).prefetch_related('categories')
    data = []
    for group in groups:
        cats = []
        for cat in group.categories.all():
            mb = MonthlyBudget.objects.filter(category=cat, month=month, year=year).first()
            budgeted = float(mb.budgeted) if mb else 0
            spent = float(Transaction.objects.filter(
                budget_id=budget_id, category=cat,
                date__month=month, date__year=year,
            ).aggregate(Sum('amount'))['amount__sum'] or 0)
            cats.append({
                'name': cat.name,
                'budgeted': budgeted,
                'spent': abs(spent) if spent < 0 else spent,
                'difference': budgeted - (abs(spent) if spent < 0 else spent),
            })
        data.append({'group': group.name, 'categories': cats})
    return render(request, 'reports/budget_vs_reality.html', {'data': data, 'month': month, 'year': year})

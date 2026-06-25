from datetime import date
from django.shortcuts import render, redirect, get_object_or_404
from django.contrib.auth.decorators import login_required
from django.contrib import messages
from django.db.models import Sum
from .models import Account
from apps.transactions.models import Transaction
from apps.credit_cards.models import CreditCard
from apps.loans.models import Loan


@login_required
def account_list(request):
    budget_id = request.session.get('active_budget_id')
    accounts = Account.objects.filter(budget_id=budget_id)
    on_budget = accounts.filter(on_budget=True)
    off_budget = accounts.filter(on_budget=False)

    total_on_budget = on_budget.aggregate(Sum('balance'))['balance__sum'] or 0
    total_off_budget = off_budget.aggregate(Sum('balance'))['balance__sum'] or 0
    grand_total = total_on_budget + total_off_budget

    credit_cards = CreditCard.objects.filter(budget_id=budget_id)
    cc_debt_cards = [cc for cc in credit_cards if cc.balance < 0]
    total_cc_debt = sum(cc.balance for cc in cc_debt_cards)

    loans = Loan.objects.filter(budget_id=budget_id, status='active')
    total_loan_debt = loans.aggregate(Sum('remaining_balance'))['remaining_balance__sum'] or 0
    total_loan_debt = -total_loan_debt

    total_debt = total_cc_debt + total_loan_debt

    return render(request, 'accounts/account_list.html', {
        'on_budget': on_budget,
        'off_budget': off_budget,
        'total_on_budget': total_on_budget,
        'total_off_budget': total_off_budget,
        'grand_total': grand_total,
        'cc_debt_cards': cc_debt_cards,
        'total_cc_debt': total_cc_debt,
        'loans': loans,
        'total_loan_debt': total_loan_debt,
        'total_debt': total_debt,
    })


@login_required
def account_create(request):
    if request.method == 'POST':
        budget_id = request.session.get('active_budget_id')
        name = request.POST.get('name')
        balance = request.POST.get('balance', 0)
        on_budget = request.POST.get('on_budget') == 'on'
        start_date = request.POST.get('start_date', '')
        account = Account.objects.create(
            budget_id=budget_id, name=name,
            balance=balance, on_budget=on_budget,
        )
        if balance and float(balance) != 0:
            Transaction.objects.create(
                budget_id=budget_id, account=account,
                date=start_date or date.today(),
                amount=balance, notes='Saldo inicial',
            )
        messages.success(request, f'Cuenta "{name}" creada.')
        return redirect('budget_view')
    return render(request, 'accounts/account_form.html')


@login_required
def account_detail(request, id):
    budget_id = request.session.get('active_budget_id')
    account = get_object_or_404(Account, id=id, budget_id=budget_id)
    transactions = Transaction.objects.filter(account=account).select_related('payee', 'category')
    return render(request, 'accounts/account_detail.html', {
        'account': account,
        'transactions': transactions,
    })


@login_required
def account_edit(request, id):
    account = get_object_or_404(Account, id=id)
    if request.method == 'POST':
        account.name = request.POST.get('name', account.name)
        account.on_budget = request.POST.get('on_budget') == 'on'
        account.save()
        messages.success(request, 'Cuenta actualizada.')
        return redirect('account_detail', id=account.id)
    return render(request, 'accounts/account_form.html', {'account': account})


@login_required
def account_delete(request, id):
    account = get_object_or_404(Account, id=id)
    if request.method == 'POST':
        name = account.name
        account.delete()
        messages.success(request, f'Cuenta "{name}" eliminada.')
        return redirect('budget_view')
    return render(request, 'accounts/account_confirm_delete.html', {'account': account})

from django.shortcuts import render, redirect, get_object_or_404
from django.contrib.auth.decorators import login_required
from django.contrib import messages
from django.db.models import Sum
from .models import Transaction
from apps.accounts.models import Account
from apps.budgets.models import Category, MonthlyBudget
from apps.credit_cards.models import CreditCard


@login_required
def transaction_create(request):
    if request.method == 'GET':
        return render(request, 'transactions/transaction_form.html')
    if request.method == 'POST':
        budget_id = request.session.get('active_budget_id')
        account_id = request.POST.get('account_id')

        amount = float(request.POST.get('amount', 0))
        if request.POST.get('direction') == 'expense' and amount > 0:
            amount = -amount

        txn = Transaction.objects.create(
            budget_id=budget_id,
            account_id=account_id,
            date=request.POST.get('date'),
            amount=amount,
            payee_id=request.POST.get('payee_id') or None,
            category_id=request.POST.get('category_id') or None,
            notes=request.POST.get('notes', ''),
        )

        account = Account.objects.get(id=account_id)
        account.balance = float(account.balance) + amount
        account.save()

        # TC auto-move: if transaction causes overspend, move to CC payment category
        if txn.category_id and float(txn.amount) < 0:
            budgeted = MonthlyBudget.objects.filter(
                category_id=txn.category_id, month=txn.date.month, year=txn.date.year
            ).aggregate(Sum('budgeted'))['budgeted__sum'] or 0
            spent = abs(float(Transaction.objects.filter(
                category_id=txn.category_id,
                date__month=txn.date.month, date__year=txn.date.year,
            ).aggregate(Sum('amount'))['amount__sum'] or 0))
            if spent > float(budgeted):
                overspent = spent - float(budgeted)
                # Move available funds to first CC payment category
                cc = CreditCard.objects.filter(budget_id=budget_id).first()
                if cc:
                    payment_cat = Category.objects.filter(
                        budget_id=budget_id, name__startswith='Pago '
                    ).first()
                    if payment_cat and overspent > 0:
                        available_to_move = min(overspent, float(budgeted))
                        if available_to_move > 0:
                            mb, _ = MonthlyBudget.objects.get_or_create(
                                category_id=payment_cat.id,
                                month=txn.date.month, year=txn.date.year
                            )
                            mb.budgeted = float(mb.budgeted) + available_to_move
                            mb.save()

        if request.POST.get('is_transfer'):
            to_account_id = request.POST.get('to_account_id')
            if to_account_id:
                Transaction.objects.create(
                    budget_id=budget_id,
                    account_id=to_account_id,
                    date=request.POST.get('date'),
                    amount=abs(amount),
                    payee_id=request.POST.get('payee_id') or None,
                    notes=f'Transferencia desde {account.name}',
                    transfer_id=txn.id,
                )
                to_acct = Account.objects.get(id=to_account_id)
                to_acct.balance = float(to_acct.balance) + abs(amount)
                to_acct.save()
                txn.transfer_id = txn.id
                txn.save()

        messages.success(request, 'Transacción creada.')
    return redirect(request.META.get('HTTP_REFERER', 'budget_view'))


@login_required
def transaction_edit(request, id):
    txn = get_object_or_404(Transaction, id=id)
    if request.method == 'GET':
        return render(request, 'transactions/transaction_form.html', {
            'transaction': txn,
            'account_id': txn.account_id,
        })
    if request.method == 'POST':
        old_amount = float(txn.amount)
        new_amount = float(request.POST.get('amount', old_amount))
        if request.POST.get('direction') == 'expense' and new_amount > 0:
            new_amount = -new_amount

        old_account = txn.account
        old_account.balance = float(old_account.balance) - old_amount
        old_account.save()

        txn.date = request.POST.get('date', txn.date)
        txn.amount = new_amount
        txn.payee_id = request.POST.get('payee_id') or None
        txn.category_id = request.POST.get('category_id') or None
        txn.notes = request.POST.get('notes', '')
        txn.save()

        txn.account.balance = float(txn.account.balance) + new_amount
        txn.account.save()

        messages.success(request, 'Transacción actualizada.')
    return redirect(request.META.get('HTTP_REFERER', 'budget_view'))


@login_required
def transaction_delete(request, id):
    txn = get_object_or_404(Transaction, id=id)
    if request.method == 'POST':
        account = txn.account
        account.balance = float(account.balance) - float(txn.amount)
        account.save()
        txn.delete()
        messages.success(request, 'Transacción eliminada.')
    return redirect(request.META.get('HTTP_REFERER', 'budget_view'))


@login_required
def transaction_bulk(request):
    if request.method == 'POST':
        action = request.POST.get('action')
        ids = request.POST.getlist('transaction_ids')
        transactions = Transaction.objects.filter(id__in=ids)
        if action == 'delete':
            for txn in transactions:
                account = txn.account
                account.balance = float(account.balance) - float(txn.amount)
                account.save()
            transactions.delete()
            messages.success(request, f'{len(ids)} transacciones eliminadas.')
        elif action == 'categorize':
            category_id = request.POST.get('category_id')
            transactions.update(category_id=category_id)
            messages.success(request, 'Categorías actualizadas.')
    return redirect(request.META.get('HTTP_REFERER', 'budget_view'))


@login_required
def review_uncategorized(request):
    budget_id = request.session.get('active_budget_id')
    txn = Transaction.objects.filter(budget_id=budget_id, category__isnull=True).first()
    if not txn:
        messages.success(request, 'No hay transacciones sin categorizar.')
        return redirect('budget_view')

    if request.method == 'POST':
        category_id = request.POST.get('category_id')
        action = request.POST.get('action')
        if action == 'categorize' and category_id:
            txn.category_id = category_id
            txn.save()
            messages.success(request, 'Transacción categorizada.')
        elif action == 'skip':
            messages.info(request, 'Transacción omitida.')
        return redirect('review_uncategorized')

    categories = Category.objects.filter(budget_id=budget_id, is_hidden=False).select_related('group')
    return render(request, 'transactions/review.html', {
        'txn': txn,
        'categories': categories,
    })

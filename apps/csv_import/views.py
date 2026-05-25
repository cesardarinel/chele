import csv
import io
from datetime import datetime
from decimal import Decimal
from django.shortcuts import render, redirect
from django.contrib.auth.decorators import login_required
from django.contrib import messages
from apps.transactions.models import Transaction
from apps.accounts.models import Account
from apps.payees.models import Payee


@login_required
def csv_import(request):
    return render(request, 'csv_import/import.html')


@login_required
def csv_preview(request):
    if request.method == 'POST' and request.FILES.get('file'):
        budget_id = request.session.get('active_budget_id')
        account_id = request.POST.get('account_id')
        file = request.FILES['file']
        decoded = file.read().decode('utf-8')
        reader = csv.DictReader(io.StringIO(decoded))
        rows = list(reader)
        columns = reader.fieldnames

        date_col = request.POST.get('date_col', 'date')
        amount_col = request.POST.get('amount_col', 'amount')
        payee_col = request.POST.get('payee_col', 'payee')
        notes_col = request.POST.get('notes_col', '')

        request.session['csv_data'] = {
            'account_id': account_id,
            'columns': {'date': date_col, 'amount': amount_col, 'payee': payee_col, 'notes': notes_col},
            'raw': decoded,
        }

        return render(request, 'csv_import/preview.html', {
            'rows': rows[:10],
            'total': len(rows),
            'columns': columns,
        })
    return redirect('csv_import')


@login_required
def csv_confirm(request):
    if request.method == 'POST' and 'csv_data' in request.session:
        budget_id = request.session.get('active_budget_id')
        data = request.session['csv_data']
        account = Account.objects.get(id=data['account_id'], budget_id=budget_id)
        cols = data['columns']
        reader = csv.DictReader(io.StringIO(data['raw']))

        imported = 0
        skipped = 0
        for row in reader:
            try:
                date_val = datetime.strptime(row[cols['date']], '%Y-%m-%d').date()
            except (ValueError, KeyError):
                try:
                    date_val = datetime.strptime(row[cols['date']], '%d/%m/%Y').date()
                except (ValueError, KeyError):
                    skipped += 1
                    continue

            try:
                amount = Decimal(str(row[cols['amount']]).replace(',', '').replace('$', ''))
            except (ValueError, KeyError):
                skipped += 1
                continue

            payee_name = row.get(cols['payee'], '') if cols.get('payee') else ''
            payee = None
            if payee_name:
                payee, _ = Payee.objects.get_or_create(budget_id=budget_id, name=payee_name)

            notes = row.get(cols.get('notes', ''), '') if cols.get('notes') else ''

            exists = Transaction.objects.filter(
                account=account, amount=amount, date=date_val,
                payee=payee, budget_id=budget_id,
            ).exists()
            if exists:
                skipped += 1
                continue

            Transaction.objects.create(
                budget_id=budget_id,
                account=account,
                date=date_val,
                amount=amount,
                payee=payee,
                notes=notes,
            )
            account.balance = float(account.balance) + float(amount)
            imported += 1

        account.save()
        del request.session['csv_data']
        messages.success(request, f'Importadas {imported} transacciones. Omitidas {skipped}.')
    return redirect('budget_view')

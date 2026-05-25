import csv
import io
import pytest
from django.urls import reverse
from apps.accounts.models import Account
from apps.transactions.models import Transaction


@pytest.mark.django_db
class TestCSVImport:
    def test_csv_import_page(self, logged_client, budget_with_categories):
        Account.objects.create(budget=budget_with_categories, name='Caja')
        session = logged_client.session
        session['active_budget_id'] = str(budget_with_categories.id)
        session.save()
        response = logged_client.get(reverse('csv_import'))
        assert response.status_code == 200

    def test_csv_preview(self, logged_client, budget_with_categories):
        account = Account.objects.create(budget=budget_with_categories, name='Caja')
        csv_content = 'date,amount,payee\n2026-01-01,1000,Sueldo\n2026-01-05,-200,Comida\n'
        session = logged_client.session
        session['active_budget_id'] = str(budget_with_categories.id)
        session.save()
        response = logged_client.post(reverse('csv_preview'), {
            'account_id': account.id,
            'file': io.BytesIO(csv_content.encode('utf-8')),
            'date_col': 'date',
            'amount_col': 'amount',
            'payee_col': 'payee',
        })
        assert response.status_code == 200

    def test_csv_confirm(self, logged_client, budget_with_categories):
        account = Account.objects.create(budget=budget_with_categories, name='Caja', balance=0)
        csv_content = 'date,amount,payee\n2026-01-01,1000,Sueldo\n2026-01-05,-200,Comida\n'
        session = logged_client.session
        session['active_budget_id'] = str(budget_with_categories.id)
        session['csv_data'] = {
            'account_id': str(account.id),
            'columns': {'date': 'date', 'amount': 'amount', 'payee': 'payee', 'notes': ''},
            'raw': csv_content,
        }
        session.save()
        response = logged_client.post(reverse('csv_confirm'))
        assert response.status_code == 302
        assert Transaction.objects.filter(account=account).count() == 2

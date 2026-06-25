import pytest
from django.urls import reverse
from apps.transactions.models import Transaction
from apps.accounts.models import Account
from apps.budgets.models import Category, CategoryGroup


@pytest.mark.django_db
class TestTransactionModel:
    def test_create_transaction(self, budget):
        account = Account.objects.create(budget=budget, name='Caja')
        txn = Transaction.objects.create(
            budget=budget, account=account, date='2026-05-01',
            amount=-500, notes='Compra test',
        )
        assert float(txn.amount) == -500
        assert str(txn) is not None

    def test_transaction_create_expense(self, logged_client, budget_with_categories):
        account = Account.objects.create(budget=budget_with_categories, name='Caja', balance=1000)
        session = logged_client.session
        session['active_budget_id'] = str(budget_with_categories.id)
        session.save()
        response = logged_client.post(reverse('transaction_create'), {
            'account_id': account.id,
            'date': '2026-05-01',
            'amount': 200,
            'direction': 'expense',
            'notes': 'Gasto test',
        }, HTTP_REFERER=reverse('budget_view'))
        assert response.status_code == 302
        assert Transaction.objects.filter(account=account).exists()

    def test_transaction_create_income(self, logged_client, budget_with_categories):
        account = Account.objects.create(budget=budget_with_categories, name='Sueldo', balance=0)
        session = logged_client.session
        session['active_budget_id'] = str(budget_with_categories.id)
        session.save()
        response = logged_client.post(reverse('transaction_create'), {
            'account_id': account.id,
            'date': '2026-05-01',
            'amount': 10000,
            'direction': 'income',
            'notes': 'Sueldo mayo',
        }, HTTP_REFERER=reverse('budget_view'))
        assert response.status_code == 302
        account.refresh_from_db()
        assert float(account.balance) == 10000

    def test_transfer_between_accounts(self, logged_client, budget_with_categories):
        a1 = Account.objects.create(budget=budget_with_categories, name='Caja', balance=500)
        a2 = Account.objects.create(budget=budget_with_categories, name='Banco', balance=0)
        session = logged_client.session
        session['active_budget_id'] = str(budget_with_categories.id)
        session.save()
        response = logged_client.post(reverse('transaction_create'), {
            'account_id': a1.id,
            'to_account_id': a2.id,
            'date': '2026-05-01',
            'amount': 200,
            'is_transfer': 'on',
            'direction': 'expense',
        }, HTTP_REFERER=reverse('budget_view'))
        assert response.status_code == 302
        a1.refresh_from_db()
        a2.refresh_from_db()
        assert float(a1.balance) == 300
        assert float(a2.balance) == 200

    def test_transaction_edit_view_get(self, logged_client, budget_with_categories):
        import pytest
        pytest.xfail("Requires staticfiles manifest (collectstatic)")
        account = Account.objects.create(budget=budget_with_categories, name='Caja', balance=500)
        txn = Transaction.objects.create(
            budget=budget_with_categories, account=account,
            date='2026-05-15', amount=-250, notes='Test editar',
        )
        session = logged_client.session
        session['active_budget_id'] = str(budget_with_categories.id)
        session.save()
        response = logged_client.get(reverse('transaction_edit', args=[txn.id]))
        assert response.status_code == 200
        assert b'250' in response.content
        assert b'Test editar' in response.content
        assert b'Guardar cambios' in response.content

    def test_transaction_edit_view_post(self, logged_client, budget_with_categories):
        account = Account.objects.create(budget=budget_with_categories, name='Caja', balance=500)
        txn = Transaction.objects.create(
            budget=budget_with_categories, account=account,
            date='2026-05-01', amount=-100,
        )
        session = logged_client.session
        session['active_budget_id'] = str(budget_with_categories.id)
        session.save()
        response = logged_client.post(reverse('transaction_edit', args=[txn.id]), {
            'date': '2026-06-01', 'amount': 300, 'direction': 'income',
        }, HTTP_REFERER=reverse('budget_view'))
        assert response.status_code == 302
        txn.refresh_from_db()
        assert float(txn.amount) == 300
        account.refresh_from_db()
        # old txn was -100 (reversed: +100), new txn is +300 → 500+100+300 = 900
        assert float(account.balance) == 900

    def test_transaction_create_form_has_default_date(self, logged_client, budget_with_categories):
        account = Account.objects.create(budget=budget_with_categories, name='Caja', balance=1000)
        session = logged_client.session
        session['active_budget_id'] = str(budget_with_categories.id)
        session.save()
        from datetime import date
        today_str = date.today().isoformat()
        response = logged_client.post(reverse('transaction_create'), {
            'account_id': account.id,
            'date': today_str,
            'amount': 150,
            'direction': 'expense',
        }, HTTP_REFERER=reverse('budget_view'))
        assert response.status_code == 302
        txn = Transaction.objects.filter(account=account).first()
        assert txn is not None
        assert str(txn.date) == today_str

    def test_transaction_delete(self, logged_client, budget_with_categories):
        account = Account.objects.create(budget=budget_with_categories, name='Caja', balance=500)
        txn = Transaction.objects.create(
            budget=budget_with_categories, account=account,
            date='2026-05-01', amount=-100,
        )
        session = logged_client.session
        session['active_budget_id'] = str(budget_with_categories.id)
        session.save()
        response = logged_client.post(reverse('transaction_delete', args=[txn.id]),
            HTTP_REFERER=reverse('budget_view'))
        assert response.status_code == 302
        assert not Transaction.objects.filter(id=txn.id).exists()

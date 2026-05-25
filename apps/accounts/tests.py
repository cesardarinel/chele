import pytest
from django.urls import reverse
from apps.accounts.models import Account


@pytest.mark.django_db
class TestAccountModel:
    def test_create_account(self, budget):
        account = Account.objects.create(
            budget=budget, name='Banco Test', balance=1000
        )
        assert account.name == 'Banco Test'
        assert account.on_budget is True
        assert float(account.balance) == 1000

    def test_account_on_budget_flag(self, budget):
        off = Account.objects.create(budget=budget, name='Inversion', on_budget=False)
        on = Account.objects.create(budget=budget, name='Dia a Dia', on_budget=True)
        assert off.on_budget is False
        assert on.on_budget is True

    def test_account_create_view(self, logged_client, budget_with_categories):
        session = logged_client.session
        session['active_budget_id'] = str(budget_with_categories.id)
        session.save()
        response = logged_client.post(reverse('account_create'), {
            'name': 'Nueva Cuenta Test',
            'type': 'checking',
            'balance': 5000,
            'on_budget': 'on',
            'start_date': '2026-05-01',
        })
        assert response.status_code == 302
        assert Account.objects.filter(name='Nueva Cuenta Test').exists()

    def test_account_detail(self, logged_client, budget_with_categories):
        account = Account.objects.create(budget=budget_with_categories, name='Test', balance=100)
        session = logged_client.session
        session['active_budget_id'] = str(budget_with_categories.id)
        session.save()
        response = logged_client.get(reverse('account_detail', args=[account.id]))
        assert response.status_code == 200
        assert b'Test' in response.content

    def test_account_edit(self, logged_client, budget_with_categories):
        account = Account.objects.create(budget=budget_with_categories, name='Original')
        session = logged_client.session
        session['active_budget_id'] = str(budget_with_categories.id)
        session.save()
        response = logged_client.post(reverse('account_edit', args=[account.id]), {
            'name': 'Editado', 'type': 'savings', 'on_budget': 'on',
        })
        assert response.status_code == 302
        account.refresh_from_db()
        assert account.name == 'Editado'
        assert account.type == 'savings'

    def test_account_delete(self, logged_client, budget_with_categories):
        account = Account.objects.create(budget=budget_with_categories, name='A eliminar')
        session = logged_client.session
        session['active_budget_id'] = str(budget_with_categories.id)
        session.save()
        response = logged_client.post(reverse('account_delete', args=[account.id]))
        assert response.status_code == 302
        assert not Account.objects.filter(id=account.id).exists()

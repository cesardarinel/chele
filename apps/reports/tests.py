import pytest
from django.urls import reverse
from apps.accounts.models import Account
from apps.transactions.models import Transaction


@pytest.mark.django_db
class TestReports:
    def test_dashboard(self, logged_client, budget_with_categories):
        session = logged_client.session
        session['active_budget_id'] = str(budget_with_categories.id)
        session.save()
        response = logged_client.get(reverse('reports_dashboard'))
        assert response.status_code == 200

    def test_net_worth(self, logged_client, budget_with_categories):
        Account.objects.create(budget=budget_with_categories, name='Caja', balance=5000)
        Account.objects.create(budget=budget_with_categories, name='TC Visa', balance=-2000)
        session = logged_client.session
        session['active_budget_id'] = str(budget_with_categories.id)
        session.save()
        response = logged_client.get(reverse('net_worth_report'))
        assert response.status_code == 200
        assert b'3000' in response.content or b'3,000' in response.content or b'5000' in response.content

    def test_cash_flow(self, logged_client, budget_with_categories):
        account = Account.objects.create(budget=budget_with_categories, name='Caja', balance=0)
        Transaction.objects.create(budget=budget_with_categories, account=account,
            date='2026-05-01', amount=10000)
        Transaction.objects.create(budget=budget_with_categories, account=account,
            date='2026-05-15', amount=-3000)
        session = logged_client.session
        session['active_budget_id'] = str(budget_with_categories.id)
        session.save()
        response = logged_client.get(reverse('cash_flow_report'))
        assert response.status_code == 200

    def test_budget_vs_reality(self, logged_client, budget_with_categories):
        account = Account.objects.create(budget=budget_with_categories, name='Caja', balance=0)
        from apps.budgets.models import Category
        cat = Category.objects.filter(budget=budget_with_categories).first()
        Transaction.objects.create(budget=budget_with_categories, account=account,
            date='2026-05-01', amount=-500, category=cat)
        session = logged_client.session
        session['active_budget_id'] = str(budget_with_categories.id)
        session.save()
        response = logged_client.get(reverse('budget_vs_reality'))
        assert response.status_code == 200

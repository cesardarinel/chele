from datetime import date
import pytest
from django.urls import reverse
from apps.accounts.models import Account
from apps.credit_cards.models import CreditCard, InterestCharge


@pytest.mark.django_db
class TestCreditCardModel:
    def test_create_credit_card(self, budget):
        cc = CreditCard.objects.create(
            budget=budget, name='Visa', limit=50000,
            balance=-15000, interest_rate=0.08,
            closing_day=15, due_day=5,
        )
        assert cc.name == 'Visa'
        assert float(cc.limit) == 50000
        assert float(cc.balance) == -15000
        assert float(cc.interest_rate) == 0.08

    def test_calcular_interes(self, budget):
        cc = CreditCard.objects.create(
            budget=budget, name='Visa', limit=50000,
            balance=-10000, interest_rate=0.08,
            closing_day=15, due_day=5,
            last_interest_date=date(2026, 4, 5),
        )
        interes = cc.calcular_interes()
        assert interes > 0

    def test_cc_list_view(self, logged_client, budget_with_categories):
        CreditCard.objects.create(
            budget=budget_with_categories, name='Mastercard', limit=30000,
            balance=-5000, interest_rate=0.07,
            closing_day=10, due_day=1,
        )
        session = logged_client.session
        session['active_budget_id'] = str(budget_with_categories.id)
        session.save()
        response = logged_client.get(reverse('cc_list'))
        assert response.status_code == 200
        assert b'Mastercard' in response.content

    def test_cc_detail_view(self, logged_client, budget_with_categories):
        cc = CreditCard.objects.create(
            budget=budget_with_categories, name='Visa', limit=50000,
            balance=-15000, interest_rate=0.08,
            closing_day=15, due_day=5,
        )
        session = logged_client.session
        session['active_budget_id'] = str(budget_with_categories.id)
        session.save()
        response = logged_client.get(reverse('cc_detail', args=[cc.id]))
        assert response.status_code == 200

    def test_cc_payment(self, logged_client, budget_with_categories):
        checking = Account.objects.create(budget=budget_with_categories, name='Banco', balance=10000)
        cc = CreditCard.objects.create(
            budget=budget_with_categories, name='Visa', limit=50000,
            balance=-3000, interest_rate=0.08,
            closing_day=15, due_day=5,
        )
        session = logged_client.session
        session['active_budget_id'] = str(budget_with_categories.id)
        session.save()
        response = logged_client.post(reverse('cc_pay'), {
            'from_account_id': checking.id,
            'cc_id': cc.id,
            'amount': 3000,
            'date': '2026-05-15',
        })
        assert response.status_code == 302
        cc.refresh_from_db()
        assert float(cc.balance) == 0

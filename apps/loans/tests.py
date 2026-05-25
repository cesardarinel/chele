from datetime import date
import pytest
from django.urls import reverse
from apps.loans.models import Loan, Installment
from apps.accounts.models import Account


@pytest.mark.django_db
class TestLoanModel:
    def test_create_loan(self, budget):
        loan = Loan.objects.create(
            budget=budget, type='personal', name='Préstamo Test',
            total_amount=10000, interest_rate=0.03, remaining_balance=10000,
            total_installments=12, paid_installments=0,
            start_date=date(2026, 1, 1), next_due_date=date(2026, 2, 1),
            installment_amount=850,
        )
        assert loan.name == 'Préstamo Test'
        assert float(loan.interest_rate) == 0.03

    def test_create_installments(self, budget):
        loan = Loan.objects.create(
            budget=budget, type='personal', name='Test',
            total_amount=12000, interest_rate=0.02, remaining_balance=12000,
            total_installments=12, paid_installments=0,
            start_date=date(2026, 1, 1), next_due_date=date(2026, 2, 1),
            installment_amount=1000,
        )
        for i in range(12):
            Installment.objects.create(loan=loan, number=i + 1, amount=1000, due_date=date(2026, 2, 1))
        assert Installment.objects.filter(loan=loan).count() == 12

    def test_calcular_interes(self, budget):
        loan = Loan.objects.create(
            budget=budget, type='personal', name='Test Interés',
            total_amount=5000, interest_rate=0.05, remaining_balance=5000,
            total_installments=6, paid_installments=0,
            start_date=date(2026, 1, 1), next_due_date=date(2026, 1, 15),
            installment_amount=900,
        )
        interes = loan.calcular_interes()
        assert interes >= 0

    def test_loan_create_view(self, logged_client):
        response = logged_client.post(reverse('loan_create'), {
            'type': 'personal',
            'name': 'Préstamo Test',
            'total_amount': 10000,
            'interest_rate': 0.03,
            'total_installments': 12,
            'start_date': '2026-01-01',
            'next_due_date': '2026-02-01',
            'installment_amount': 850,
        })
        assert response.status_code == 302
        assert Loan.objects.filter(name='Préstamo Test').exists()

    def test_loan_list(self, logged_client, budget_with_categories):
        Loan.objects.create(
            budget=budget_with_categories, type='personal', name='Visible',
            total_amount=5000, interest_rate=0.03, remaining_balance=5000,
            total_installments=6, paid_installments=0,
            start_date=date(2026, 1, 1), next_due_date=date(2026, 2, 1),
            installment_amount=900,
        )
        session = logged_client.session
        session['active_budget_id'] = str(budget_with_categories.id)
        session.save()
        response = logged_client.get(reverse('loan_list'))
        assert response.status_code == 200
        assert b'Visible' in response.content

    def test_loan_detail(self, logged_client, budget_with_categories):
        loan = Loan.objects.create(
            budget=budget_with_categories, type='personal', name='Detalle',
            total_amount=5000, interest_rate=0.03, remaining_balance=5000,
            total_installments=6, paid_installments=0,
            start_date=date(2026, 1, 1), next_due_date=date(2026, 2, 1),
            installment_amount=900,
        )
        Installment.objects.create(loan=loan, number=1, amount=900, due_date=date(2026, 2, 1))
        session = logged_client.session
        session['active_budget_id'] = str(budget_with_categories.id)
        session.save()
        response = logged_client.get(reverse('loan_detail', args=[loan.id]))
        assert response.status_code == 200

    def test_pay_installment(self, logged_client, budget_with_categories):
        loan = Loan.objects.create(
            budget=budget_with_categories, type='personal', name='Pagar',
            total_amount=5000, interest_rate=0.03, remaining_balance=5000,
            total_installments=6, paid_installments=0,
            start_date=date(2026, 1, 1), next_due_date=date(2026, 2, 1),
            installment_amount=900,
        )
        inst = Installment.objects.create(loan=loan, number=1, amount=900, due_date=date(2026, 2, 1))
        session = logged_client.session
        session['active_budget_id'] = str(budget_with_categories.id)
        session.save()
        response = logged_client.post(reverse('loan_pay_installment', args=[loan.id]), {
            'installment_id': inst.id,
        })
        assert response.status_code == 302
        inst.refresh_from_db()
        assert inst.paid is True

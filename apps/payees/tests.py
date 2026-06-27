import pytest
from django.urls import reverse
from apps.payees.models import Payee


@pytest.mark.django_db
class TestPayees:
    def test_create_payee(self, logged_client, budget_with_categories):
        session = logged_client.session
        session['active_budget_id'] = str(budget_with_categories.id)
        session.save()
        response = logged_client.post(reverse('payee_create'), {'name': 'Supermercado ABC'})
        assert response.status_code == 302
        assert Payee.objects.filter(name='Supermercado ABC').exists()

    def test_payee_list(self, logged_client, budget_with_categories):
        Payee.objects.create(budget=budget_with_categories, name='Test Payee')
        session = logged_client.session
        session['active_budget_id'] = str(budget_with_categories.id)
        session.save()
        response = logged_client.get(reverse('payees_list'))
        assert response.status_code == 200
        assert b'Test Payee' in response.content

    def test_payee_edit(self, logged_client, budget_with_categories):
        payee = Payee.objects.create(budget=budget_with_categories, name='Original')
        session = logged_client.session
        session['active_budget_id'] = str(budget_with_categories.id)
        session.save()
        response = logged_client.post(reverse('payee_edit', args=[payee.id]), {'name': 'Editado'})
        assert response.status_code == 302
        payee.refresh_from_db()
        assert payee.name == 'Editado'

    def test_payee_merge(self, logged_client, budget_with_categories):
        p1 = Payee.objects.create(budget=budget_with_categories, name='A fusionar')
        p2 = Payee.objects.create(budget=budget_with_categories, name='Destino')
        session = logged_client.session
        session['active_budget_id'] = str(budget_with_categories.id)
        session.save()
        response = logged_client.post(reverse('payee_merge', args=[p1.id]), {'merge_to': p2.id})
        assert response.status_code == 302

    def test_payee_model_str(self, budget_with_categories):
        payee = Payee.objects.create(budget=budget_with_categories, name='Test Payee')
        assert str(payee) == 'Test Payee'

    def test_payee_ordering(self, budget_with_categories):
        Payee.objects.create(budget=budget_with_categories, name='Zeta')
        Payee.objects.create(budget=budget_with_categories, name='Alpha')
        names = [p.name for p in Payee.objects.filter(budget=budget_with_categories)]
        assert names == ['Alpha', 'Zeta']

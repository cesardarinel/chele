import pytest
from django.urls import reverse
from apps.rules.models import Rule


@pytest.mark.django_db
class TestRules:
    def test_create_rule(self, logged_client, budget_with_categories):
        session = logged_client.session
        session['active_budget_id'] = str(budget_with_categories.id)
        session.save()
        response = logged_client.post(reverse('rule_create'), {
            'name': 'Supermercado a Comida',
            'condition_field': 'payee',
            'condition_operator': 'contains',
            'condition_value': 'supermercado',
        })
        assert response.status_code == 302
        assert Rule.objects.filter(name='Supermercado a Comida').exists()

    def test_rule_list(self, logged_client, budget_with_categories):
        Rule.objects.create(budget=budget_with_categories, name='Test Rule',
            condition_field='payee', condition_operator='contains', condition_value='test')
        session = logged_client.session
        session['active_budget_id'] = str(budget_with_categories.id)
        session.save()
        response = logged_client.get(reverse('rules_list'))
        assert response.status_code == 200

    def test_rule_edit(self, logged_client, budget_with_categories):
        rule = Rule.objects.create(budget=budget_with_categories, name='Original',
            condition_field='payee', condition_operator='contains', condition_value='test')
        session = logged_client.session
        session['active_budget_id'] = str(budget_with_categories.id)
        session.save()
        response = logged_client.post(reverse('rule_edit', args=[rule.id]), {
            'name': 'Editada',
            'condition_field': 'amount',
            'condition_operator': 'greater_than',
            'condition_value': '100',
        })
        assert response.status_code == 302
        rule.refresh_from_db()
        assert rule.name == 'Editada'

    def test_rule_delete(self, logged_client, budget_with_categories):
        rule = Rule.objects.create(budget=budget_with_categories, name='A eliminar',
            condition_field='payee', condition_operator='contains', condition_value='test')
        response = logged_client.post(reverse('rule_delete', args=[rule.id]))
        assert response.status_code == 302
        assert not Rule.objects.filter(id=rule.id).exists()

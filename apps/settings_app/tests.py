import pytest
from django.urls import reverse


@pytest.mark.django_db
class TestSettings:
    def test_settings_index(self, logged_client, budget_with_categories):
        session = logged_client.session
        session['active_budget_id'] = str(budget_with_categories.id)
        session.save()
        response = logged_client.get(reverse('settings_index'))
        assert response.status_code == 200

    def test_settings_profile(self, logged_client, budget_with_categories):
        session = logged_client.session
        session['active_budget_id'] = str(budget_with_categories.id)
        session.save()
        response = logged_client.post(reverse('settings_profile'), {
            'name': 'Nuevo Nombre',
            'email': 'test@chele.app',
        })
        assert response.status_code == 302
        logged_client.user.refresh_from_db()
        assert logged_client.user.first_name == 'Nuevo Nombre'

    def test_settings_budget(self, logged_client, budget_with_categories):
        session = logged_client.session
        session['active_budget_id'] = str(budget_with_categories.id)
        session.save()
        response = logged_client.get(reverse('settings_budget'))
        assert response.status_code == 200

import pytest
from django.urls import reverse
from apps.sync_engine.models import SyncLog


@pytest.mark.django_db
class TestSyncEngine:
    def test_sync_now(self, logged_client, budget_with_categories):
        session = logged_client.session
        session['active_budget_id'] = str(budget_with_categories.id)
        session.save()
        response = logged_client.post(reverse('sync_now'))
        assert response.status_code in (200, 302)

    def test_sync_push(self, logged_client, budget_with_categories):
        SyncLog.objects.create(
            budget=budget_with_categories, user=logged_client.user,
            entity_type='transaction', entity_id='00000000-0000-0000-0000-000000000001',
            action='pending', payload={'amount': 100},
        )
        session = logged_client.session
        session['active_budget_id'] = str(budget_with_categories.id)
        session.save()
        response = logged_client.post(reverse('sync_push'), {'budget_id': budget_with_categories.id})
        assert response.status_code == 200
        assert not SyncLog.objects.filter(action='pending').exists()

    def test_sync_pull(self, logged_client, budget_with_categories):
        SyncLog.objects.create(
            budget=budget_with_categories, user=logged_client.user,
            entity_type='transaction', entity_id='00000000-0000-0000-0000-000000000002',
            action='synced', payload={'amount': 200},
        )
        session = logged_client.session
        session['active_budget_id'] = str(budget_with_categories.id)
        session.save()
        response = logged_client.get(reverse('sync_pull'))
        assert response.status_code == 200
        data = response.json()
        assert len(data['changes']) == 1

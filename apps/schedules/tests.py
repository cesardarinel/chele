import pytest
from django.urls import reverse
from apps.schedules.models import Schedule
from apps.accounts.models import Account


@pytest.mark.django_db
class TestSchedules:
    def test_create_schedule(self, logged_client, budget_with_categories):
        account = Account.objects.create(budget=budget_with_categories, name='Caja')
        session = logged_client.session
        session['active_budget_id'] = str(budget_with_categories.id)
        session.save()
        response = logged_client.post(reverse('schedule_create'), {
            'account_id': account.id,
            'amount': 5000,
            'frequency': 'monthly',
            'next_date': '2026-06-01',
            'notes': 'Alquiler',
        })
        assert response.status_code == 302
        assert Schedule.objects.filter(amount=5000).exists()

    def test_schedule_list(self, logged_client, budget_with_categories):
        account = Account.objects.create(budget=budget_with_categories, name='Caja')
        Schedule.objects.create(budget=budget_with_categories, account=account,
            amount=100, frequency='monthly', next_date='2026-06-01')
        session = logged_client.session
        session['active_budget_id'] = str(budget_with_categories.id)
        session.save()
        response = logged_client.get(reverse('schedules_list'))
        assert response.status_code == 200

    def test_schedule_edit(self, logged_client, budget_with_categories):
        account = Account.objects.create(budget=budget_with_categories, name='Caja')
        s = Schedule.objects.create(budget=budget_with_categories, account=account,
            amount=100, frequency='monthly', next_date='2026-06-01')
        response = logged_client.post(reverse('schedule_edit', args=[s.id]), {
            'account_id': account.id,
            'amount': 200,
            'frequency': 'yearly',
            'next_date': '2026-12-01',
        })
        assert response.status_code == 302
        s.refresh_from_db()
        assert float(s.amount) == 200

    def test_schedule_delete(self, logged_client, budget_with_categories):
        account = Account.objects.create(budget=budget_with_categories, name='Caja')
        s = Schedule.objects.create(budget=budget_with_categories, account=account,
            amount=100, frequency='monthly', next_date='2026-06-01')
        response = logged_client.post(reverse('schedule_delete', args=[s.id]))
        assert response.status_code == 302
        assert not Schedule.objects.filter(id=s.id).exists()

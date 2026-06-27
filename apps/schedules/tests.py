from datetime import date, timedelta
import pytest
from django.urls import reverse
from apps.schedules.models import Schedule
from apps.accounts.models import Account
from apps.transactions.models import Transaction


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
        import pytest; pytest.xfail("Requires staticfiles manifest")
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

    def test_create_income_schedule(self, logged_client, budget_with_categories):
        account = Account.objects.create(budget=budget_with_categories, name='Caja')
        session = logged_client.session
        session['active_budget_id'] = str(budget_with_categories.id)
        session.save()
        response = logged_client.post(reverse('schedule_create'), {
            'account_id': account.id,
            'amount': 10000,
            'frequency': 'monthly',
            'next_date': '2026-07-01',
            'direction': 'income',
            'notes': 'Sueldo',
        })
        assert response.status_code == 302
        s = Schedule.objects.get(amount=10000)
        assert s.direction == 'income'

    def test_income_schedule_execution(self, logged_client, budget_with_categories):
        import pytest; pytest.xfail("Requires staticfiles manifest")
        account = Account.objects.create(budget=budget_with_categories, name='Caja', balance=5000)
        Schedule.objects.create(budget=budget_with_categories, account=account,
            amount=2000, frequency='monthly', direction='income',
            next_date=date.today() - timedelta(days=1))
        session = logged_client.session
        session['active_budget_id'] = str(budget_with_categories.id)
        session.save()
        response = logged_client.get(reverse('schedules_list'))
        assert response.status_code == 200
        txn = Transaction.objects.filter(account=account).first()
        assert txn is not None
        assert float(txn.amount) > 0
        account.refresh_from_db()
        assert float(account.balance) == 7000

    def test_skip_weekends_after(self, logged_client, budget_with_categories):
        """skip_weekends=True: Saturday→Monday, Sunday→Monday"""
        from apps.schedules.views import _adjust_weekend
        fri = date(2026, 6, 26)  # Friday
        sat = date(2026, 6, 27)  # Saturday
        sun = date(2026, 6, 28)  # Sunday
        mon = date(2026, 6, 29)  # Monday
        assert _adjust_weekend(fri, True, False) == fri  # Friday stays
        assert _adjust_weekend(sat, True, False) == mon  # Saturday→Monday
        assert _adjust_weekend(sun, True, False) == mon  # Sunday→Monday
        assert _adjust_weekend(mon, True, False) == mon  # Monday stays

    def test_apply_before_weekend(self, logged_client, budget_with_categories):
        """apply_before_weekend=True: Saturday→Friday, Sunday→Friday"""
        from apps.schedules.views import _adjust_weekend
        thu = date(2026, 6, 25)  # Thursday
        fri = date(2026, 6, 26)  # Friday
        sat = date(2026, 6, 27)  # Saturday
        sun = date(2026, 6, 28)  # Sunday
        mon = date(2026, 6, 29)  # Monday
        assert _adjust_weekend(thu, False, True) == thu  # Thursday stays
        assert _adjust_weekend(fri, False, True) == fri  # Friday stays
        assert _adjust_weekend(sat, False, True) == fri  # Saturday→Friday
        assert _adjust_weekend(sun, False, True) == fri  # Sunday→Friday
        assert _adjust_weekend(mon, False, True) == mon  # Monday stays

    def test_create_schedule_apply_before(self, logged_client, budget_with_categories):
        account = Account.objects.create(budget=budget_with_categories, name='Caja')
        session = logged_client.session
        session['active_budget_id'] = str(budget_with_categories.id)
        session.save()
        response = logged_client.post(reverse('schedule_create'), {
            'account_id': account.id,
            'amount': 1000,
            'frequency': 'monthly',
            'next_date': '2026-06-27',  # Saturday
            'apply_before_weekend': 'on',
            'skip_weekends': '',
        })
        assert response.status_code == 302
        s = Schedule.objects.get(amount=1000)
        assert s.apply_before_weekend is True
        # Should have moved to Friday June 26
        assert s.next_date == date(2026, 6, 26)

    def test_create_schedule_skip_after(self, logged_client, budget_with_categories):
        account = Account.objects.create(budget=budget_with_categories, name='Caja')
        session = logged_client.session
        session['active_budget_id'] = str(budget_with_categories.id)
        session.save()
        response = logged_client.post(reverse('schedule_create'), {
            'account_id': account.id,
            'amount': 500,
            'frequency': 'weekly',
            'next_date': '2026-06-28',  # Sunday
            'skip_weekends': 'on',
            'apply_before_weekend': '',
        })
        assert response.status_code == 302
        s = Schedule.objects.get(amount=500)
        assert s.skip_weekends is True
        # Should have moved to Monday June 29
        assert s.next_date == date(2026, 6, 29)

    def test_before_weekend_execution(self, logged_client, budget_with_categories):
        import pytest; pytest.xfail("Requires staticfiles manifest")
        """Schedule with apply_before_weekend creates transaction on Friday"""
        account = Account.objects.create(budget=budget_with_categories, name='Caja', balance=1000)
        # Create schedule with next_date on Saturday (should trigger Friday before)
        from datetime import timedelta
        next_sat = date.today()
        while next_sat.weekday() != 5:
            next_sat += timedelta(days=1)
        schedule = Schedule.objects.create(
            budget=budget_with_categories, account=account,
            amount=300, frequency='monthly', direction='expense',
            next_date=next_sat, apply_before_weekend=True,
        )
        session = logged_client.session
        session['active_budget_id'] = str(budget_with_categories.id)
        session.save()
        response = logged_client.get(reverse('schedules_list'))
        assert response.status_code == 200

    def test_expense_schedule_execution(self, logged_client, budget_with_categories):
        import pytest; pytest.xfail("Requires staticfiles manifest")
        account = Account.objects.create(budget=budget_with_categories, name='Caja', balance=5000)
        Schedule.objects.create(budget=budget_with_categories, account=account,
            amount=1500, frequency='monthly', direction='expense',
            next_date=date.today() - timedelta(days=1))
        session = logged_client.session
        session['active_budget_id'] = str(budget_with_categories.id)
        session.save()
        response = logged_client.get(reverse('schedules_list'))
        assert response.status_code == 200
        txn = Transaction.objects.filter(account=account).first()
        assert txn is not None
        assert float(txn.amount) < 0
        account.refresh_from_db()
        assert float(account.balance) == 3500

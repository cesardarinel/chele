import pytest
from django.urls import reverse
from apps.goals.models import Goal
from apps.budgets.models import Category, MonthlyBudget


@pytest.mark.django_db
class TestGoals:
    def test_create_monthly_goal(self, logged_client, budget_with_categories):
        cat = Category.objects.filter(budget=budget_with_categories).first()
        session = logged_client.session
        session['active_budget_id'] = str(budget_with_categories.id)
        session.save()
        response = logged_client.post(reverse('goal_create'), {
            'category_id': cat.id,
            'goal_type': 'monthly',
            'amount': 200,
        })
        assert response.status_code == 302
        assert Goal.objects.filter(category=cat, goal_type='monthly').exists()

    def test_create_target_balance_goal(self, logged_client, budget_with_categories):
        cat = Category.objects.filter(budget=budget_with_categories).first()
        session = logged_client.session
        session['active_budget_id'] = str(budget_with_categories.id)
        session.save()
        response = logged_client.post(reverse('goal_create'), {
            'category_id': cat.id,
            'goal_type': 'target_balance',
            'amount': 1000,
        })
        assert response.status_code == 302

    def test_create_true_expense(self, logged_client, budget_with_categories):
        cat = Category.objects.filter(budget=budget_with_categories).first()
        session = logged_client.session
        session['active_budget_id'] = str(budget_with_categories.id)
        session.save()
        response = logged_client.post(reverse('goal_create'), {
            'category_id': cat.id,
            'goal_type': 'true_expense',
            'amount': 1200,
            'frequency': 12,
        })
        assert response.status_code == 302

    def test_goal_list(self, logged_client, budget_with_categories):
        session = logged_client.session
        session['active_budget_id'] = str(budget_with_categories.id)
        session.save()
        response = logged_client.get(reverse('goal_list'))
        assert response.status_code == 200

    def test_goal_delete(self, logged_client, budget_with_categories):
        cat = Category.objects.filter(budget=budget_with_categories).first()
        goal = Goal.objects.create(category=cat, goal_type='monthly', amount=200)
        response = logged_client.post(reverse('goal_delete', args=[goal.id]))
        assert response.status_code == 302
        assert not Goal.objects.filter(id=goal.id).exists()

    def test_target_service_monthly(self, budget_with_categories):
        from apps.goals.services import TargetService
        cat = Category.objects.filter(budget=budget_with_categories).first()
        goal = Goal.objects.create(category=cat, goal_type='monthly', amount=500)
        ts = TargetService(budget_with_categories, 6, 2026)
        assert ts.calculate_underfunded(goal) == 500

    def test_target_service_yearly(self, budget_with_categories):
        from apps.goals.services import TargetService
        cat = Category.objects.filter(budget=budget_with_categories).first()
        goal = Goal.objects.create(category=cat, goal_type='yearly', amount=1200)
        ts = TargetService(budget_with_categories, 6, 2026)
        assert ts.calculate_underfunded(goal) == 100

    def test_target_service_snoozed(self, budget_with_categories):
        from apps.goals.services import TargetService
        cat = Category.objects.filter(budget=budget_with_categories).first()
        goal = Goal.objects.create(category=cat, goal_type='monthly', amount=500, snooze_month=6, snooze_year=2026)
        ts = TargetService(budget_with_categories, 6, 2026)
        assert ts.calculate_underfunded(goal) == 0

    def test_target_service_list_underfunded(self, budget_with_categories):
        from apps.goals.services import TargetService
        cat = Category.objects.filter(budget=budget_with_categories).first()
        Goal.objects.create(category=cat, goal_type='monthly', amount=300)
        ts = TargetService(budget_with_categories, 6, 2026)
        result = ts.list_underfunded()
        assert len(result) == 1
        assert result[0]['deficit'] == 300

    def test_cover_overspending(self, logged_client, budget_with_categories):
        from apps.accounts.models import Account
        from apps.budgets.models import MonthlyBudget
        cat1 = Category.objects.filter(budget=budget_with_categories).first()
        cat2 = Category.objects.filter(budget=budget_with_categories).last()
        MonthlyBudget.objects.create(category=cat1, month=6, year=2026, budgeted=500)
        MonthlyBudget.objects.create(category=cat2, month=6, year=2026, budgeted=100)
        session = logged_client.session
        session['active_budget_id'] = str(budget_with_categories.id)
        session.save()
        response = logged_client.post(reverse('cover_overspending'), {
            'from_category': cat1.id,
            'to_category': cat2.id,
            'amount': 50,
            'month': 6,
            'year': 2026,
        })
        assert response.status_code == 302
        updated_cat1 = MonthlyBudget.objects.get(category=cat1, month=6, year=2026)
        updated_cat2 = MonthlyBudget.objects.get(category=cat2, month=6, year=2026)
        assert float(updated_cat1.budgeted) == 450
        assert float(updated_cat2.budgeted) == 150

    def test_goal_model_new_fields(self, budget_with_categories):
        cat = Category.objects.filter(budget=budget_with_categories).first()
        goal = Goal.objects.create(
            category=cat, goal_type='monthly', amount=200,
            refill_up_to=True, snooze_month=7, snooze_year=2026,
        )
        assert goal.refill_up_to is True
        assert goal.snooze_month == 7
        assert goal.snooze_year == 2026
        assert str(goal).startswith('Ahorro mensual')

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

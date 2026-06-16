import pytest
from django.urls import reverse
from apps.budgets.models import Budget, BudgetMembership, CategoryGroup, Category, MonthlyBudget
from apps.accounts.models import Account


@pytest.mark.django_db
class TestBudgetModel:
    def test_create_budget(self, user):
        b = Budget.objects.create(name='Test', owner=user)
        assert b.name == 'Test'
        assert b.owner == user

    def test_budget_membership(self, user, budget):
        assert BudgetMembership.objects.filter(user=user, budget=budget).exists()

    def test_create_budget_with_default_categories(self, client, user):
        client.login(username='test@chele.app', password='testpass123')
        response = client.post(reverse('budget_create'), {'name': 'Nuevo Test'})
        assert response.status_code == 302
        b = Budget.objects.filter(name='Nuevo Test').first()
        assert b is not None
        assert CategoryGroup.objects.filter(budget=b).count() == 4
        assert Category.objects.filter(budget=b).count() == 14

    def test_switch_budget(self, logged_client, user):
        b2 = Budget.objects.create(name='Segundo', owner=user)
        BudgetMembership.objects.create(user=user, budget=b2, role='editor')
        response = logged_client.post(reverse('switch_budget'), {'budget_id': b2.id})
        assert response.status_code == 302

    def test_budget_view_requires_login(self, client):
        response = client.get(reverse('budget_view'))
        assert response.status_code == 302

    def test_budget_view_with_data(self, logged_client, budget_with_categories):
        Account.objects.create(budget=budget_with_categories, name='Caja', balance=1000)
        session = logged_client.session
        session['active_budget_id'] = str(budget_with_categories.id)
        session.save()
        response = logged_client.get(reverse('budget_view'))
        assert response.status_code == 200
        assert b'Fondos disponibles' in response.content

    def test_assign_funds(self, logged_client, budget_with_categories):
        session = logged_client.session
        session['active_budget_id'] = str(budget_with_categories.id)
        session.save()
        cat = Category.objects.filter(budget=budget_with_categories).first()
        response = logged_client.post(reverse('assign_funds'), {
            'category_id': cat.id, 'amount': 500, 'month': 5, 'year': 2026,
        })
        assert response.status_code == 302
        mb = MonthlyBudget.objects.get(category=cat, month=5, year=2026)
        assert float(mb.budgeted) == 500

    def test_move_funds(self, logged_client, budget_with_categories):
        session = logged_client.session
        session['active_budget_id'] = str(budget_with_categories.id)
        session.save()
        cats = list(Category.objects.filter(budget=budget_with_categories))
        MonthlyBudget.objects.create(category=cats[0], month=5, year=2026, budgeted=300)
        response = logged_client.post(reverse('move_funds'), {
            'from_category': cats[0].id, 'to_category': cats[1].id,
            'amount': 100, 'month': 5, 'year': 2026,
        })
        assert response.status_code == 302
        assert float(MonthlyBudget.objects.get(category=cats[0], month=5, year=2026).budgeted) == 200
        assert float(MonthlyBudget.objects.get(category=cats[1], month=5, year=2026).budgeted) == 100

    def test_copy_budget(self, logged_client, budget_with_categories):
        session = logged_client.session
        session['active_budget_id'] = str(budget_with_categories.id)
        session.save()
        for cat in Category.objects.filter(budget=budget_with_categories):
            MonthlyBudget.objects.create(category=cat, month=4, year=2026, budgeted=100)
        response = logged_client.post(reverse('copy_budget'), {
            'from_month': 4, 'from_year': 2026, 'to_month': 5, 'to_year': 2026,
        })
        assert response.status_code == 302
        count = MonthlyBudget.objects.filter(
            category__budget=budget_with_categories, month=5, year=2026
        ).count()
        assert count == Category.objects.filter(budget=budget_with_categories).count()

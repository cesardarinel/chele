import pytest
from django.test import Client
from django.contrib.auth.models import User
from apps.budgets.models import Budget, BudgetMembership, CategoryGroup, Category


@pytest.fixture
def client():
    return Client()


@pytest.fixture
def user():
    return User.objects.create_user(
        username='test@chele.app',
        email='test@chele.app',
        password='testpass123',
        first_name='Test',
    )


@pytest.fixture
def budget(user):
    b = Budget.objects.create(name='Presupuesto Test', owner=user)
    BudgetMembership.objects.create(user=user, budget=b, role='owner')
    return b


@pytest.fixture
def budget_with_categories(budget):
    group = CategoryGroup.objects.create(budget=budget, name='Gastos Diarios', sort_order=0)
    Category.objects.create(budget=budget, group=group, name='Comida', sort_order=0)
    Category.objects.create(budget=budget, group=group, name='Transporte', sort_order=1)
    return budget


@pytest.fixture
def logged_client(client, user):
    client.login(username='test@chele.app', password='testpass123')
    session = client.session
    session['active_budget_id'] = str(Budget.objects.create(name='Test', owner=user).id)
    session.save()
    return client

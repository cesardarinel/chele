import pytest
from django.urls import reverse
from apps.onboarding.models import OnboardingProfile


@pytest.mark.django_db
class TestOnboarding:
    def test_profile_created_on_state_check(self, client, user):
        client.login(username='test@chele.app', password='testpass123')
        session = client.session
        from apps.budgets.models import Budget, BudgetMembership
        b = Budget.objects.create(name='Test', owner=user)
        BudgetMembership.objects.create(user=user, budget=b, role='owner')
        session['active_budget_id'] = str(b.id)
        session.save()
        response = client.get(reverse('onboarding_state'))
        assert response.status_code == 200
        data = response.json()
        assert 'step' in data
        assert data['step'] == 0

    def test_onboarding_profile_created(self, user):
        profile = OnboardingProfile.objects.create(user=user, step=0)
        assert profile.step == 0
        assert str(profile) == f'{user.username}: paso 0'

    def test_advance_step(self, client, user):
        client.login(username='test@chele.app', password='testpass123')
        OnboardingProfile.objects.create(user=user, step=1)
        from django.middleware.csrf import get_token
        response = client.post(reverse('onboarding_advance'))
        assert response.status_code == 200
        data = response.json()
        assert data['step'] == 2

    def test_new_user_starts_at_step_0(self, user):
        profile = OnboardingProfile.objects.create(user=user)
        assert profile.step == 0

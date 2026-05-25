import pytest
from django.urls import reverse
from django.contrib.auth.models import User


@pytest.mark.django_db
class TestAuth:
    def test_registration(self, client):
        response = client.post(reverse('register'), {
            'name': 'Nuevo',
            'email': 'nuevo@test.com',
            'password': 'securepass123',
        })
        assert response.status_code == 302
        assert User.objects.filter(email='nuevo@test.com').exists()

    def test_duplicate_registration(self, client, user):
        response = client.post(reverse('register'), {
            'name': 'Otro',
            'email': 'test@chele.app',
            'password': 'securepass123',
        })
        assert response.status_code == 200
        assert b'registrado' in response.content.lower() or b'email' in response.content.lower()

    def test_login_page(self, client):
        response = client.get(reverse('login'))
        assert response.status_code == 200

    def test_login_success(self, client, user):
        response = client.post(reverse('login'), {
            'username': 'test@chele.app',
            'password': 'testpass123',
        })
        assert response.status_code == 302

    def test_login_failure(self, client):
        response = client.post(reverse('login'), {
            'username': 'wrong@test.com',
            'password': 'wrongpass',
        })
        assert response.status_code == 200

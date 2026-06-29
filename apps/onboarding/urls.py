from django.urls import path
from . import views

urlpatterns = [
    path('state/', views.onboarding_state, name='onboarding_state'),
    path('avanzar/', views.advance_step, name='onboarding_advance'),
    path('descartar/', views.dismiss_condition, name='onboarding_dismiss'),
]

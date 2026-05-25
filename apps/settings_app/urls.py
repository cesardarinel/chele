from django.urls import path
from . import views

urlpatterns = [
    path('', views.settings_index, name='settings_index'),
    path('perfil/', views.settings_profile, name='settings_profile'),
    path('presupuesto/', views.settings_budget, name='settings_budget'),
    path('invitar/', views.settings_invite, name='settings_invite'),
]

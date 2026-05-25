from django.urls import path
from . import views

urlpatterns = [
    path('', views.payees_list, name='payees_list'),
    path('crear/', views.payee_create, name='payee_create'),
    path('<uuid:id>/editar/', views.payee_edit, name='payee_edit'),
    path('<uuid:id>/fusionar/', views.payee_merge, name='payee_merge'),
]

from django.urls import path
from . import views

urlpatterns = [
    path('', views.account_list, name='account_list'),
    path('crear/', views.account_create, name='account_create'),
    path('<uuid:id>/', views.account_detail, name='account_detail'),
    path('<uuid:id>/editar/', views.account_edit, name='account_edit'),
    path('<uuid:id>/eliminar/', views.account_delete, name='account_delete'),
]

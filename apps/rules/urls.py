from django.urls import path
from . import views

urlpatterns = [
    path('', views.rules_list, name='rules_list'),
    path('crear/', views.rule_create, name='rule_create'),
    path('<uuid:id>/editar/', views.rule_edit, name='rule_edit'),
    path('<uuid:id>/eliminar/', views.rule_delete, name='rule_delete'),
]

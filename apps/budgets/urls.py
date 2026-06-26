from django.urls import path
from . import views

urlpatterns = [
    path('', views.budget_view, name='budget_view'),
    path('crear/', views.budget_create, name='budget_create'),
    path('cambiar/', views.switch_budget, name='switch_budget'),
    path('categorias/crear/', views.category_create, name='category_create'),
    path('categorias/<uuid:id>/editar/', views.category_edit, name='category_edit'),
    path('categorias/<uuid:id>/eliminar/', views.category_delete, name='category_delete'),
    path('grupos/crear/', views.category_group_create, name='category_group_create'),
    path('grupos/<uuid:id>/editar/', views.category_group_edit, name='category_group_edit'),
    path('grupos/<uuid:id>/eliminar/', views.category_group_delete, name='category_group_delete'),
    path('categoria/<uuid:id>/inspector/', views.category_inspector, name='category_inspector'),
    path('asignar/', views.assign_funds, name='assign_funds'),
    path('mover/', views.move_funds, name='move_funds'),
    path('cubrir/', views.cover_overspending, name='cover_overspending'),
    path('auto-asignar/', views.auto_assign, name='auto_assign'),
    path('reservar/', views.hold_for_next_month, name='hold_for_next_month'),
    path('copiar/', views.copy_budget, name='copy_budget'),
]

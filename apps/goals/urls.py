from django.urls import path
from . import views

urlpatterns = [
    path('', views.goal_list, name='goal_list'),
    path('crear/', views.goal_create, name='goal_create'),
    path('<uuid:id>/editar/', views.goal_edit, name='goal_edit'),
    path('<uuid:id>/eliminar/', views.goal_delete, name='goal_delete'),
]

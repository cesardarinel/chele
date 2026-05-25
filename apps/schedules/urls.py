from django.urls import path
from . import views

urlpatterns = [
    path('', views.schedules_list, name='schedules_list'),
    path('crear/', views.schedule_create, name='schedule_create'),
    path('<uuid:id>/editar/', views.schedule_edit, name='schedule_edit'),
    path('<uuid:id>/eliminar/', views.schedule_delete, name='schedule_delete'),
]

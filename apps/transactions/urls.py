from django.urls import path
from . import views

urlpatterns = [
    path('crear/', views.transaction_create, name='transaction_create'),
    path('<uuid:id>/editar/', views.transaction_edit, name='transaction_edit'),
    path('<uuid:id>/eliminar/', views.transaction_delete, name='transaction_delete'),
    path('masiva/', views.transaction_bulk, name='transaction_bulk'),
]

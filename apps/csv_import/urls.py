from django.urls import path
from . import views

urlpatterns = [
    path('', views.csv_import, name='csv_import'),
    path('preview/', views.csv_preview, name='csv_preview'),
    path('confirmar/', views.csv_confirm, name='csv_confirm'),
]

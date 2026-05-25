from django.urls import path
from . import views

urlpatterns = [
    path('', views.sync_now, name='sync_now'),
    path('push/', views.sync_push, name='sync_push'),
    path('pull/', views.sync_pull, name='sync_pull'),
]

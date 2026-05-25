from django.urls import path
from . import views

urlpatterns = [
    path('', views.cc_list, name='cc_list'),
    path('crear/', views.cc_create, name='cc_create'),
    path('<uuid:id>/', views.cc_detail, name='cc_detail'),
    path('pagar/', views.cc_pay, name='cc_pay'),
]

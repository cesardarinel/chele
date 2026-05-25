from django.urls import path
from . import views

urlpatterns = [
    path('', views.reports_dashboard, name='reports_dashboard'),
    path('net-worth/', views.net_worth_report, name='net_worth_report'),
    path('cash-flow/', views.cash_flow_report, name='cash_flow_report'),
    path('budget-vs-reality/', views.budget_vs_reality_report, name='budget_vs_reality'),
]

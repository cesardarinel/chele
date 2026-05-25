from django.urls import path
from . import views

urlpatterns = [
    path('', views.loan_list, name='loan_list'),
    path('crear/', views.loan_create, name='loan_create'),
    path('<uuid:id>/', views.loan_detail, name='loan_detail'),
    path('<uuid:id>/editar/', views.loan_edit, name='loan_edit'),
    path('<uuid:id>/eliminar/', views.loan_delete, name='loan_delete'),
    path('<uuid:id>/pagar-cuota/', views.loan_pay_installment, name='loan_pay_installment'),
]

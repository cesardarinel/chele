from django.contrib import admin
from django.urls import path, include
from django.contrib.auth import views as auth_views
from django.views.generic import RedirectView

urlpatterns = [
    path('admin/', admin.site.urls),
    path('', RedirectView.as_view(pattern_name='budget_view', permanent=False), name='home'),
    path('presupuestos/', include('apps.budgets.urls')),
    path('cuentas/', include('apps.accounts.urls')),
    path('transacciones/', include('apps.transactions.urls')),
    path('reportes/', include('apps.reports.urls')),
    path('metas/', include('apps.goals.urls')),
    path('importar/', include('apps.csv_import.urls')),
    path('programaciones/', include('apps.schedules.urls')),
    path('beneficiarios/', include('apps.payees.urls')),
    path('reglas/', include('apps.rules.urls')),
    path('configuracion/', include('apps.settings_app.urls')),
    path('sync/', include('apps.sync_engine.urls')),
    path('onboarding/', include('apps.onboarding.urls')),
    path('tc/', include('apps.credit_cards.urls')),
    path('prestamos/', include('apps.loans.urls')),
    path('registro/', include('apps.budgets.urls_registration')),
    path('acceder/', auth_views.LoginView.as_view(template_name='registration/login.html'), name='login'),
    path('salir/', auth_views.LogoutView.as_view(), name='logout'),
]

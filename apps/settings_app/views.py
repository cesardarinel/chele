from datetime import datetime
from django.shortcuts import render, redirect, get_object_or_404
from django.contrib.auth.decorators import login_required
from django.contrib import messages
from django.contrib.auth.models import User
from apps.budgets.models import Budget, BudgetMembership


@login_required
def settings_index(request):
    return render(request, 'settings_app/index.html')


@login_required
def settings_profile(request):
    if request.method == 'POST':
        user = request.user
        user.first_name = request.POST.get('name', user.first_name)
        user.email = request.POST.get('email', user.email)
        if request.POST.get('password'):
            user.set_password(request.POST.get('password'))
        user.save()
        messages.success(request, '🐷 Perfil actualizado.')
        return redirect('settings_index')
    return render(request, 'settings_app/profile.html')


@login_required
def settings_budget(request):
    budget_id = request.session.get('active_budget_id')
    budget = Budget.objects.filter(id=budget_id, owner=request.user).first()
    if request.method == 'POST' and request.POST.get('name'):
        budget.name = request.POST.get('name', budget.name)
        budget.description = request.POST.get('description', '')
        budget.save()
        messages.success(request, '🐷 Presupuesto actualizado.')
        return redirect('settings_budget')
    members = BudgetMembership.objects.filter(budget_id=budget_id).select_related('user')
    return render(request, 'settings_app/budget.html', {
        'budget': budget,
        'members': members,
    })


@login_required
def settings_invite(request):
    if request.method == 'POST':
        budget_id = request.session.get('active_budget_id')
        budget = get_object_or_404(Budget, id=budget_id, owner=request.user)
        email = request.POST.get('email', '').strip()
        if not email:
            messages.error(request, '🐷 Ingresá un email.')
            return redirect('settings_budget')
        try:
            user = User.objects.get(email=email)
        except User.DoesNotExist:
            messages.error(request, f'🐷 No existe un usuario con email {email}.')
            return redirect('settings_budget')
        if BudgetMembership.objects.filter(user=user, budget=budget).exists():
            messages.error(request, f'🐷 {email} ya es miembro de este presupuesto.')
            return redirect('settings_budget')
        BudgetMembership.objects.create(
            user=user, budget=budget, role='editor',
            invited_at=datetime.now(), accepted_at=datetime.now(),
        )
        messages.success(request, f'🐷 {email} agregado al presupuesto.')
    return redirect('settings_budget')

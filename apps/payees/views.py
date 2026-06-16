from django.shortcuts import render, redirect, get_object_or_404
from django.contrib.auth.decorators import login_required
from django.contrib import messages
from .models import Payee
from apps.budgets.models import Budget


@login_required
def payees_list(request):
    budget_id = request.session.get('active_budget_id')
    if not budget_id:
        first = Budget.objects.filter(members=request.user).first()
        if first:
            budget_id = str(first.id)
            request.session['active_budget_id'] = budget_id
    payees = Payee.objects.filter(budget_id=budget_id) if budget_id else []
    return render(request, 'payees/payee_list.html', {'payees': payees})


@login_required
def payee_create(request):
    if request.method == 'POST':
        budget_id = request.session.get('active_budget_id')
        if not budget_id:
            first = Budget.objects.filter(members=request.user).first()
            if first:
                budget_id = str(first.id)
                request.session['active_budget_id'] = budget_id
        if not budget_id:
            messages.error(request, 'Primero creá un presupuesto antes de agregar beneficiarios.')
            return redirect('budget_create')
        Payee.objects.create(budget_id=budget_id, name=request.POST.get('name'))
        messages.success(request, 'Beneficiario creado.')
        return redirect('payees_list')
    return render(request, 'payees/payee_form.html')


@login_required
def payee_edit(request, id):
    payee = get_object_or_404(Payee, id=id)
    if request.method == 'POST':
        payee.name = request.POST.get('name', payee.name)
        payee.save()
        messages.success(request, 'Beneficiario actualizado.')
        return redirect('payees_list')
    return render(request, 'payees/payee_form.html', {'payee': payee})


@login_required
def payee_merge(request, id):
    payee = get_object_or_404(Payee, id=id)
    if request.method == 'POST':
        target_id = request.POST.get('merge_to')
        target = get_object_or_404(Payee, id=target_id)
        payee.transactions.update(payee=target)
        payee.merge_to = target
        payee.save()
        messages.success(request, f'"{payee.name}" fusionado con "{target.name}".')
        return redirect('payees_list')
    return render(request, 'payees/payee_merge.html', {'payee': payee})

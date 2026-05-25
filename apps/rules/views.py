from django.shortcuts import render, redirect, get_object_or_404
from django.contrib.auth.decorators import login_required
from django.contrib import messages
from .models import Rule


@login_required
def rules_list(request):
    budget_id = request.session.get('active_budget_id')
    rules = Rule.objects.filter(budget_id=budget_id).select_related('action_category', 'action_payee')
    return render(request, 'rules/rule_list.html', {'rules': rules})


@login_required
def rule_create(request):
    if request.method == 'POST':
        budget_id = request.session.get('active_budget_id')
        Rule.objects.create(
            budget_id=budget_id,
            name=request.POST.get('name'),
            condition_field=request.POST.get('condition_field'),
            condition_operator=request.POST.get('condition_operator'),
            condition_value=request.POST.get('condition_value'),
            action_category_id=request.POST.get('action_category') or None,
            action_payee_id=request.POST.get('action_payee') or None,
            action_notes=request.POST.get('action_notes', ''),
        )
        messages.success(request, 'Regla creada.')
        return redirect('rules_list')
    return render(request, 'rules/rule_form.html')


@login_required
def rule_edit(request, id):
    rule = get_object_or_404(Rule, id=id)
    if request.method == 'POST':
        rule.name = request.POST.get('name', rule.name)
        rule.condition_field = request.POST.get('condition_field')
        rule.condition_operator = request.POST.get('condition_operator')
        rule.condition_value = request.POST.get('condition_value')
        rule.action_category_id = request.POST.get('action_category') or None
        rule.action_payee_id = request.POST.get('action_payee') or None
        rule.action_notes = request.POST.get('action_notes', '')
        rule.save()
        messages.success(request, 'Regla actualizada.')
        return redirect('rules_list')
    return render(request, 'rules/rule_form.html', {'rule': rule})


@login_required
def rule_delete(request, id):
    rule = get_object_or_404(Rule, id=id)
    if request.method == 'POST':
        rule.delete()
        messages.success(request, 'Regla eliminada.')
    return redirect('rules_list')

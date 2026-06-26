from django.shortcuts import render, redirect, get_object_or_404
from django.contrib.auth.decorators import login_required
from django.contrib import messages
from .models import CreditCard, InterestCharge
from apps.transactions.models import Transaction
from apps.accounts.models import Account
from apps.budgets.models import CategoryGroup, Category
from core.interest import aplicar_interes
from datetime import date


@login_required
def cc_pay(request):
    if request.method == 'POST':
        budget_id = request.session.get('active_budget_id')
        from_account_id = request.POST.get('from_account_id')
        cc_id = request.POST.get('cc_id')
        amount = float(request.POST.get('amount', 0))

        cc = get_object_or_404(CreditCard, id=cc_id, budget_id=budget_id)

        aplicar_interes(cc, cc.last_interest_date or date.today())

        Transaction.objects.create(
            budget_id=budget_id,
            account_id=from_account_id,
            date=request.POST.get('date'),
            amount=-amount,
            notes=f'Pago TC {cc.name}',
        )

        cc.balance = float(cc.balance) + amount
        cc.last_interest_date = date.today()
        cc.save()

        from_acct = Account.objects.get(id=from_account_id)
        from_acct.balance = float(from_acct.balance) - amount
        from_acct.save()

        messages.success(request, f'Pago a {cc.name} registrado.')
    return redirect('budget_view')


@login_required
def cc_create(request):
    if request.method == 'POST':
        budget_id = request.session.get('active_budget_id')
        if not budget_id:
            messages.error(request, 'Primero creá un presupuesto.')
            return redirect('budget_create')
        cc = CreditCard.objects.create(
            budget_id=budget_id,
            name=request.POST.get('name'),
            limit=request.POST.get('limit', 0),
            balance=request.POST.get('balance', 0),
            interest_rate=request.POST.get('interest_rate', 0),
            closing_day=request.POST.get('closing_day', 15),
            due_day=request.POST.get('due_day', 5),
            notes=request.POST.get('notes', ''),
        )

        # Create Payment category for this credit card (for TC auto-move)
        from django.utils.text import slugify
        payment_group, _ = CategoryGroup.objects.get_or_create(
            budget_id=budget_id, name='Pagos TC', is_income=False,
            defaults={'sort_order': 99}
        )
        payment_cat_name = f'Pago {cc.name}'
        if not Category.objects.filter(budget_id=budget_id, name=payment_cat_name).exists():
            Category.objects.create(
                budget_id=budget_id, group=payment_group,
                name=payment_cat_name,
                sort_order=0, notes=f'Categoría de pago para {cc.name}',
            )

        messages.success(request, f'Tarjeta "{cc.name}" creada.')
        return redirect('cc_list')
    return render(request, 'credit_cards/form.html')


@login_required
def cc_list(request):
    budget_id = request.session.get('active_budget_id')
    cards = CreditCard.objects.filter(budget_id=budget_id) if budget_id else []
    return render(request, 'credit_cards/list.html', {'cards': cards})


@login_required
def cc_detail(request, id):
    budget_id = request.session.get('active_budget_id')
    card = get_object_or_404(CreditCard, id=id, budget_id=budget_id)
    interes = card.calcular_interes()
    return render(request, 'credit_cards/detail.html', {
        'card': card,
        'interes_calculado': interes,
    })

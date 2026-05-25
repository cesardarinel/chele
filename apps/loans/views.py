from datetime import date
from django.shortcuts import render, redirect, get_object_or_404
from django.contrib.auth.decorators import login_required
from django.contrib import messages
from .models import Loan, Installment
from apps.transactions.models import Transaction
from core.interest import aplicar_interes


@login_required
def loan_list(request):
    budget_id = request.session.get('active_budget_id')
    loans = Loan.objects.filter(budget_id=budget_id)
    return render(request, 'loans/list.html', {'loans': loans})


@login_required
def loan_create(request):
    if request.method == 'POST':
        budget_id = request.session.get('active_budget_id')
        loan = Loan.objects.create(
            budget_id=budget_id,
            type=request.POST.get('type'),
            name=request.POST.get('name'),
            total_amount=request.POST.get('total_amount'),
            interest_rate=request.POST.get('interest_rate'),
            remaining_balance=request.POST.get('total_amount'),
            total_installments=request.POST.get('total_installments'),
            start_date=request.POST.get('start_date'),
            next_due_date=request.POST.get('next_due_date'),
            installment_amount=request.POST.get('installment_amount'),
            notes=request.POST.get('notes', ''),
        )

        total = int(request.POST.get('total_installments', 0))
        amount = float(request.POST.get('installment_amount', 0))
        start = date.fromisoformat(request.POST.get('next_due_date'))
        for i in range(total):
            due = start.replace(month=((start.month - 1 + i) % 12) + 1,
                                year=start.year + (start.month - 1 + i) // 12)
            Installment.objects.create(
                loan=loan, number=i + 1, amount=amount, due_date=due
            )

        Transaction.objects.create(
            budget_id=budget_id, account=None,
            date=request.POST.get('start_date'),
            amount=request.POST.get('total_amount'),
            notes=f'Préstamo {loan.name} - {loan.get_type_display()}',
        )
        messages.success(request, f'🐷 Préstamo "{loan.name}" registrado.')
        return redirect('loan_list')
    return render(request, 'loans/form.html')


@login_required
def loan_detail(request, id):
    budget_id = request.session.get('active_budget_id')
    loan = get_object_or_404(Loan, id=id, budget_id=budget_id)
    installments = Installment.objects.filter(loan=loan)
    interes = loan.calcular_interes()
    return render(request, 'loans/detail.html', {
        'loan': loan,
        'installments': installments,
        'interes_calculado': interes,
    })


@login_required
def loan_edit(request, id):
    budget_id = request.session.get('active_budget_id')
    loan = get_object_or_404(Loan, id=id, budget_id=budget_id)
    if request.method == 'POST':
        loan.name = request.POST.get('name', loan.name)
        loan.type = request.POST.get('type', loan.type)
        loan.total_amount = request.POST.get('total_amount', loan.total_amount)
        loan.interest_rate = request.POST.get('interest_rate', loan.interest_rate)
        loan.total_installments = request.POST.get('total_installments', loan.total_installments)
        loan.installment_amount = request.POST.get('installment_amount', loan.installment_amount)
        loan.start_date = request.POST.get('start_date') or loan.start_date
        loan.next_due_date = request.POST.get('next_due_date') or loan.next_due_date
        loan.notes = request.POST.get('notes', '')
        loan.save()
        messages.success(request, f'🐷 Préstamo "{loan.name}" actualizado.')
        return redirect('loan_detail', id=id)
    return render(request, 'loans/form.html', {'loan': loan, 'editing': True})


@login_required
def loan_delete(request, id):
    budget_id = request.session.get('active_budget_id')
    loan = get_object_or_404(Loan, id=id, budget_id=budget_id)
    if request.method == 'POST':
        name = loan.name
        loan.delete()
        messages.success(request, f'🐷 Préstamo "{name}" eliminado.')
        return redirect('loan_list')
    return redirect('loan_detail', id=id)


@login_required
def loan_pay_installment(request, id):
    if request.method == 'POST':
        budget_id = request.session.get('active_budget_id')
        installment_id = request.POST.get('installment_id')
        installment = get_object_or_404(Installment, id=installment_id, loan_id=id)

        aplicar_interes(installment.loan, installment.loan.next_due_date)

        Transaction.objects.create(
            budget_id=budget_id, account=None,
            date=date.today(),
            amount=-float(installment.amount),
            notes=f'Cuota {installment.number} - {installment.loan.name}',
        )

        installment.paid = True
        installment.paid_date = date.today()
        installment.save()

        loan = installment.loan
        loan.paid_installments += 1
        loan.remaining_balance = max(0, float(loan.remaining_balance) - float(installment.amount))
        if loan.remaining_balance <= 0:
            loan.next_due_date = None
            loan.status = 'completed'
            messages.success(request, f'🐷 ¡Préstamo "{loan.name}" completamente pagado!')
        else:
            if loan.paid_installments < loan.total_installments:
                next_inst = Installment.objects.filter(loan=loan, paid=False).order_by('number').first()
                if next_inst:
                    loan.next_due_date = next_inst.due_date
            messages.success(request, f'🐷 Cuota {installment.number} pagada.')
        loan.save()
    return redirect('loan_detail', id=id)

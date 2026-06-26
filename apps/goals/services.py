from decimal import Decimal
from datetime import date, datetime
from django.db.models import Sum
from apps.transactions.models import Transaction
from apps.budgets.models import MonthlyBudget


class TargetService:
    def __init__(self, budget, month, year):
        self.budget = budget
        self.month = month
        self.year = year

    def get_category_available(self, category_id):
        budgeted = MonthlyBudget.objects.filter(
            category_id=category_id, month=self.month, year=self.year
        ).aggregate(Sum('budgeted'))['budgeted__sum'] or 0

        spent = Transaction.objects.filter(
            category_id=category_id,
            date__month=self.month, date__year=self.year,
        ).aggregate(Sum('amount'))['amount__sum'] or 0
        if spent < 0:
            spent = abs(spent)
        else:
            spent = 0

        return float(budgeted) - float(spent)

    def calculate_underfunded(self, goal):
        if goal.snooze_month == self.month and goal.snooze_year == self.year:
            return 0

        amount = float(goal.amount)

        if goal.goal_type == 'monthly':
            if goal.refill_up_to:
                rollover = self._get_rollover(goal.category_id)
                return max(0, amount - rollover)
            return amount

        elif goal.goal_type == 'yearly':
            monthly = amount / 12
            if goal.refill_up_to:
                rollover = self._get_rollover(goal.category_id)
                return max(0, monthly - rollover)
            return monthly

        elif goal.goal_type == 'target_balance':
            avail = self.get_category_available(goal.category_id)
            return max(0, amount - avail)

        elif goal.goal_type == 'target_date':
            if not goal.target_date:
                return amount
            td = goal.target_date
            now = date(self.year, self.month, 1)
            if td < now:
                return 0
            months_left = (td.year - now.year) * 12 + (td.month - now.month) + 1
            if months_left <= 0:
                months_left = 1
            already = self.get_category_available(goal.category_id)
            return max(0, (amount - already) / months_left)

        elif goal.goal_type == 'true_expense':
            freq = goal.frequency or 12
            monthly = amount / freq
            if goal.refill_up_to:
                rollover = self._get_rollover(goal.category_id)
                return max(0, monthly - rollover)
            return monthly

        return 0

    def _get_rollover(self, category_id):
        prev_month = self.month - 1
        prev_year = self.year
        if prev_month < 1:
            prev_month = 12
            prev_year -= 1

        ts = TargetService(self.budget, prev_month, prev_year)
        avail = ts.get_category_available(category_id)
        return max(0, avail)

    def list_underfunded(self):
        from .models import Goal
        result = []
        goals = Goal.objects.filter(
            category__budget=self.budget, is_completed=False
        ).select_related('category')
        for g in goals:
            deficit = self.calculate_underfunded(g)
            if deficit > 0:
                result.append({
                    'category_id': str(g.category_id),
                    'category_name': g.category.name,
                    'deficit': round(deficit, 2),
                    'goal': g,
                })
        return result

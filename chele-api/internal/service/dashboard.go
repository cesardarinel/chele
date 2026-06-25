package service

import (
	"fmt"
	"math"

	"github.com/cesardarinel/chele-api/internal/models"
	"github.com/jmoiron/sqlx"
)

type DashboardService struct {
	DB *sqlx.DB
}

func NewDashboardService(db *sqlx.DB) *DashboardService {
	return &DashboardService{DB: db}
}

func (s *DashboardService) GetBudgetView(budgetID string, month, year, rangeVal int) (*models.DashboardResponse, error) {
	var budget models.Budget
	if err := s.DB.Get(&budget, "SELECT * FROM budgets_budget WHERE id=?", budgetID); err != nil {
		return nil, err
	}

	var totalBalance float64
	s.DB.Get(&totalBalance, "SELECT COALESCE(SUM(balance),0) FROM accounts_account WHERE budget_id=? AND on_budget=1", budgetID)

	var months []models.DashboardMonth
	for i := -rangeVal; i <= rangeVal; i++ {
		m, y := month+i, year
		for m > 12 { m -= 12; y++ }
		for m < 1 { m += 12; y-- }

		var income, expenses float64
		s.DB.Get(&income,
			"SELECT COALESCE(SUM(amount),0) FROM transactions_transaction WHERE budget_id=? AND amount>0 AND CAST(strftime('%m',date) AS INTEGER)=? AND CAST(strftime('%Y',date) AS INTEGER)=?",
			budgetID, m, y)
		s.DB.Get(&expenses,
			"SELECT COALESCE(SUM(amount),0) FROM transactions_transaction WHERE budget_id=? AND amount<0 AND CAST(strftime('%m',date) AS INTEGER)=? AND CAST(strftime('%Y',date) AS INTEGER)=?",
			budgetID, m, y)

		var scheduledIncome float64
		s.DB.Get(&scheduledIncome,
			"SELECT COALESCE(SUM(amount),0) FROM schedules_schedule WHERE budget_id=? AND is_active=1 AND direction='income' AND CAST(strftime('%m',next_date) AS INTEGER)=? AND CAST(strftime('%Y',next_date) AS INTEGER)=?",
			budgetID, m, y)

		var totalBudgeted float64
		s.DB.Get(&totalBudgeted,
			"SELECT COALESCE(SUM(budgeted),0) FROM budgets_monthlybudget mb JOIN budgets_category c ON c.id=mb.category_id WHERE c.budget_id=? AND mb.month=? AND mb.year=?",
			budgetID, m, y)

		dm := models.DashboardMonth{
			Month:          m,
			Year:           y,
			Income:         income,
			ScheduledIncome: scheduledIncome,
			Expenses:       math.Abs(expenses),
			TotalBudgeted:  totalBudgeted,
		}
		dm.Income += scheduledIncome
		dm.AvailableFunds = dm.CarriedOver + dm.Income
		dm.ForNextMonth = dm.AvailableFunds - dm.TotalBudgeted

		if len(months) > 0 {
			prev := months[len(months)-1]
			dm.CarriedOver = math.Round((prev.Income-prev.TotalBudgeted-prev.Expenses)*100)/100
		}
		dm.AvailableFunds = math.Round((dm.CarriedOver+dm.Income)*100)/100
		dm.ForNextMonth = math.Round((dm.AvailableFunds-dm.TotalBudgeted)*100)/100

		monthNames := []string{"Ene","Feb","Mar","Abr","May","Jun","Jul","Ago","Sep","Oct","Nov","Dic"}
		dm.Label = fmt.Sprintf("%s %d", monthNames[m-1], y)
		months = append(months, dm)
	}

	var groups []models.CategoryGroup
	s.DB.Select(&groups, "SELECT * FROM budgets_categorygroup WHERE budget_id=? ORDER BY sort_order", budgetID)

	var dashboardGroups []models.DashboardGroup
	for _, g := range groups {
		var cats []models.Category
		s.DB.Select(&cats, "SELECT * FROM budgets_category WHERE group_id=? AND is_hidden=0 ORDER BY sort_order", g.ID)
		var dCats []models.DashboardCategory
		for _, cat := range cats {
			var catMonths []models.DashboardCatMonth
			for _, vm := range months {
				var budgeted float64
				s.DB.Get(&budgeted,
					"SELECT COALESCE(budgeted,0) FROM budgets_monthlybudget WHERE category_id=? AND month=? AND year=?",
					cat.ID, vm.Month, vm.Year)
				var spent float64
				s.DB.Get(&spent,
					"SELECT COALESCE(SUM(amount),0) FROM transactions_transaction WHERE category_id=? AND CAST(strftime('%m',date) AS INTEGER)=? AND CAST(strftime('%Y',date) AS INTEGER)=?",
					cat.ID, vm.Month, vm.Year)
				if spent < 0 { spent = -spent }
				catMonths = append(catMonths, models.DashboardCatMonth{
					Month:    vm.Month,
					Year:     vm.Year,
					Budgeted: budgeted,
					Spent:    spent,
					Balance:  budgeted - spent,
				})
			}
			dCats = append(dCats, models.DashboardCategory{Category: cat, Months: catMonths})
		}
		dashboardGroups = append(dashboardGroups, models.DashboardGroup{Group: g, Categories: dCats})
	}

	availableToBudget := 0.0
	if len(months) > 0 {
		a := months[len(months)-1-rangeVal]
		availableToBudget = math.Round((a.Income+a.CarriedOver-a.TotalBudgeted)*100)/100
	}
	if availableToBudget == 0 && totalBalance > 0 {
		availableToBudget = math.Round(totalBalance*100)/100
	}

	return &models.DashboardResponse{
		Budget:            budget,
		Groups:            dashboardGroups,
		Months:            months,
		AvailableToBudget: availableToBudget,
		TotalBalance:      totalBalance,
	}, nil
}

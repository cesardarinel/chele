package service

import (
	"math"
	"time"

	"github.com/jmoiron/sqlx"
)

type TargetService struct {
	DB *sqlx.DB
}

func NewTargetService(db *sqlx.DB) *TargetService {
	return &TargetService{DB: db}
}

func (s *TargetService) GetCategoryAvailable(categoryID string, month, year int) float64 {
	var budgeted float64
	s.DB.Get(&budgeted,
		"SELECT COALESCE(budgeted,0) FROM budgets_monthlybudget WHERE category_id=? AND month=? AND year=?",
		categoryID, month, year)
	var spent float64
	s.DB.Get(&spent,
		"SELECT COALESCE(SUM(amount),0) FROM transactions_transaction WHERE category_id=? AND CAST(strftime('%m',date) AS INTEGER)=? AND CAST(strftime('%Y',date) AS INTEGER)=?",
		categoryID, month, year)
	if spent < 0 {
		spent = -spent
	}
	return budgeted - spent
}

func (s *TargetService) getLastMonthRollover(categoryID string, month, year int) float64 {
	prevMonth, prevYear := month-1, year
	if prevMonth < 1 {
		prevMonth = 12
		prevYear--
	}
	avail := s.GetCategoryAvailable(categoryID, prevMonth, prevYear)
	if avail < 0 {
		avail = 0
	}
	return avail
}

type TargetInfo struct {
	GoalType     string  `json:"goal_type"`
	Amount       float64 `json:"amount"`
	RefillUpTo   bool    `json:"refill_up_to"`
	Underfunded  float64 `json:"underfunded"`
	IsSnoozed    bool    `json:"is_snoozed"`
}

func (s *TargetService) CalculateUnderfunded(goalType string, amount float64, refillUpTo bool, snoozeMonth, snoozeYear *int, categoryID string, month, year int, frequency int, targetDate *string) float64 {
	if snoozeMonth != nil && snoozeYear != nil && *snoozeMonth == month && *snoozeYear == year {
		return 0
	}

	switch goalType {
	case "monthly":
		if refillUpTo {
			rollover := s.getLastMonthRollover(categoryID, month, year)
			return math.Max(0, amount-rollover)
		}
		return amount

	case "yearly":
		monthly := amount / 12
		if refillUpTo {
			rollover := s.getLastMonthRollover(categoryID, month, year)
			return math.Max(0, monthly-rollover)
		}
		return monthly

	case "target_balance":
		avail := s.GetCategoryAvailable(categoryID, month, year)
		return math.Max(0, amount-avail)

	case "target_date":
		if targetDate == nil {
			return amount
		}
		td, err := time.Parse("2006-01-02", *targetDate)
		if err != nil {
			return amount
		}
		now := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
		if td.Before(now) {
			return 0
		}
		monthsRemaining := int(td.Sub(now).Hours()/24/30) + 1
		if monthsRemaining <= 0 {
			monthsRemaining = 1
		}
		alreadyAssigned := s.GetCategoryAvailable(categoryID, month, year)
		return math.Max(0, (amount-alreadyAssigned)/float64(monthsRemaining))

	case "true_expense":
		if frequency <= 0 {
			frequency = 12
		}
		monthly := amount / float64(frequency)
		if refillUpTo {
			rollover := s.getLastMonthRollover(categoryID, month, year)
			return math.Max(0, monthly-rollover)
		}
		return monthly
	}

	return 0
}

func (s *TargetService) ListUnderfunded(budgetID string, month, year int) []struct {
	CategoryID  string  `json:"category_id"`
	CategoryName string `json:"category_name"`
	Deficit     float64 `json:"deficit"`
} {
	type catTarget struct {
		ID         string  `db:"id"`
		CategoryID string  `db:"category_id"`
		Name       string  `db:"name"`
		GoalType   string  `db:"goal_type"`
		Amount     float64 `db:"amount"`
		RefillUpTo bool    `db:"refill_up_to"`
		SnoozeM    *int    `db:"snooze_month"`
		SnoozeY    *int    `db:"snooze_year"`
		Frequency  int     `db:"frequency"`
		TargetDate *string `db:"target_date"`
	}

	var rows []catTarget
	s.DB.Select(&rows,
		`SELECT g.id,g.category_id,c.name,g.goal_type,g.amount,g.refill_up_to,g.snooze_month,g.snooze_year,g.frequency,g.target_date
		 FROM goals_goal g JOIN budgets_category c ON c.id=g.category_id
		 WHERE c.budget_id=? AND g.is_completed=0`, budgetID)

	var result []struct {
		CategoryID  string  `json:"category_id"`
		CategoryName string `json:"category_name"`
		Deficit     float64 `json:"deficit"`
	}

	for _, r := range rows {
		deficit := s.CalculateUnderfunded(r.GoalType, r.Amount, r.RefillUpTo, r.SnoozeM, r.SnoozeY, r.CategoryID, month, year, r.Frequency, r.TargetDate)
		if deficit > 0 {
			result = append(result, struct {
				CategoryID  string  `json:"category_id"`
				CategoryName string `json:"category_name"`
				Deficit     float64 `json:"deficit"`
			}{r.CategoryID, r.Name, deficit})
		}
	}
	return result
}

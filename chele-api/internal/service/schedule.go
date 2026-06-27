package service

import (
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type ScheduleService struct {
	DB *sqlx.DB
}

func NewScheduleService(db *sqlx.DB) *ScheduleService {
	return &ScheduleService{DB: db}
}

type dueSchedule struct {
	ID                 string  `db:"id"`
	AccountID          string  `db:"account_id"`
	BudgetID           string  `db:"budget_id"`
	Amount             float64 `db:"amount"`
	Direction          string  `db:"direction"`
	NextDate           string  `db:"next_date"`
	Frequency          string  `db:"frequency"`
	SkipWeekends       bool    `db:"skip_weekends"`
	ApplyBeforeWeekend bool    `db:"apply_before_weekend"`
	PayeeID            *string `db:"payee_id"`
	CatID              *string `db:"category_id"`
}

func (s *ScheduleService) ProcessDue(budgetID string) (int, error) {
	today := time.Now().UTC().Format("2006-01-02")
	var due []dueSchedule
	err := s.DB.Select(&due,
		`SELECT id,account_id,budget_id,amount,direction,next_date,frequency,skip_weekends,apply_before_weekend,payee_id,category_id
		 FROM schedules_schedule WHERE budget_id=? AND is_active=1 AND next_date<=?`,
		budgetID, today)
	if err != nil {
		return 0, err
	}

	count := 0
	for _, sch := range due {
		amount := sch.Amount
		if sch.Direction == "expense" {
			amount = -amount
		}

		actualDate := adjustDate(sch.NextDate, sch.SkipWeekends, sch.ApplyBeforeWeekend)

		txnID := strings.ReplaceAll(uuid.New().String(), "-", "")
		s.DB.Exec(
			`INSERT INTO transactions_transaction
			 (id,budget_id,account_id,date,amount,payee_id,category_id,notes,reconciled,cleared,created_at,updated_at)
			 VALUES (?,?,?,?,?,?,?,?,0,0,datetime('now'),datetime('now'))`,
			txnID, sch.BudgetID, sch.AccountID, actualDate, amount, sch.PayeeID, sch.CatID, "Programación",
		)
		s.DB.Exec("UPDATE accounts_account SET balance = balance + ? WHERE id = ?", amount, sch.AccountID)
		newDate := AdvanceDate(sch.NextDate, sch.Frequency, sch.SkipWeekends, sch.ApplyBeforeWeekend)
		s.DB.Exec("UPDATE schedules_schedule SET next_date = ? WHERE id = ?", newDate, sch.ID)
		count++
	}
	return count, nil
}

func adjustDate(dateStr string, skipWeekends, applyBefore bool) string {
	dt, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return dateStr
	}
	if applyBefore {
		if dt.Weekday() == time.Saturday {
			dt = dt.AddDate(0, 0, -1)
		}
		if dt.Weekday() == time.Sunday {
			dt = dt.AddDate(0, 0, -2)
		}
	} else if skipWeekends {
		if dt.Weekday() == time.Saturday {
			dt = dt.AddDate(0, 0, 2)
		}
		if dt.Weekday() == time.Sunday {
			dt = dt.AddDate(0, 0, 1)
		}
	}
	return dt.Format("2006-01-02")
}

func AdvanceDate(currentDate, frequency string, skipWeekends, applyBefore bool) string {
	dt, err := time.Parse("2006-01-02", currentDate)
	if err != nil {
		return currentDate
	}
	switch frequency {
	case "weekly":
		dt = dt.AddDate(0, 0, 7)
	case "biweekly":
		dt = dt.AddDate(0, 0, 14)
	case "monthly":
		dt = dt.AddDate(0, 1, 0)
	case "quarterly":
		dt = dt.AddDate(0, 3, 0)
	case "yearly":
		dt = dt.AddDate(1, 0, 0)
	}
	if applyBefore {
		if dt.Weekday() == time.Saturday {
			dt = dt.AddDate(0, 0, -1)
		}
		if dt.Weekday() == time.Sunday {
			dt = dt.AddDate(0, 0, -2)
		}
	} else if skipWeekends {
		for dt.Weekday() == time.Saturday {
			dt = dt.AddDate(0, 0, 2)
		}
		for dt.Weekday() == time.Sunday {
			dt = dt.AddDate(0, 0, 1)
		}
	}
	return dt.Format("2006-01-02")
}

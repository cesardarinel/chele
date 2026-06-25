package service

import "github.com/jmoiron/sqlx"

type BalanceService struct {
	DB *sqlx.DB
}

func NewBalanceService(db *sqlx.DB) *BalanceService {
	return &BalanceService{DB: db}
}

func (s *BalanceService) Apply(accountID string, amount float64) error {
	_, err := s.DB.Exec(
		"UPDATE accounts_account SET balance = balance + ? WHERE id = ?",
		amount, accountID,
	)
	return err
}

func (s *BalanceService) Reverse(accountID string, amount float64) error {
	_, err := s.DB.Exec(
		"UPDATE accounts_account SET balance = balance - ? WHERE id = ?",
		amount, accountID,
	)
	return err
}

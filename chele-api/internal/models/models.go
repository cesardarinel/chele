package models

import "time"

// accounts_account
type Account struct {
	ID        string    `db:"id" json:"id"`
	BudgetID  string    `db:"budget_id" json:"budget_id"`
	Name      string    `db:"name" json:"name"`
	OnBudget  bool      `db:"on_budget" json:"on_budget"`
	Balance   float64   `db:"balance" json:"balance"`
	Notes     string    `db:"notes" json:"notes"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}

// budgets_budget
type Budget struct {
	ID          string    `db:"id" json:"id"`
	Name        string    `db:"name" json:"name"`
	Description string    `db:"description" json:"description"`
	OwnerID     int       `db:"owner_id" json:"owner_id"`
	CreatedAt   time.Time `db:"created_at" json:"created_at"`
	UpdatedAt   time.Time `db:"updated_at" json:"updated_at"`
}

// budgets_budgetmembership
type BudgetMembership struct {
	ID          string     `db:"id" json:"id"`
	UserID      int        `db:"user_id" json:"user_id"`
	BudgetID    string     `db:"budget_id" json:"budget_id"`
	Role        string     `db:"role" json:"role"`
	InvitedAt   time.Time  `db:"invited_at" json:"invited_at"`
	AcceptedAt  *time.Time `db:"accepted_at" json:"accepted_at"`
}

// budgets_categorygroup
type CategoryGroup struct {
	ID        string    `db:"id" json:"id"`
	BudgetID  string    `db:"budget_id" json:"budget_id"`
	Name      string    `db:"name" json:"name"`
	SortOrder int       `db:"sort_order" json:"sort_order"`
	IsIncome  bool      `db:"is_income" json:"is_income"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}

// budgets_category
type Category struct {
	ID        string    `db:"id" json:"id"`
	BudgetID  string    `db:"budget_id" json:"budget_id"`
	GroupID   string    `db:"group_id" json:"group_id"`
	Name      string    `db:"name" json:"name"`
	SortOrder int       `db:"sort_order" json:"sort_order"`
	IsHidden  bool      `db:"is_hidden" json:"is_hidden"`
	Notes     string    `db:"notes" json:"notes"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}

// budgets_monthlybudget
type MonthlyBudget struct {
	ID        string    `db:"id" json:"id"`
	CategoryID string   `db:"category_id" json:"category_id"`
	Month     int       `db:"month" json:"month"`
	Year      int       `db:"year" json:"year"`
	Budgeted  float64   `db:"budgeted" json:"budgeted"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}

// transactions_transaction
type Transaction struct {
	ID         string     `db:"id" json:"id"`
	BudgetID   string     `db:"budget_id" json:"budget_id"`
	AccountID  *string    `db:"account_id" json:"account_id"`
	Date       string     `db:"date" json:"date"`
	Amount     float64    `db:"amount" json:"amount"`
	PayeeID    *string    `db:"payee_id" json:"payee_id"`
	CategoryID *string    `db:"category_id" json:"category_id"`
	Notes      string     `db:"notes" json:"notes"`
	TransferID *string    `db:"transfer_id" json:"transfer_id"`
	Reconciled bool       `db:"reconciled" json:"reconciled"`
	Cleared    bool       `db:"cleared" json:"cleared"`
	CreatedAt  time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt  time.Time  `db:"updated_at" json:"updated_at"`
}

// schedules_schedule
type Schedule struct {
	ID           string    `db:"id" json:"id"`
	BudgetID     string    `db:"budget_id" json:"budget_id"`
	PayeeID      *string   `db:"payee_id" json:"payee_id"`
	CategoryID   *string   `db:"category_id" json:"category_id"`
	AccountID    string    `db:"account_id" json:"account_id"`
	Amount       float64   `db:"amount" json:"amount"`
	Frequency    string    `db:"frequency" json:"frequency"`
	NextDate     string    `db:"next_date" json:"next_date"`
	Notes        string    `db:"notes" json:"notes"`
	SkipWeekends bool      `db:"skip_weekends" json:"skip_weekends"`
	IsActive     bool      `db:"is_active" json:"is_active"`
	Direction    string    `db:"direction" json:"direction"`
	CreatedAt    time.Time `db:"created_at" json:"created_at"`
	UpdatedAt    time.Time `db:"updated_at" json:"updated_at"`
}

// payees_payee
type Payee struct {
	ID        string     `db:"id" json:"id"`
	BudgetID  string     `db:"budget_id" json:"budget_id"`
	Name      string     `db:"name" json:"name"`
	MergeTo   *string    `db:"merge_to" json:"merge_to"`
	CreatedAt time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt time.Time  `db:"updated_at" json:"updated_at"`
}

// goals_goal
type Goal struct {
	ID          string     `db:"id" json:"id"`
	CategoryID  string     `db:"category_id" json:"category_id"`
	GoalType    string     `db:"goal_type" json:"goal_type"`
	Amount      float64    `db:"amount" json:"amount"`
	TargetDate  *string    `db:"target_date" json:"target_date"`
	Frequency   int        `db:"frequency" json:"frequency"`
	IsCompleted bool       `db:"is_completed" json:"is_completed"`
	CreatedAt   time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt   time.Time  `db:"updated_at" json:"updated_at"`
}

// credit_cards_creditcard
type CreditCard struct {
	ID              string     `db:"id" json:"id"`
	BudgetID        string     `db:"budget_id" json:"budget_id"`
	Name            string     `db:"name" json:"name"`
	Limit           float64    `db:"limit" json:"limit"`
	Balance         float64    `db:"balance" json:"balance"`
	InterestRate    float64    `db:"interest_rate" json:"interest_rate"`
	ClosingDay      int        `db:"closing_day" json:"closing_day"`
	DueDay          int        `db:"due_day" json:"due_day"`
	LastInterestDate *string   `db:"last_interest_date" json:"last_interest_date"`
	Notes           string     `db:"notes" json:"notes"`
	CreatedAt       time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt       time.Time  `db:"updated_at" json:"updated_at"`
}

// credit_cards_interestcharge
type InterestCharge struct {
	ID           string    `db:"id" json:"id"`
	CreditCardID string    `db:"credit_card_id" json:"credit_card_id"`
	Amount       float64   `db:"amount" json:"amount"`
	Date         string    `db:"date" json:"date"`
	DaysOverdue  int       `db:"days_overdue" json:"days_overdue"`
	Applied      bool      `db:"applied" json:"applied"`
}

// loans_loan
type Loan struct {
	ID                string    `db:"id" json:"id"`
	BudgetID          string    `db:"budget_id" json:"budget_id"`
	AccountID         *string   `db:"account_id" json:"account_id"`
	Type              string    `db:"type" json:"type"`
	Name              string    `db:"name" json:"name"`
	Status            string    `db:"status" json:"status"`
	TotalAmount       float64   `db:"total_amount" json:"total_amount"`
	InterestRate      float64   `db:"interest_rate" json:"interest_rate"`
	RemainingBalance  float64   `db:"remaining_balance" json:"remaining_balance"`
	TotalInstallments int       `db:"total_installments" json:"total_installments"`
	PaidInstallments  int       `db:"paid_installments" json:"paid_installments"`
	StartDate         string    `db:"start_date" json:"start_date"`
	NextDueDate       string    `db:"next_due_date" json:"next_due_date"`
	InstallmentAmount float64   `db:"installment_amount" json:"installment_amount"`
	Notes             string    `db:"notes" json:"notes"`
	CreatedAt         time.Time `db:"created_at" json:"created_at"`
	UpdatedAt         time.Time `db:"updated_at" json:"updated_at"`
}

// loans_installment
type Installment struct {
	ID       string  `db:"id" json:"id"`
	LoanID   string  `db:"loan_id" json:"loan_id"`
	Number   int     `db:"number" json:"number"`
	Amount   float64 `db:"amount" json:"amount"`
	DueDate  string  `db:"due_date" json:"due_date"`
	Paid     bool    `db:"paid" json:"paid"`
	PaidDate *string `db:"paid_date" json:"paid_date"`
	Notes    string  `db:"notes" json:"notes"`
}

// rules_rule
type Rule struct {
	ID               string    `db:"id" json:"id"`
	BudgetID         string    `db:"budget_id" json:"budget_id"`
	Name             string    `db:"name" json:"name"`
	ConditionField   string    `db:"condition_field" json:"condition_field"`
	ConditionOperator string   `db:"condition_operator" json:"condition_operator"`
	ConditionValue   string    `db:"condition_value" json:"condition_value"`
	ActionCategoryID *string   `db:"action_category_id" json:"action_category_id"`
	ActionPayeeID    *string   `db:"action_payee_id" json:"action_payee_id"`
	ActionNotes      string    `db:"action_notes" json:"action_notes"`
	SortOrder        int       `db:"sort_order" json:"sort_order"`
	IsActive         bool      `db:"is_active" json:"is_active"`
	CreatedAt        time.Time `db:"created_at" json:"created_at"`
	UpdatedAt        time.Time `db:"updated_at" json:"updated_at"`
}

// auth_user
type User struct {
	ID           int        `db:"id" json:"id"`
	Password     string     `db:"password" json:"-"`
	LastLogin    *time.Time `db:"last_login" json:"last_login"`
	IsSuperuser  bool       `db:"is_superuser" json:"-"`
	Username     string     `db:"username" json:"username"`
	FirstName    string     `db:"first_name" json:"first_name"`
	LastName     string     `db:"last_name" json:"last_name"`
	Email        string     `db:"email" json:"email"`
	IsStaff      bool       `db:"is_staff" json:"-"`
	IsActive     bool       `db:"is_active" json:"is_active"`
	DateJoined   time.Time  `db:"date_joined" json:"date_joined"`
}

// sync_engine_synclog
type SyncLog struct {
	ID         string    `db:"id" json:"id"`
	BudgetID   string    `db:"budget_id" json:"budget_id"`
	UserID     int       `db:"user_id" json:"user_id"`
	EntityType string    `db:"entity_type" json:"entity_type"`
	EntityID   string    `db:"entity_id" json:"entity_id"`
	Action     string    `db:"action" json:"action"`
	Payload    string    `db:"payload" json:"payload"`
	CreatedAt  time.Time `db:"created_at" json:"created_at"`
}

// Helper response types
type AccountSummary struct {
	OnBudgetAccounts  []Account `json:"on_budget"`
	OffBudgetAccounts []Account `json:"off_budget"`
	TotalOnBudget     float64   `json:"total_on_budget"`
	TotalOffBudget    float64   `json:"total_off_budget"`
	GrandTotal        float64   `json:"grand_total"`
	CCDebtCards       []CreditCard `json:"cc_debt_cards"`
	TotalCCDebt       float64   `json:"total_cc_debt"`
	Loans             []Loan    `json:"loans"`
	TotalLoanDebt     float64   `json:"total_loan_debt"`
	TotalDebt         float64   `json:"total_debt"`
}

type DashboardMonth struct {
	Month          int     `json:"month"`
	Year           int     `json:"year"`
	Active         bool    `json:"active"`
	Label          string  `json:"label"`
	Income         float64 `json:"income"`
	ScheduledIncome float64 `json:"scheduled_income"`
	Expenses       float64 `json:"expenses"`
	TotalBudgeted  float64 `json:"total_budgeted"`
	CarriedOver    float64 `json:"carried_over"`
	AvailableFunds float64 `json:"available_funds"`
	ForNextMonth   float64 `json:"for_next_month"`
}

type DashboardResponse struct {
	Budget            Budget                `json:"budget"`
	Groups            []DashboardGroup      `json:"groups"`
	Months            []DashboardMonth      `json:"months"`
	AvailableToBudget float64               `json:"available_to_budget"`
	TotalBalance      float64               `json:"total_balance"`
}

type DashboardGroup struct {
	Group      CategoryGroup    `json:"group"`
	Categories []DashboardCategory `json:"categories"`
}

type DashboardCategory struct {
	Category Category          `json:"category"`
	Months   []DashboardCatMonth `json:"months"`
}

type DashboardCatMonth struct {
	Month    int     `json:"month"`
	Year     int     `json:"year"`
	Budgeted float64 `json:"budgeted"`
	Spent    float64 `json:"spent"`
	Balance  float64 `json:"balance"`
}

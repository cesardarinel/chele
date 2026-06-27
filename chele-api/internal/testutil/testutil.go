package testutil

import (
	"log"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

func SetupDB() *sqlx.DB {
	db, err := sqlx.Connect("sqlite3", ":memory:")
	if err != nil {
		log.Fatal(err)
	}
	db.MustExec("PRAGMA foreign_keys=ON")
	db.MustExec(schema)
	return db
}

const schema = `
CREATE TABLE auth_user (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	password varchar(128) NOT NULL,
	last_login datetime NULL,
	is_superuser bool NOT NULL,
	username varchar(150) NOT NULL UNIQUE,
	last_name varchar(150) NOT NULL,
	email varchar(254) NOT NULL,
	is_staff bool NOT NULL,
	is_active bool NOT NULL,
	date_joined datetime NOT NULL,
	first_name varchar(150) NOT NULL
);
CREATE TABLE budgets_budget (
	id char(32) NOT NULL PRIMARY KEY,
	name varchar(200) NOT NULL,
	description text NOT NULL,
	owner_id integer NOT NULL REFERENCES auth_user(id),
	created_at datetime NOT NULL,
	updated_at datetime NOT NULL
);
CREATE TABLE budgets_budgetmembership (
	id char(32) NOT NULL PRIMARY KEY,
	user_id integer NOT NULL REFERENCES auth_user(id),
	budget_id char(32) NOT NULL REFERENCES budgets_budget(id),
	role varchar(20) NOT NULL,
	invited_at datetime NOT NULL,
	accepted_at datetime NULL
);
CREATE TABLE budgets_categorygroup (
	id char(32) NOT NULL PRIMARY KEY,
	budget_id char(32) NOT NULL REFERENCES budgets_budget(id),
	name varchar(200) NOT NULL,
	sort_order integer NOT NULL,
	is_income bool NOT NULL,
	created_at datetime NOT NULL,
	updated_at datetime NOT NULL
);
CREATE TABLE budgets_category (
	id char(32) NOT NULL PRIMARY KEY,
	budget_id char(32) NOT NULL REFERENCES budgets_budget(id),
	group_id char(32) NOT NULL REFERENCES budgets_categorygroup(id),
	name varchar(200) NOT NULL,
	sort_order integer NOT NULL,
	is_hidden bool NOT NULL,
	notes text NOT NULL,
	created_at datetime NOT NULL,
	updated_at datetime NOT NULL
);
CREATE TABLE budgets_monthlybudget (
	id char(32) NOT NULL PRIMARY KEY,
	category_id char(32) NOT NULL REFERENCES budgets_category(id),
	month integer NOT NULL,
	year integer NOT NULL,
	budgeted decimal NOT NULL,
	created_at datetime NOT NULL,
	updated_at datetime NOT NULL
);
CREATE TABLE accounts_account (
	id char(32) NOT NULL PRIMARY KEY,
	budget_id char(32) NOT NULL REFERENCES budgets_budget(id),
	name varchar(200) NOT NULL,
	on_budget bool NOT NULL,
	balance decimal NOT NULL,
	notes text NOT NULL,
	created_at datetime NOT NULL,
	updated_at datetime NOT NULL
);
CREATE TABLE transactions_transaction (
	id char(32) NOT NULL PRIMARY KEY,
	budget_id char(32) NOT NULL REFERENCES budgets_budget(id),
	account_id char(32) NULL REFERENCES accounts_account(id),
	date date NOT NULL,
	amount decimal NOT NULL,
	payee_id char(32) NULL,
	category_id char(32) NULL REFERENCES budgets_category(id),
	notes text NOT NULL,
	transfer_id char(32) NULL,
	reconciled bool NOT NULL,
	cleared bool NOT NULL,
	created_at datetime NOT NULL,
	updated_at datetime NOT NULL
);
CREATE TABLE schedules_schedule (
	id char(32) NOT NULL PRIMARY KEY,
	budget_id char(32) NOT NULL REFERENCES budgets_budget(id),
	payee_id char(32) NULL,
	category_id char(32) NULL,
	account_id char(32) NOT NULL REFERENCES accounts_account(id),
	amount decimal NOT NULL,
	frequency varchar(20) NOT NULL,
	next_date date NOT NULL,
	notes text NOT NULL,
	skip_weekends bool NOT NULL,
	apply_before_weekend bool NOT NULL,
	is_active bool NOT NULL,
	direction varchar(10) NOT NULL,
	created_at datetime NOT NULL,
	updated_at datetime NOT NULL
);
CREATE TABLE credit_cards_creditcard (
	id char(32) NOT NULL PRIMARY KEY,
	budget_id char(32) NOT NULL REFERENCES budgets_budget(id),
	name varchar(200) NOT NULL,
	"limit" decimal NOT NULL,
	balance decimal NOT NULL,
	interest_rate decimal NOT NULL,
	closing_day integer NOT NULL,
	due_day integer NOT NULL,
	last_interest_date date NULL,
	notes text NOT NULL,
	created_at datetime NOT NULL,
	updated_at datetime NOT NULL
);
CREATE TABLE loans_loan (
	id char(32) NOT NULL PRIMARY KEY,
	budget_id char(32) NOT NULL REFERENCES budgets_budget(id),
	account_id char(32) NULL REFERENCES accounts_account(id),
	type varchar(20) NOT NULL,
	name varchar(200) NOT NULL,
	status varchar(20) NOT NULL,
	total_amount decimal NOT NULL,
	interest_rate decimal NOT NULL,
	remaining_balance decimal NOT NULL,
	total_installments integer NOT NULL,
	paid_installments integer NOT NULL,
	start_date date NOT NULL,
	next_due_date date NOT NULL,
	installment_amount decimal NOT NULL,
	notes text NOT NULL,
	created_at datetime NOT NULL,
	updated_at datetime NOT NULL
);
CREATE TABLE goals_goal (
	id char(32) NOT NULL PRIMARY KEY,
	category_id char(32) NOT NULL REFERENCES budgets_category(id),
	goal_type varchar(20) NOT NULL,
	amount decimal NOT NULL,
	target_date date NULL,
	frequency integer NOT NULL,
	is_completed bool NOT NULL,
	refill_up_to bool NOT NULL,
	snooze_month integer NULL,
	snooze_year integer NULL,
	created_at datetime NOT NULL,
	updated_at datetime NOT NULL
);
CREATE TABLE loans_installment (
	id char(32) NOT NULL PRIMARY KEY,
	loan_id char(32) NOT NULL REFERENCES loans_loan(id),
	number integer NOT NULL,
	amount decimal NOT NULL,
	due_date date NOT NULL,
	paid bool NOT NULL,
	paid_date date NULL,
	notes text NOT NULL
);
`

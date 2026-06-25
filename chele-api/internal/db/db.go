package db

import (
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	_ "github.com/jackc/pgx/v5/stdlib"
)

func Connect(driver, dsn string) *sqlx.DB {
	var d string
	switch driver {
	case "sqlite":
		d = "sqlite3"
	case "postgres":
		d = "pgx"
	default:
		log.Fatalf("unsupported DB_DRIVER: %s", driver)
	}
	db, err := sqlx.Connect(d, dsn)
	if err != nil {
		log.Fatalf("db connect: %v", err)
	}
	db.SetMaxOpenConns(10)
	if driver == "sqlite" {
		db.MustExec("PRAGMA journal_mode=WAL")
		db.MustExec("PRAGMA foreign_keys=ON")
	}
	fmt.Printf("connected to %s (%s)\n", driver, dsn)
	return db
}

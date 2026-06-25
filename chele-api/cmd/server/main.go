package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/cesardarinel/chele-api/internal/config"
	"github.com/cesardarinel/chele-api/internal/db"
	"github.com/cesardarinel/chele-api/internal/router"
)

func main() {
	cfg := config.Load()
	database := db.Connect(cfg.DBDriver, cfg.DatabaseURL)
	r := router.New(database, cfg)

	addr := fmt.Sprintf(":%d", cfg.Port)
	log.Printf("chele-api listening on %s", addr)
	if err := http.ListenAndServe(addr, r); err != nil {
		log.Fatal(err)
	}
}

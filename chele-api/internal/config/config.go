package config

import (
	"os"
	"strconv"
)

type Config struct {
	DBDriver  string
	DatabaseURL string
	JWTSecert string
	Port      int
}

func Load() *Config {
	port, _ := strconv.Atoi(os.Getenv("PORT"))
	if port == 0 {
		port = 8080
	}
	return &Config{
		DBDriver:   getEnv("DB_DRIVER", "sqlite"),
		DatabaseURL: getEnv("DATABASE_URL", "db/sqlite.db"),
		JWTSecert:  getEnv("JWT_SECRET", "dev-secret-change-in-production"),
		Port:       port,
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

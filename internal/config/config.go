package config

import (
	"os"
)

type Config struct {
	DatabaseURL string
	Port        string
	AdminPassword string
}

func Load() *Config {
	return &Config{
		DatabaseURL:   getEnv("DATABASE_URL", "postgres://dev:dev@localhost:5432/personal_site?sslmode=disable"),
		Port:          getEnv("PORT", "3000"),
		AdminPassword: getEnv("ADMIN_PASSWORD", "admin"),
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

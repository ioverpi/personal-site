package config

import (
	"os"
	"strconv"
)

type Config struct {
	DatabaseURL          string
	Port                 string
	SecureCookies        bool   // Set to true in production (HTTPS)
	SessionDurationHours int    // How long sessions last
	BaseURL              string // For invite links
}

func Load() *Config {
	return &Config{
		DatabaseURL:          getEnv("DATABASE_URL", "postgres://dev:dev@localhost:5432/personal_site?sslmode=disable"),
		Port:                 getEnv("PORT", "3000"),
		SecureCookies:        getEnvBool("SECURE_COOKIES", false),
		SessionDurationHours: getEnvInt("SESSION_DURATION_HOURS", 24*7), // 1 week default
		BaseURL:              getEnv("BASE_URL", "http://localhost:3000"),
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func getEnvBool(key string, fallback bool) bool {
	if value, ok := os.LookupEnv(key); ok {
		return value == "true" || value == "1"
	}
	return fallback
}

func getEnvInt(key string, fallback int) int {
	if value, ok := os.LookupEnv(key); ok {
		if i, err := strconv.Atoi(value); err == nil {
			return i
		}
	}
	return fallback
}

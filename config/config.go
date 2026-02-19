package config

import (
	"log/slog"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	AppPort string
}

func Load() Config {
	if os.Getenv("APP_ENV") != "production" {
		if err := godotenv.Load(); err != nil {
			slog.Warn("No .env file found, fallback used", "error", err)
		}
	}

	return Config{
		AppPort: getOrFallback("APP_PORT", "8080"),
	}
}

func getOrFallback(key, fallback string) string {
	if val, ok := os.LookupEnv(key); ok {
		return val
	}
	return fallback
}

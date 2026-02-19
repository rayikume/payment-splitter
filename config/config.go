package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	AppPort string
}

func Load() Config {
	if os.Getenv("APP_ENV") != "production" {
		if err := godotenv.Load(); err != nil {
			log.Println("No [.env] file found. Fallback used")
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

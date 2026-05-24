package config

import (
	"errors"
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port               string
	DatabaseURL        string
	RedisURL           string
	TradeLockerBaseURL string
	TradeLockerServer  string
	CORSOrigin         string
	SessionCookieName  string
}

var requiredEnvKeys = []string{
	"PORT",
	"DATABASE_URL",
	"REDIS_URL",
	"TRADELOCKER_BASE_URL",
	"TRADELOCKER_SERVER",
	"CORS_ORIGIN",
	"SESSION_COOKIE_NAME",
}

func loadEnvFile() error {
	if err := godotenv.Load(".env.local"); err != nil {
		if errors.Is(err, os.ErrNotExist) || os.IsNotExist(err) {
			return nil
		}
		if _, statErr := os.Stat(".env.local"); os.IsNotExist(statErr) {
			return nil
		}
		return fmt.Errorf("load .env.local: %w", err)
	}
	return nil
}

func Load() (Config, error) {
	if err := loadEnvFile(); err != nil {
		return Config{}, err
	}

	var missing []string
	for _, key := range requiredEnvKeys {
		if os.Getenv(key) == "" {
			missing = append(missing, key)
		}
	}
	if len(missing) > 0 {
		return Config{}, fmt.Errorf("missing required environment variables: %v", missing)
	}

	return Config{
		Port:               os.Getenv("PORT"),
		DatabaseURL:        os.Getenv("DATABASE_URL"),
		RedisURL:           os.Getenv("REDIS_URL"),
		TradeLockerBaseURL: os.Getenv("TRADELOCKER_BASE_URL"),
		TradeLockerServer:  os.Getenv("TRADELOCKER_SERVER"),
		CORSOrigin:         os.Getenv("CORS_ORIGIN"),
		SessionCookieName:  os.Getenv("SESSION_COOKIE_NAME"),
	}, nil
}

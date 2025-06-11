package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"github.com/korolev-n/merch/internal/logger"
)

type Config struct {
	DBdsn string
}

func NewConfig() (*Config, error) {
	err := godotenv.Load()
	if err != nil {
		logger.Log.Warn(".env not found")
	}

	dsn := os.Getenv("MERCH_DB_DSN")
	if dsn == "" {
		logger.Log.Error("MERCH_DB_DSN is not set")
		return nil, fmt.Errorf("MERCH_DB_DSN is not set")
	}

	logger.Log.Info("Configuration loaded")
	return &Config{
		DBdsn: dsn,
	}, nil
}

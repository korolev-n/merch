package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DBdsn string
}

func NewConfig() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found")
	}

	dsn := os.Getenv("MERCH_DB_DSN")
	if dsn == "" {
		log.Fatal("environment variable MERCH_DB_DSN is not set")
	}
	return &Config{
		DBdsn: dsn,
	}
}

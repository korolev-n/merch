package main

import (
	"database/sql"
	"os"

	"github.com/korolev-n/merch-auth/internal/config"
	"github.com/korolev-n/merch-auth/internal/logger"
	"github.com/korolev-n/merch-auth/internal/server"
	_ "github.com/lib/pq"
)

func main() {
	logger.Init()

	cfg, err := config.NewConfig()
	if err != nil {
		logger.Log.Error("failed to load config", "error", err)
		os.Exit(1)
	}

	db, err := openDB(cfg)
	if err != nil {
		logger.Log.Error("failed to open DB", "error", err)
		os.Exit(1)
	}
	defer db.Close()

	srv := server.New(db)
	srv.Run()
}

func openDB(cfg *config.Config) (*sql.DB, error) {
	db, err := sql.Open("postgres", cfg.DBdsn)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}

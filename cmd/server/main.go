package main

import (
	"database/sql"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/korolev-n/merch-auth/internal/config"
	"github.com/korolev-n/merch-auth/internal/repository/mocks"
	"github.com/korolev-n/merch-auth/internal/service"
	transport "github.com/korolev-n/merch-auth/internal/transport/http"
	_ "github.com/lib/pq"
)

func main() {

	cfg := config.NewConfig()

	db, err := openDB(cfg)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// userRepo := repository.NewUserRepository(db)
	// walletRepo := repository.NewWalletRepository(db)

	userRepo := mocks.NewMockUserRepository()
	walletRepo := mocks.NewMockWalletRepository()
	regService := service.NewRegistrationService(userRepo, walletRepo)
	handler := &transport.Handler{Reg: regService}

	r := gin.Default()
	r.POST("/api/auth", handler.Register)

	if err := r.Run(); err != nil {
		log.Fatalf("could not start server: %v", err)
	}
}

func openDB(cfg *config.Config) (*sql.DB, error) {
	db, err := sql.Open("postgres", cfg.DBdsn)
	if err != nil {
		return nil, err
	}

	return db, nil
}

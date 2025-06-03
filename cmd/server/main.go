package main

import (
	"log"

	"github.com/gin-gonic/gin"

	"github.com/korolev-n/merch-auth/internal/repository/mocks"
	"github.com/korolev-n/merch-auth/internal/service"

	transport "github.com/korolev-n/merch-auth/internal/transport/http"
)

func main() {
	// db, _ := sql.Open("postgres", "postgres://user:pass@localhost:5432/dbname?sslmode=disable")
	// defer db.Close()

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

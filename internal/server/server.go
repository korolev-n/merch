package server

import (
	"database/sql"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/korolev-n/merch-auth/internal/repository"
	"github.com/korolev-n/merch-auth/internal/service"
	transport "github.com/korolev-n/merch-auth/internal/transport/http"
)

type Server struct {
	db     *sql.DB
	router *gin.Engine
}

func New(db *sql.DB) *Server {

	// userRepo := mocks.NewMockUserRepository()
	// walletRepo := mocks.NewMockWalletRepository()
	
	userRepo := repository.NewUserRepository(db)
	walletRepo := repository.NewWalletRepository(db)
	regService := service.NewRegistrationService(userRepo, walletRepo)
	handler := &transport.Handler{Reg: regService}

	router := gin.Default()
	router.POST("/api/auth", handler.Register)

	return &Server{
		db:     db,
		router: router,
	}
}

func (s *Server) Run() {
	log.Println("Starting server on :8080")
	if err := s.router.Run(":8080"); err != nil {
		log.Fatalf("could not start server: %v", err)
	}
}

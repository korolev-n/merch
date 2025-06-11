package server

import (
	"database/sql"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/korolev-n/merch/internal/logger"
	"github.com/korolev-n/merch/internal/repository"
	"github.com/korolev-n/merch/internal/service"
	transport "github.com/korolev-n/merch/internal/transport/http"
	"github.com/korolev-n/merch/internal/transport/http/middleware"
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
	jwtService := service.NewJWTService()
	regService := service.NewRegistrationService(userRepo, walletRepo, jwtService)
	transferService := service.NewTransferService(userRepo, walletRepo)
	shopRepo := repository.NewShopRepository(db)
	shopService := service.NewShopService(shopRepo)
	infoRepo := repository.NewInfoRepository(db)
	infoService := service.NewInfoService(infoRepo)
	handler := &transport.Handler{
		Reg:      regService,
		Transfer: transferService,
		Shop:     shopService,
		Info:     infoService,
	}

	router := gin.Default()
	router.POST("/api/auth", handler.Register)

	protected := router.Group("/api")
	protected.Use(middleware.AuthMiddleware(jwtService))
	{
		protected.POST("/sendCoin", handler.SendCoin)
		protected.GET("/buy/:item", handler.BuyItem)
		protected.GET("/info", handler.GetInfo)
	}

	return &Server{
		db:     db,
		router: router,
	}
}

func (s *Server) Run() {
	logger.Log.Info("Starting server", "addr", ":8080")
	if err := s.router.Run(":8080"); err != nil {
		logger.Log.Error("could not start server", "error", err)
		os.Exit(1)
	}
}

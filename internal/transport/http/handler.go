package http

import (
	"github.com/gin-gonic/gin"
	"github.com/korolev-n/merch-auth/internal/service"
	"github.com/korolev-n/merch-auth/internal/transport/http/request"
	"github.com/korolev-n/merch-auth/internal/transport/http/response"
)

type Handler struct {
	Reg *service.RegistrationService
}

type RegisterRequest struct {
    Email    string `json:"email"`
    Password string `json:"password"`
}

func (h *Handler) Register(c *gin.Context) {
	var req request.AuthRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "invalid input"})
		return
	}

	if err := h.Reg.RegisterUser(); err != nil {
		c.JSON(500, response.ErrorResponse{Errors: "could not register"})
		return
	}

	c.JSON(201, response.AuthResponse{Token: ""})
}

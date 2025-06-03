package http

import (
	"github.com/gin-gonic/gin"
	"github.com/korolev-n/merch-auth/internal/service"
)

type Handler struct {
	Reg *service.RegistrationService
}

type RegisterRequest struct {
    Email    string `json:"email"`
    Password string `json:"password"`
}

func (h *Handler) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "invalid input"})
		return
	}

	if err := h.Reg.RegisterUser(); err != nil {
		c.JSON(500, gin.H{"error": "could not register"})
		return
	}

	c.JSON(201, gin.H{"status": "registered"})
}

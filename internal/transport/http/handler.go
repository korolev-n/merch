package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/korolev-n/merch-auth/internal/logger"
	"github.com/korolev-n/merch-auth/internal/service"
	"github.com/korolev-n/merch-auth/internal/transport/http/helper"
	"github.com/korolev-n/merch-auth/internal/transport/http/request"
	"github.com/korolev-n/merch-auth/internal/transport/http/response"
)

type Handler struct {
	Reg service.Registration
}

type RegisterRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (h *Handler) Register(c *gin.Context) {
	var req request.AuthRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Log.Warn("Invalid input in register", "error", err)
		helper.JSONError(c, http.StatusBadRequest, "invalid input")
		return
	}

	ctx := c.Request.Context()

	token, err := h.Reg.RegisterUser(ctx, req.Username, req.Password)
	if err != nil {
		logger.Log.Warn("Registration failed", "username", req.Username, "error", err)
		
		switch err {
		case service.ErrIncorrectPassword:
			helper.JSONError(c, http.StatusUnauthorized, "incorrect password")
		case service.ErrUserAlreadyExists:
			helper.JSONError(c, http.StatusConflict, "user already exists")
		case service.ErrTokenGeneration:
			helper.JSONError(c, http.StatusInternalServerError, "token generation failed")
		default:
			helper.JSONError(c, http.StatusInternalServerError, "could not register")
		}
		return
	}

	c.JSON(200, response.AuthResponse{Token: token})
}

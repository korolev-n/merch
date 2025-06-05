package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/korolev-n/merch-auth/internal/service"
	"github.com/korolev-n/merch-auth/internal/transport/http/helper"
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
		helper.JSONError(c, http.StatusBadRequest, "invalid input")
		return
	}

	ctx := c.Request.Context()

	if err := h.Reg.RegisterUser(ctx, req.Username, req.Password); err != nil {
		helper.JSONError(c, http.StatusInternalServerError, "could not register")
		return
	}

	c.JSON(200, response.AuthResponse{Token: ""})
}

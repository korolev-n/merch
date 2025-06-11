package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/korolev-n/merch/internal/logger"
	"github.com/korolev-n/merch/internal/service"
	"github.com/korolev-n/merch/internal/transport/http/helper"
	"github.com/korolev-n/merch/internal/transport/http/request"
	"github.com/korolev-n/merch/internal/transport/http/response"
)

type Handler struct {
	Reg      service.Registration
	Transfer service.Transfer
	Shop     service.Shop
	Info     service.Info
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

	c.JSON(http.StatusOK, response.AuthResponse{Token: token})
}

func (h *Handler) SendCoin(c *gin.Context) {
	var req request.SendCoinRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helper.JSONError(c, http.StatusBadRequest, "invalid input")
		return
	}

	fromUserIDVal, exists := c.Get("user_id")
	if !exists {
		helper.JSONError(c, http.StatusUnauthorized, "missing user id")
		return
	}
	fromUserID := fromUserIDVal.(int)

	err := h.Transfer.SendCoins(c.Request.Context(), fromUserID, req.ToUser, req.Amount)
	if err != nil {
		switch err {
		case service.ErrUserNotFound:
			helper.JSONError(c, http.StatusNotFound, "recipient not found")
		default:
			helper.JSONError(c, http.StatusBadRequest, err.Error())
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func (h *Handler) BuyItem(c *gin.Context) {
	itemType := c.Param("item")
	userID, exists := c.Get("user_id")
	if !exists {
		helper.JSONError(c, http.StatusUnauthorized, "missing user id")
		return
	}

	err := h.Shop.BuyItem(c.Request.Context(), userID.(int), itemType)
	if err != nil {
		switch err {
		case service.ErrItemNotFound:
			helper.JSONError(c, http.StatusNotFound, "item not found")
		case service.ErrInsufficientBalance:
			helper.JSONError(c, http.StatusBadRequest, "not enough coins")
		default:
			helper.JSONError(c, http.StatusInternalServerError, "purchase failed")
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func (h *Handler) GetInfo(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		helper.JSONError(c, http.StatusUnauthorized, "missing user id")
		return
	}

	info, err := h.Info.GetUserInfo(c.Request.Context(), userID.(int))
	if err != nil {
		helper.JSONError(c, http.StatusInternalServerError, "failed to get info")
		return
	}

	resp := response.InfoResponse{
		Coins: info.Coins,
		Inventory: func() []response.InventoryItem {
			items := make([]response.InventoryItem, len(info.Inventory))
			for i, item := range info.Inventory {
				items[i] = response.InventoryItem{
					Type:     item.Type,
					Quantity: item.Quantity,
				}
			}
			return items
		}(),
		CoinHistory: response.CoinHistory{
			Received: func() []response.CoinReceived {
				received := make([]response.CoinReceived, len(info.CoinHistory.Received))
				for i, r := range info.CoinHistory.Received {
					received[i] = response.CoinReceived{
						FromUser: r.FromUser,
						Amount:   r.Amount,
					}
				}
				return received
			}(),
			Sent: func() []response.CoinSent {
				sent := make([]response.CoinSent, len(info.CoinHistory.Sent))
				for i, s := range info.CoinHistory.Sent {
					sent[i] = response.CoinSent{
						ToUser: s.ToUser,
						Amount: s.Amount,
					}
				}
				return sent
			}(),
		},
	}

	c.JSON(http.StatusOK, resp)
}

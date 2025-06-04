package helper

import (
	"github.com/gin-gonic/gin"
	"github.com/korolev-n/merch-auth/internal/transport/http/response"
)

func JSONError(c *gin.Context, code int, msg string) {
	c.JSON(code, response.ErrorResponse{Errors: msg})
}

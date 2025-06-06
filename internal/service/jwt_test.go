package service

import (
	"os"
	"testing"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
)

func TestJWTService_GenerateToken(t *testing.T) {
	// Arrange: установить переменную окружения для секрета
	secret := "test_secret_key"
	os.Setenv("JWT_SECRET", secret)

	jwtService := NewJWTService()

	userID := 42
	username := "testuser"

	tokenStr, err := jwtService.GenerateToken(userID, username)

	assert.NoError(t, err)
	assert.NotEmpty(t, tokenStr)

	parsedToken, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	assert.NoError(t, err)
	assert.True(t, parsedToken.Valid)

	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	assert.True(t, ok)

	assert.Equal(t, float64(userID), claims["user_id"])
	assert.Equal(t, username, claims["username"])
	_, hasExp := claims["exp"]
	assert.True(t, hasExp, "token should have expiration")
}

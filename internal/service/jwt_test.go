package service

import (
	"os"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
)

func TestJWTService_GenerateToken(t *testing.T) {
	secret := "test_secret_key"
	os.Setenv("JWT_SECRET", secret)

	jwtService := NewJWTService()

	userID := 42
	username := "testuser"

	tokenStr, err := jwtService.GenerateToken(userID, username)

	assert.NoError(t, err)
	assert.NotEmpty(t, tokenStr)

	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})

	assert.NoError(t, err)
	assert.True(t, token.Valid)

	assert.Equal(t, userID, claims.UserID)
	assert.Equal(t, username, claims.Username)
	assert.WithinDuration(t, time.Now().Add(24*time.Hour), claims.ExpiresAt.Time, time.Minute)
}

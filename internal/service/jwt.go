package service

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type TokenGenerator interface {
	GenerateToken(userID int, username string) (string, error)
	ParseToken(tokenStr string) (*Claims, error)
}

type JWTService struct {
	secretKey string
}

type Claims struct {
	UserID   int    `json:"user_id"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

func NewJWTService() *JWTService {
	return &JWTService{
		secretKey: os.Getenv("JWT_SECRET"),
	}
}

func (j *JWTService) GenerateToken(userID int, username string) (string, error) {
	claims := Claims{
		UserID:   userID,
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(j.secretKey))
}

func (j *JWTService) ParseToken(tokenStr string) (*Claims, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (interface{}, error) {
		return []byte(j.secretKey), nil
	})

	if err != nil || !token.Valid {
		return nil, ErrInvalidToken
	}

	return claims, nil
}

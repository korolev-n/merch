package http

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// Мок RegistrationService
type mockRegistrationService struct {
	RegisterUserFunc func(ctx context.Context, username, password string) (string, error)
}

func (m *mockRegistrationService) RegisterUser(ctx context.Context, username, password string) (string, error) {
	return m.RegisterUserFunc(ctx, username, password)
}

func TestHandler_Register_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := &mockRegistrationService{
		RegisterUserFunc: func(ctx context.Context, username, password string) (string, error) {
			return "mocked.jwt.token", nil
		},
	}

	handler := &Handler{Reg: mockService}

	router := gin.Default()
	router.POST("/api/auth", handler.Register)

	body := map[string]string{
		"username": "testuser",
		"password": "testpass",
	}
	jsonBody, _ := json.Marshal(body)
	req, _ := http.NewRequest(http.MethodPost, "/api/auth", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "mocked.jwt.token")
}

func TestHandler_Register_InvalidInput(t *testing.T) {
	gin.SetMode(gin.TestMode)

	handler := &Handler{Reg: &mockRegistrationService{}}
	router := gin.Default()
	router.POST("/api/auth", handler.Register)

	req, _ := http.NewRequest(http.MethodPost, "/api/auth", bytes.NewBufferString(`{invalid json}`))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "invalid input")
}

func TestHandler_Register_ServiceError(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := &mockRegistrationService{
		RegisterUserFunc: func(ctx context.Context, username, password string) (string, error) {
			return "", errors.New("db error")
		},
	}

	handler := &Handler{Reg: mockService}
	router := gin.Default()
	router.POST("/api/auth", handler.Register)

	body := map[string]string{
		"username": "user",
		"password": "pass",
	}
	jsonBody, _ := json.Marshal(body)
	req, _ := http.NewRequest(http.MethodPost, "/api/auth", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, w.Body.String(), "could not register")
}

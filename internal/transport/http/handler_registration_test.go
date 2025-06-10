package http_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/korolev-n/merch-auth/internal/logger"
	"github.com/korolev-n/merch-auth/internal/service"
	myhttp "github.com/korolev-n/merch-auth/internal/transport/http"
	"github.com/korolev-n/merch-auth/internal/transport/http/request"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Реализует интерфейс Registration
type mockRegistrationService struct {
	RegisterUserFunc func(ctx context.Context, username, password string) (string, error)
}

func (m *mockRegistrationService) RegisterUser(ctx context.Context, username, password string) (string, error) {
	return m.RegisterUserFunc(ctx, username, password)
}

// helper создает handler и роутер
func setupHandler(service service.Registration) *gin.Engine {
	handler := &myhttp.Handler{Reg: service}
	router := gin.Default()
	router.POST("/api/auth", handler.Register)
	return router
}

func TestHandler_Register(t *testing.T) {
	logger.Init()

	tests := []struct {
		name           string
		body           any
		mockService    service.Registration
		expectedCode   int
		expectedSubstr string
	}{
		{
			name: "valid request returns token",
			body: request.AuthRequest{Username: "user", Password: "pass"},
			mockService: &mockRegistrationService{
				RegisterUserFunc: func(ctx context.Context, username, password string) (string, error) {
					return "token.jwt.mock", nil
				},
			},
			expectedCode:   http.StatusOK,
			expectedSubstr: "token.jwt.mock",
		},
		{
			name: "invalid JSON returns 400",
			body: `{invalid json}`, // строка, не структура
			mockService: &mockRegistrationService{
				RegisterUserFunc: func(ctx context.Context, username, password string) (string, error) {
					return "", nil
				},
			},
			expectedCode:   http.StatusBadRequest,
			expectedSubstr: "invalid input",
		},
		{
			name: "internal error returns 500",
			body: request.AuthRequest{Username: "user", Password: "pass"},
			mockService: &mockRegistrationService{
				RegisterUserFunc: func(ctx context.Context, username, password string) (string, error) {
					return "", errors.New("db failure")
				},
			},
			expectedCode:   http.StatusInternalServerError,
			expectedSubstr: "could not register",
		},
		{
			name: "incorrect password returns 401",
			body: request.AuthRequest{Username: "user", Password: "wrongpass"},
			mockService: &mockRegistrationService{
				RegisterUserFunc: func(ctx context.Context, username, password string) (string, error) {
					return "", service.ErrIncorrectPassword
				},
			},
			expectedCode:   http.StatusUnauthorized,
			expectedSubstr: "incorrect password",
		},
		{
			name: "user exists returns 409",
			body: request.AuthRequest{Username: "user", Password: "pass"},
			mockService: &mockRegistrationService{
				RegisterUserFunc: func(ctx context.Context, username, password string) (string, error) {
					return "", service.ErrUserAlreadyExists
				},
			},
			expectedCode:   http.StatusConflict,
			expectedSubstr: "user already exists",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router := setupHandler(tt.mockService)

			var reqBody *bytes.Buffer
			switch b := tt.body.(type) {
			case string:
				reqBody = bytes.NewBufferString(b)
			default:
				j, err := json.Marshal(b)
				require.NoError(t, err)
				reqBody = bytes.NewBuffer(j)
			}

			req, err := http.NewRequest(http.MethodPost, "/api/auth", reqBody)
			require.NoError(t, err)
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedCode, w.Code)
			assert.Contains(t, w.Body.String(), tt.expectedSubstr)
		})
	}
}

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
	"github.com/korolev-n/merch/internal/service"
	myhttp "github.com/korolev-n/merch/internal/transport/http"
	"github.com/stretchr/testify/assert"
)

type mockTransferService struct {
	SendCoinsFunc func(ctx context.Context, fromUserID int, toUsername string, amount int) error
}

func (m *mockTransferService) SendCoins(ctx context.Context, fromUserID int, toUsername string, amount int) error {
	return m.SendCoinsFunc(ctx, fromUserID, toUsername, amount)
}

func TestHandler_SendCoin(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		body           interface{}
		setupMock      func() *mockTransferService
		expectedStatus int
		expectedBody   string
	}{
		{
			name: "success",
			body: map[string]interface{}{"toUser": "user2", "amount": 100},
			setupMock: func() *mockTransferService {
				return &mockTransferService{
					SendCoinsFunc: func(ctx context.Context, fromUserID int, toUsername string, amount int) error {
						assert.Equal(t, 1, fromUserID)
						assert.Equal(t, "user2", toUsername)
						assert.Equal(t, 100, amount)
						return nil
					},
				}
			},
			expectedStatus: http.StatusOK,
			expectedBody:   `"status":"ok"`,
		},
		{
			name: "invalid input",
			body: `{invalid json}`,
			setupMock: func() *mockTransferService {
				return &mockTransferService{}
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "invalid input",
		},
		{
			name: "recipient not found",
			body: map[string]interface{}{"toUser": "ghost", "amount": 50},
			setupMock: func() *mockTransferService {
				return &mockTransferService{
					SendCoinsFunc: func(ctx context.Context, fromUserID int, toUsername string, amount int) error {
						return service.ErrUserNotFound
					},
				}
			},
			expectedStatus: http.StatusNotFound,
			expectedBody:   "recipient not found",
		},
		{
			name: "unexpected error",
			body: map[string]interface{}{"toUser": "user2", "amount": 10},
			setupMock: func() *mockTransferService {
				return &mockTransferService{
					SendCoinsFunc: func(ctx context.Context, fromUserID int, toUsername string, amount int) error {
						return errors.New("db timeout")
					},
				}
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "db timeout",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockSvc := tt.setupMock()
			handler := &myhttp.Handler{Transfer: mockSvc}

			router := gin.New()
			router.POST("/api/sendCoin", func(c *gin.Context) {
				c.Set("user_id", 1)
				handler.SendCoin(c)
			})

			var reqBody *bytes.Buffer
			switch v := tt.body.(type) {
			case string:
				reqBody = bytes.NewBufferString(v)
			default:
				jsonData, _ := json.Marshal(v)
				reqBody = bytes.NewBuffer(jsonData)
			}

			req, _ := http.NewRequest(http.MethodPost, "/api/sendCoin", reqBody)
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			assert.Contains(t, w.Body.String(), tt.expectedBody)
		})
	}
}

package http_test

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type InventoryItem struct {
	ID    int    `json:"id"`
	Type  string `json:"type"`
	Price int    `json:"price"`
}

type InventoryService interface {
	BuyItem(userID int, itemType string) (*InventoryItem, error)
}

type MockInventoryService struct {
	mock.Mock
}

func (m *MockInventoryService) BuyItem(userID int, itemType string) (*InventoryItem, error) {
	args := m.Called(userID, itemType)
	item, _ := args.Get(0).(*InventoryItem)
	return item, args.Error(1)
}

func BuyItemHandler(service InventoryService) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := 1 // жестко задан для теста
		item := c.Param("item")

		result, err := service.BuyItem(userID, item)
		if err != nil {
			switch err.Error() {
			case "insufficient funds":
				c.JSON(http.StatusPaymentRequired, gin.H{"error": err.Error()})
			default:
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			}
			return
		}

		c.JSON(http.StatusOK, result)
	}
}

func TestBuyItemHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name         string
		item         string
		mockBehavior func(m *MockInventoryService)
		expectedCode int
		expectedBody string
	}{
		{
			name: "success",
			item: "sword",
			mockBehavior: func(m *MockInventoryService) {
				m.On("BuyItem", 1, "sword").Return(&InventoryItem{ID: 1, Type: "sword", Price: 100}, nil)
			},
			expectedCode: http.StatusOK,
			expectedBody: `{"id":1,"type":"sword","price":100}`,
		},
		{
			name: "item not found",
			item: "potion",
			mockBehavior: func(m *MockInventoryService) {
				m.On("BuyItem", 1, "potion").Return(nil, errors.New("item not found"))
			},
			expectedCode: http.StatusBadRequest,
			expectedBody: `{"error":"item not found"}`,
		},
		{
			name: "insufficient funds",
			item: "armor",
			mockBehavior: func(m *MockInventoryService) {
				m.On("BuyItem", 1, "armor").Return(nil, errors.New("insufficient funds"))
			},
			expectedCode: http.StatusPaymentRequired,
			expectedBody: `{"error":"insufficient funds"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := new(MockInventoryService)
			tt.mockBehavior(mockService)

			router := gin.Default()
			router.POST("/api/buy/:item", BuyItemHandler(mockService))

			req := httptest.NewRequest(http.MethodPost, "/api/buy/"+tt.item, nil)
			resp := httptest.NewRecorder()

			router.ServeHTTP(resp, req)

			assert.Equal(t, tt.expectedCode, resp.Code)

			// JSON сравнение без учёта порядка
			var expected, actual map[string]interface{}
			_ = json.Unmarshal([]byte(tt.expectedBody), &expected)
			_ = json.Unmarshal(resp.Body.Bytes(), &actual)
			assert.Equal(t, expected, actual)
		})
	}
}

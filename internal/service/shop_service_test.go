package service_test

import (
	"context"
	"testing"

	"github.com/korolev-n/merch/internal/domain"
	"github.com/korolev-n/merch/internal/repository"
	"github.com/korolev-n/merch/internal/service"
)

type shopRepoMock struct {
	repository.ShopRepository
	GetItemByTypeFunc       func(ctx context.Context, itemType string) (*domain.InventoryItem, error)
	GetWalletByUserIDFunc   func(ctx context.Context, userID int) (*domain.Wallet, error)
	UpdateWalletBalanceFunc func(ctx context.Context, userID, newBalance int) error
	AddToUserInventoryFunc  func(ctx context.Context, userID, inventoryID int) error
}

func (m *shopRepoMock) GetItemByType(ctx context.Context, itemType string) (*domain.InventoryItem, error) {
	return m.GetItemByTypeFunc(ctx, itemType)
}

func (m *shopRepoMock) GetWalletByUserID(ctx context.Context, userID int) (*domain.Wallet, error) {
	return m.GetWalletByUserIDFunc(ctx, userID)
}

func (m *shopRepoMock) UpdateWalletBalance(ctx context.Context, userID, newBalance int) error {
	return m.UpdateWalletBalanceFunc(ctx, userID, newBalance)
}

func (m *shopRepoMock) AddToUserInventory(ctx context.Context, userID, inventoryID int) error {
	return m.AddToUserInventoryFunc(ctx, userID, inventoryID)
}

func TestBuyItem_Success(t *testing.T) {
	mock := &shopRepoMock{
		GetItemByTypeFunc: func(ctx context.Context, itemType string) (*domain.InventoryItem, error) {
			return &domain.InventoryItem{ID: 1, Price: 100}, nil
		},
		GetWalletByUserIDFunc: func(ctx context.Context, userID int) (*domain.Wallet, error) {
			return &domain.Wallet{Balance: 200}, nil
		},
		UpdateWalletBalanceFunc: func(ctx context.Context, userID, newBalance int) error {
			return nil
		},
		AddToUserInventoryFunc: func(ctx context.Context, userID, inventoryID int) error {
			return nil
		},
	}

	s := service.NewShopService(mock)
	err := s.BuyItem(context.Background(), 1, "cup")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

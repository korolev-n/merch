package mocks

import (
	"context"

	"github.com/korolev-n/merch/internal/domain"
	"github.com/korolev-n/merch/internal/repository"
)

type WalletRepositoryMock struct {
	CreateFunc          func(ctx context.Context, wallet *domain.Wallet) error
	TransferCoinsTxFunc func(ctx context.Context, fromUserID, toUserID, amount int) error
}

func NewMockWalletRepository() repository.WalletRepository {
	return &WalletRepositoryMock{}
}

func (m *WalletRepositoryMock) Create(ctx context.Context, wallet *domain.Wallet) error {
	if m.CreateFunc != nil {
		return m.CreateFunc(ctx, wallet)
	}
	return nil
}

func (m *WalletRepositoryMock) TransferCoinsTx(ctx context.Context, fromUserID, toUserID, amount int) error {
	if m.TransferCoinsTxFunc != nil {
		return m.TransferCoinsTxFunc(ctx, fromUserID, toUserID, amount)
	}
	return nil
}

package mocks

import (
	"context"

	"github.com/korolev-n/merch-auth/internal/domain"
	"github.com/korolev-n/merch-auth/internal/repository"
)

type WalletRepositoryMock struct {
	CreateFunc func(ctx context.Context, wallet *domain.Wallet) error
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

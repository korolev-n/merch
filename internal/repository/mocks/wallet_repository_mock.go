package mocks

import "github.com/korolev-n/merch-auth/internal/repository"

type WalletRepositoryMock struct {
	CreateFunc func() error
}

func NewMockWalletRepository() repository.WalletRepository {
	return &WalletRepositoryMock{}
}

func (m *WalletRepositoryMock) Create() error {
	return m.CreateFunc()
}

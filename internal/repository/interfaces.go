package repository

import (
	"context"

	"github.com/korolev-n/merch-auth/internal/domain"
)

type UserRepository interface {
	GetByUsername(ctx context.Context, username string) (*domain.User, error)
	Create(ctx context.Context, user *domain.User) (int, error)
}

type WalletRepository interface {
	Create(ctx context.Context, wallet *domain.Wallet) error
}

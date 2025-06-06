package repository

import (
	"context"
	"database/sql"

	"github.com/korolev-n/merch-auth/internal/domain"
)

type WalletRepository interface {
	Create(ctx context.Context, wallet *domain.Wallet) error
}

type walletRepo struct {
	db *sql.DB
}

func NewWalletRepository(db *sql.DB) WalletRepository {
	return &walletRepo{db: db}
}

func (r *walletRepo) Create(ctx context.Context, wallet *domain.Wallet) error {
	_, err := r.db.ExecContext(ctx, "INSERT INTO wallets (user_id, balance) VALUES ($1, $2)", wallet.UserID, wallet.Balance)
	return err
}

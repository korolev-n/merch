package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/korolev-n/merch-auth/internal/domain"
)

type WalletRepository interface {
	Create(ctx context.Context, wallet *domain.Wallet) error
	TransferCoinsTx(ctx context.Context, fromUserID, toUserID, amount int) error
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

func (r *walletRepo) TransferCoinsTx(ctx context.Context, fromUserID, toUserID, amount int) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Проверка баланса
	var balance int
	err = tx.QueryRowContext(ctx, "SELECT balance FROM wallets WHERE user_id = $1 FOR UPDATE", fromUserID).Scan(&balance)
	if err != nil {
		return err
	}
	if balance < amount {
		return errors.New("insufficient balance")
	}

	// Списываем
	_, err = tx.ExecContext(ctx, "UPDATE wallets SET balance = balance - $1 WHERE user_id = $2", amount, fromUserID)
	if err != nil {
		return err
	}

	// Зачисляем
	_, err = tx.ExecContext(ctx, "UPDATE wallets SET balance = balance + $1 WHERE user_id = $2", amount, toUserID)
	if err != nil {
		return err
	}

	// Добавляем запись о транзакции
	_, err = tx.ExecContext(ctx, `
		INSERT INTO coin_transactions (from_user_id, to_user_id, amount)
		VALUES ($1, $2, $3)
	`, fromUserID, toUserID, amount)
	if err != nil {
		return err
	}

	return tx.Commit()
}

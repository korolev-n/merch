package repository

import (
	"context"
	"database/sql"

	"github.com/korolev-n/merch-auth/internal/domain"
)

type InfoRepository interface {
	GetWalletBalance(ctx context.Context, userID int) (int, error)
	GetUserInventory(ctx context.Context, userID int) ([]domain.InventoryInfo, error)
	GetReceivedCoins(ctx context.Context, userID int) ([]domain.CoinHistoryEntry, error)
	GetSentCoins(ctx context.Context, userID int) ([]domain.CoinSentEntry, error)
}

type infoRepo struct {
	db *sql.DB
}

func NewInfoRepository(db *sql.DB) InfoRepository {
	return &infoRepo{db: db}
}

func (r *infoRepo) GetWalletBalance(ctx context.Context, userID int) (int, error) {
	var balance int
	err := r.db.QueryRowContext(ctx, "SELECT balance FROM wallets WHERE user_id = $1", userID).Scan(&balance)
	return balance, err
}

func (r *infoRepo) GetUserInventory(ctx context.Context, userID int) ([]domain.InventoryInfo, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT i.type, ui.quantity
		FROM users_inventory ui
		JOIN inventory i ON ui.inventory_id = i.id
		WHERE ui.user_id = $1
	`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []domain.InventoryInfo
	for rows.Next() {
		var item domain.InventoryInfo
		if err := rows.Scan(&item.Type, &item.Quantity); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, nil
}

func (r *infoRepo) GetReceivedCoins(ctx context.Context, userID int) ([]domain.CoinHistoryEntry, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT u.username, ct.amount
		FROM coin_transactions ct
		JOIN users u ON ct.from_user_id = u.id
		WHERE ct.to_user_id = $1
	`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var history []domain.CoinHistoryEntry
	for rows.Next() {
		var entry domain.CoinHistoryEntry
		if err := rows.Scan(&entry.FromUser, &entry.Amount); err != nil {
			return nil, err
		}
		history = append(history, entry)
	}
	return history, nil
}

func (r *infoRepo) GetSentCoins(ctx context.Context, userID int) ([]domain.CoinSentEntry, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT u.username, ct.amount
		FROM coin_transactions ct
		JOIN users u ON ct.to_user_id = u.id
		WHERE ct.from_user_id = $1
	`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var history []domain.CoinSentEntry
	for rows.Next() {
		var entry domain.CoinSentEntry
		if err := rows.Scan(&entry.ToUser, &entry.Amount); err != nil {
			return nil, err
		}
		history = append(history, entry)
	}
	return history, nil
}

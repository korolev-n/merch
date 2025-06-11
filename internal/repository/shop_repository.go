package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/korolev-n/merch-auth/internal/domain"
)

type ShopRepository interface {
	GetItemByType(ctx context.Context, itemType string) (*domain.InventoryItem, error)
	GetWalletByUserID(ctx context.Context, userID int) (*domain.Wallet, error)
	UpdateWalletBalance(ctx context.Context, userID, newBalance int) error
	AddToUserInventory(ctx context.Context, userID, inventoryID int) error
}

type shopRepo struct {
	db *sql.DB
}

func NewShopRepository(db *sql.DB) ShopRepository {
	return &shopRepo{db: db}
}

func (r *shopRepo) GetItemByType(ctx context.Context, itemType string) (*domain.InventoryItem, error) {
	row := r.db.QueryRowContext(ctx, "SELECT id, type, price FROM inventory WHERE type = $1", itemType)

	var item domain.InventoryItem
	if err := row.Scan(&item.ID, &item.Type, &item.Price); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &item, nil
}

func (r *shopRepo) GetWalletByUserID(ctx context.Context, userID int) (*domain.Wallet, error) {
	row := r.db.QueryRowContext(ctx, "SELECT id, user_id, balance FROM wallets WHERE user_id = $1", userID)

	var wallet domain.Wallet
	if err := row.Scan(&wallet.ID, &wallet.UserID, &wallet.Balance); err != nil {
		return nil, err
	}
	return &wallet, nil
}

func (r *shopRepo) UpdateWalletBalance(ctx context.Context, userID, newBalance int) error {
	_, err := r.db.ExecContext(ctx, "UPDATE wallets SET balance = $1 WHERE user_id = $2", newBalance, userID)
	return err
}

func (r *shopRepo) AddToUserInventory(ctx context.Context, userID, inventoryID int) error {
	_, err := r.db.ExecContext(ctx, `
		INSERT INTO users_inventory (user_id, inventory_id, quantity)
		VALUES ($1, $2, 1)
		ON CONFLICT (user_id, inventory_id)
		DO UPDATE SET quantity = users_inventory.quantity + 1
	`, userID, inventoryID)
	return err
}

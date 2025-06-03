package repository

import "database/sql"

type walletRepo struct {
	db *sql.DB
}

func NewWalletRepository(db *sql.DB) WalletRepository {
	return &walletRepo{db: db}
}

func (r *walletRepo) Create() error {
	return nil
}

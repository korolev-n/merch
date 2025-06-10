package service

import (
	"context"

	"github.com/korolev-n/merch-auth/internal/repository"
)

type Transfer interface {
	SendCoins(ctx context.Context, fromUserID int, toUsername string, amount int) error
}

type TransferService struct {
	Users   repository.UserRepository
	Wallets repository.WalletRepository
}

func NewTransferService(users repository.UserRepository, wallets repository.WalletRepository) *TransferService {
	return &TransferService{
		Users:   users,
		Wallets: wallets,
	}
}

func (s *TransferService) SendCoins(ctx context.Context, fromUserID int, toUsername string, amount int) error {
	if amount <= 0 {
		return ErrNegativeBalance
	}

	toUser, err := s.Users.GetByUsername(ctx, toUsername)
	if err != nil || toUser == nil {
		return ErrUserNotFound
	}

	return s.Wallets.TransferCoinsTx(ctx, fromUserID, toUser.ID, amount)
}

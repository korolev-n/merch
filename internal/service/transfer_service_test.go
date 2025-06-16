package service_test

import (
	"context"
	"errors"
	"testing"

	"github.com/korolev-n/merch/internal/domain"
	"github.com/korolev-n/merch/internal/repository/mocks"
	"github.com/korolev-n/merch/internal/service"
)

func TestSendCoins_Success(t *testing.T) {
	userRepo := mocks.NewMockUserRepository()
	walletRepo := mocks.NewMockWalletRepository()

	userRepo.GetByUsernameFunc = func(ctx context.Context, username string) (*domain.User, error) {
		return &domain.User{ID: 2}, nil
	}
	walletRepo.TransferCoinsTxFunc = func(ctx context.Context, from, to, amount int) error {
		if from == 1 && to == 2 && amount == 100 {
			return nil
		}
		return errors.New("wrong data")
	}

	s := service.NewTransferService(userRepo, walletRepo)
	err := s.SendCoins(context.Background(), 1, "receiver", 100)
	if err != nil {
		t.Fatalf("expected success, got: %v", err)
	}
}

func TestSendCoins_NegativeAmount(t *testing.T) {
	s := service.NewTransferService(nil, nil)
	err := s.SendCoins(context.Background(), 1, "someone", -1)
	if err != service.ErrNegativeBalance {
		t.Fatalf("expected ErrNegativeBalance, got: %v", err)
	}
}

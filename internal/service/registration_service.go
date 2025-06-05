package service

import (
	"context"

	"github.com/korolev-n/merch-auth/internal/domain"
	"github.com/korolev-n/merch-auth/internal/repository"
)

type RegistrationService struct {
	Users   repository.UserRepository
	Wallets repository.WalletRepository
}

func NewRegistrationService(users repository.UserRepository, wallets repository.WalletRepository) *RegistrationService {
	return &RegistrationService{
		Users:   users,
		Wallets: wallets,
	}
}

func (s *RegistrationService) RegisterUser(ctx context.Context, username, password string) error {
	user, err := s.Users.GetByUsername(ctx, username)
	if err != nil {
		return err
	}

	if user != nil {
		// пользователь уже существует
		return nil
	}

	newUser := &domain.User{
		Username: username,
		Password: password, // хешировать перед сохранением
	}

	userID, err := s.Users.Create(ctx, newUser)
	if err != nil {
		return err
	}

	wallet := &domain.Wallet{
		UserID:  userID,
		Balance: 1000,
	}

	return s.Wallets.Create(ctx, wallet)
}

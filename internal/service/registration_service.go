package service

import (
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

func (s *RegistrationService) RegisterUser() error {
	// if err := s.Users.Create(); err != nil {
	// 	return err
	// }
	// if err := s.Wallets.Create(); err != nil {
	// 	return err
	// }
	return nil
}

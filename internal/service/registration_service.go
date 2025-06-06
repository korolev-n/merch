package service

import (
	"context"

	"github.com/korolev-n/merch-auth/internal/domain"
	"github.com/korolev-n/merch-auth/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

type Registration interface {
	RegisterUser(ctx context.Context, username, password string) (string, error)
}

type RegistrationService struct {
	Users   repository.UserRepository
	Wallets repository.WalletRepository
	JWT     TokenGenerator
}

func NewRegistrationService(users repository.UserRepository, wallets repository.WalletRepository, jwt TokenGenerator) *RegistrationService {
	return &RegistrationService{
		Users:   users,
		Wallets: wallets,
		JWT:     jwt,
	}
}

func (s *RegistrationService) RegisterUser(ctx context.Context, username, password string) (string, error) {
	user, err := s.Users.GetByUsername(ctx, username)
	if err != nil {
		return "", err
	}

	if user != nil {
		err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
		if err != nil {
			return "", ErrIncorrectPassword
		}
		return s.JWT.GenerateToken(user.ID, user.Username)
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	newUser := &domain.User{
		Username: username,
		Password: string(hashedPassword),
	}

	userID, err := s.Users.Create(ctx, newUser)
	if err != nil {
		return "", ErrUserAlreadyExists
	}

	wallet := &domain.Wallet{
		UserID:  userID,
		Balance: 1000,
	}

	if err := s.Wallets.Create(ctx, wallet); err != nil {
		return "", err
	}

	return s.JWT.GenerateToken(userID, username)
}

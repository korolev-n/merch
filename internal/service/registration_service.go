package service

import (
	"context"

	"github.com/korolev-n/merch/internal/domain"
	"github.com/korolev-n/merch/internal/logger"
	"github.com/korolev-n/merch/internal/repository"
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
		logger.Log.Error("error getting user by username", "error", err, "username", username)
		return "", err
	}

	if user != nil {
		err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
		if err != nil {
			logger.Log.Warn("Incorrect password attempt", "username", username)
			return "", ErrIncorrectPassword
		}
		return s.JWT.GenerateToken(user.ID, user.Username)
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		logger.Log.Error("password hashing failed", "error", err)
		return "", err
	}

	newUser := &domain.User{
		Username: username,
		Password: string(hashedPassword),
	}

	userID, err := s.Users.Create(ctx, newUser)
	if err != nil {
		logger.Log.Warn("User creation failed", "username", username, "error", err)
		return "", ErrUserAlreadyExists
	}

	wallet := &domain.Wallet{
		UserID:  userID,
		Balance: 1000,
	}

	if err := s.Wallets.Create(ctx, wallet); err != nil {
		logger.Log.Error("Wallet creation failed", "userID", userID, "error", err)
		return "", err
	}

	return s.JWT.GenerateToken(userID, username)
}

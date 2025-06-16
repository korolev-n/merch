package service_test

import (
	"context"
	"os"
	"testing"

	"github.com/korolev-n/merch/internal/domain"
	"github.com/korolev-n/merch/internal/logger"
	"github.com/korolev-n/merch/internal/repository/mocks"
	"github.com/korolev-n/merch/internal/service"
)

type mockJWT struct{}

func (m *mockJWT) GenerateToken(userID int, username string) (string, error) {
	return "mock-token", nil
}

func (m *mockJWT) ParseToken(token string) (*service.Claims, error) {
	return nil, nil
}

func TestMain(m *testing.M) {
	logger.Init()
	os.Exit(m.Run())
}

func TestRegisterUser_NewUser(t *testing.T) {
	userRepo := mocks.NewMockUserRepository()
	walletRepo := mocks.NewMockWalletRepository()

	userRepo.GetByUsernameFunc = func(ctx context.Context, username string) (*domain.User, error) {
		return nil, nil
	}
	userRepo.CreateFunc = func(ctx context.Context, user *domain.User) (int, error) {
		return 42, nil
	}
	walletRepo.CreateFunc = func(ctx context.Context, wallet *domain.Wallet) error {
		return nil
	}

	s := service.NewRegistrationService(userRepo, walletRepo, &mockJWT{})
	token, err := s.RegisterUser(context.Background(), "new_user", "123456")
	if err != nil || token != "mock-token" {
		t.Fatalf("expected token 'mock-token', got '%s', error: %v", token, err)
	}
}

func TestRegisterUser_ExistingUser_WrongPassword(t *testing.T) {
	userRepo := mocks.NewMockUserRepository()
	walletRepo := mocks.NewMockWalletRepository()

	hashed := "$2a$10$F9NOn4K6Iq4ZrPZmdQOVMOLQyl7bEj2BjxtLEb8YHHTMnEvYJ8Mf6" // hash of "correct"

	userRepo.GetByUsernameFunc = func(ctx context.Context, username string) (*domain.User, error) {
		return &domain.User{ID: 1, Username: username, Password: hashed}, nil
	}

	s := service.NewRegistrationService(userRepo, walletRepo, &mockJWT{})
	_, err := s.RegisterUser(context.Background(), "existing_user", "wrong")
	if err != service.ErrIncorrectPassword {
		t.Fatalf("expected ErrIncorrectPassword, got %v", err)
	}
}

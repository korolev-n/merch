package mocks

import (
	"context"

	"github.com/korolev-n/merch-auth/internal/domain"
	"github.com/korolev-n/merch-auth/internal/repository"
)

type UserRepositoryMock struct {
	CreateFunc        func(ctx context.Context, user *domain.User) (int, error)
	GetByUsernameFunc func(ctx context.Context, username string) (*domain.User, error)
}

func NewMockUserRepository() repository.UserRepository {
	return &UserRepositoryMock{}
}

func (m *UserRepositoryMock) Create(ctx context.Context, user *domain.User) (int, error) {
	if m.CreateFunc != nil {
		return m.CreateFunc(ctx, user)
	}
	return 0, nil
}

func (m *UserRepositoryMock) GetByUsername(ctx context.Context, username string) (*domain.User, error) {
	if m.GetByUsernameFunc != nil {
		return m.GetByUsernameFunc(ctx, username)
	}
	return nil, nil
}

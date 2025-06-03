package mocks

import "github.com/korolev-n/merch-auth/internal/repository"

type UserRepositoryMock struct {
	CreateFunc func() error
}

func NewMockUserRepository() repository.UserRepository {
	return &UserRepositoryMock{}
}

func (m *UserRepositoryMock) Create() error {
	return m.CreateFunc()
}

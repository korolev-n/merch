package repository

import (
	"context"
	"database/sql"

	"github.com/korolev-n/merch/internal/domain"
)

type UserRepository interface {
	GetByUsername(ctx context.Context, username string) (*domain.User, error)
	Create(ctx context.Context, user *domain.User) (int, error)
}

type userRepo struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) UserRepository {
	return &userRepo{db: db}
}

func (r *userRepo) GetByUsername(ctx context.Context, username string) (*domain.User, error) {
	row := r.db.QueryRowContext(ctx, "SELECT id, username, password_hash FROM users WHERE username = $1", username)

	var user domain.User
	if err := row.Scan(&user.ID, &user.Username, &user.Password); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (r *userRepo) Create(ctx context.Context, user *domain.User) (int, error) {
	var id int
	err := r.db.QueryRowContext(ctx, "INSERT INTO users (username, password_hash) VALUES ($1, $2) RETURNING id", user.Username, user.Password).Scan(&id)
	return id, err
}

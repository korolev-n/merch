package service

import "errors"

var (
	ErrIncorrectPassword   = errors.New("incorrect password")
	ErrUserAlreadyExists   = errors.New("user already exists")
	ErrUserNotFound        = errors.New("user not found")
	ErrTokenGeneration     = errors.New("could not generate token")
	ErrInvalidToken        = errors.New("invalid JWT token")
	ErrInsufficientBalance = errors.New("insufficient balance")
	ErrNegativeBalance     = errors.New("amount must be greater than zero")
	ErrItemNotFound        = errors.New("item not found")
)

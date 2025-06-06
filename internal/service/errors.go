package service

import "errors"

var (
	ErrIncorrectPassword  = errors.New("incorrect password")
	ErrUserAlreadyExists  = errors.New("user already exists")
	ErrUserNotFound       = errors.New("user not found")
	ErrTokenGeneration    = errors.New("could not generate token")
)

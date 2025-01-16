package domain

import (
	"context"
)

type User struct {
	ID           string
	Email        string
	PasswordHash string
	UserType     string // "admin" or "user"
}

type AuthRepository interface {
	GetUserByEmail(ctx context.Context, email string) (*User, error)
	InitializeAdmin(ctx context.Context, email, passwordHash string) error
	BlacklistToken(ctx context.Context, token string) error
	IsTokenBlacklisted(ctx context.Context, token string) (bool, error)
}

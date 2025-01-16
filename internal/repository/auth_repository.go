package repository

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/mohamedfawas/rmshop-auth-service/internal/domain"
)

type authRepository struct {
	db *sql.DB
}

func NewAuthRepository(db *sql.DB) domain.AuthRepository {
	return &authRepository{db: db}
}

func (r *authRepository) GetUserByEmail(ctx context.Context, email string) (*domain.User, error) {
	// Try admin first
	var user domain.User
	err := r.db.QueryRowContext(ctx,
		"SELECT id, email, password_hash FROM admins WHERE email = $1",
		email,
	).Scan(&user.ID, &user.Email, &user.PasswordHash)

	if err == nil {
		user.UserType = "admin"
		return &user, nil
	}

	if err != sql.ErrNoRows {
		return nil, err
	}

	// Try regular user
	err = r.db.QueryRowContext(ctx,
		"SELECT id, email, password_hash FROM users WHERE email = $1",
		email,
	).Scan(&user.ID, &user.Email, &user.PasswordHash)

	if err != nil {
		return nil, err
	}

	user.UserType = "user"
	return &user, nil
}

func (r *authRepository) InitializeAdmin(ctx context.Context, email, passwordHash string) error {
	query := `
        INSERT INTO admins (email, password_hash) 
        VALUES ($1, $2) 
        ON CONFLICT (email) DO UPDATE 
        SET password_hash = EXCLUDED.password_hash
    `
	_, err := r.db.ExecContext(ctx, query, email, passwordHash)
	if err != nil {
		return fmt.Errorf("failed to initialize admin: %v", err)
	}
	return nil
}

func (r *authRepository) BlacklistToken(ctx context.Context, token string) error {
	// Set expiration time
	expiresAt := time.Now().UTC()

	query := `INSERT INTO blacklisted_tokens (token, expires_at) VALUES ($1, $2)`
	_, err := r.db.ExecContext(ctx, query, token, expiresAt)
	if err != nil {
		log.Printf("Error blacklisting token: %v", err) // Log the error
	}
	return err
}

func (r *authRepository) IsTokenBlacklisted(ctx context.Context, token string) (bool, error) {
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM blacklisted_tokens WHERE token = $1)`
	err := r.db.QueryRowContext(ctx, query, token).Scan(&exists)
	return exists, err
}

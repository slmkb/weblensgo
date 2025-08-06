package models

import (
	"crypto/sha256"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/slmkb/weblensgo/rand"
)

const (
	DefaultTokenDuration = 60 * time.Minute
)

type PasswordReset struct {
	UserID     uuid.UUID
	Email      string
	RToken     string
	RTokenHash string
	CreatedAt  time.Time
}

type PasswordResetService struct {
	DB            *sql.DB
	TokenDuration time.Duration
}

func (rs *PasswordResetService) Create(email string) (*PasswordReset, error) {
	token, err := rand.ResetToken()
	if err != nil {
		return nil, fmt.Errorf("password reset service create: %w", err)
	}
	tokenHash := rs.hash(token)

	pr := PasswordReset{
		Email:      email,
		RToken:     token,
		RTokenHash: tokenHash,
		CreatedAt:  time.Now(),
	}

	result, err := rs.DB.Exec(`
		WITH u_id AS (
			SELECT id
			FROM users
			WHERE email = $1
		)
		INSERT INTO reset_tokens (user_id, token_hash, created_at)
		SELECT id, $2, $3
		FROM u_id
		ON CONFLICT (user_id) DO
		UPDATE
			SET 
			token_hash = EXCLUDED.token_hash,
			created_at = EXCLUDED.created_at
	`, pr.Email, pr.RTokenHash, pr.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("password reset service create: %w", err)
	}
	if rows, err := result.RowsAffected(); err != nil {
		return nil, fmt.Errorf("password reset service create: %w", err)
	} else if rows != 1 {
		return nil, fmt.Errorf("password reset service create: user %s not found", email)
	}
	return &pr, nil
}

func (rs *PasswordResetService) Consume(token string) (uuid.UUID, error) {
	tokenHash := rs.hash(token)

	row := rs.DB.QueryRow(`
		DELETE FROM reset_tokens
		WHERE token_hash = $1
		RETURNING created_at, user_id;
	`, tokenHash)

	var result struct {
		UserID    uuid.UUID
		CreatedAt time.Time
	}
	if err := row.Scan(&result.CreatedAt, &result.UserID); err != nil {
		return uuid.Nil, fmt.Errorf("password reset service consume: %w", err)
	}
	if time.Now().After(result.CreatedAt.Add(DefaultTokenDuration)) {
		return uuid.Nil, errors.New("password reset service consume: token expired")
	}
	return result.UserID, nil
}

func (rs *PasswordResetService) hash(token string) string {
	return fmt.Sprintf("%x", sha256.Sum256([]byte(token)))
}

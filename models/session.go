package models

import (
	"crypto/sha256"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
	"github.com/slmkb/weblensgo/rand"
)

type Session struct {
	ID        uuid.UUID
	UserID    uuid.UUID
	TokenHash string
	Token     string //used only during session creation to create a cookie
}

type SessionService struct {
	DB *sql.DB
}

func (ss *SessionService) Create(u *User) (*Session, error) {
	token, err := rand.SessionToken()
	if err != nil {
		return nil, fmt.Errorf("session service create: %w", err)
	}
	tokenHash := ss.hash(token)

	s := Session{
		ID:        uuid.New(),
		UserID:    u.ID,
		TokenHash: tokenHash,
		Token:     token,
	}

	_, err = ss.DB.Exec(`
		INSERT INTO sessions (id, user_id, token_hash)
		VALUES($1, $2, $3) ON CONFLICT (user_id) DO
		UPDATE
		SET token_hash = EXCLUDED.token_hash
	`, s.ID, s.UserID, s.TokenHash)
	if err != nil {
		return nil, fmt.Errorf("sessions service create: %w", err)
	}

	return &s, nil
}

func (ss *SessionService) GetUser(token string) (*User, error) {
	tokenHash := ss.hash(token)

	row := ss.DB.QueryRow(`
		SELECT sessions.user_id,
			users.email,
			users.password_hash
		FROM sessions
		INNER JOIN users ON sessions.user_id = users.id
		WHERE token_hash = $1;
	`, tokenHash)

	var user User
	if err := row.Scan(&user.ID, &user.Email, &user.PasswordHash); err != nil {
		return nil, fmt.Errorf("session service getuser: %w", err)
	}
	return &user, nil
}

func (ss *SessionService) DeleteSession(token string) error {
	tokenHash := ss.hash(token)
	_, err := ss.DB.Exec(`
		DELETE
		FROM sessions
		WHERE token_hash = $1;
	`, tokenHash)
	if err != nil {
		return fmt.Errorf("delete session: %w", err)
	}
	return nil
}

func (ss *SessionService) hash(token string) string {
	return fmt.Sprintf("%x", sha256.Sum256([]byte(token)))
}

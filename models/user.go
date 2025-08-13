package models

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"

	uuid "github.com/google/uuid"
	"github.com/jackc/pgx"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID           uuid.UUID
	Email        string
	PasswordHash string
}

type UserService struct {
	DB *sql.DB
}

var (
	ErrEmailTaken    = errors.New("models: email already taken")
	ErrUUIDCollision = errors.New("models: user uuid collision")
)

func (us *UserService) Create(email, password string) (*User, error) {
	email = strings.ToLower(email)
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("create user: %w", err)
	}

	user := User{
		ID:           uuid.New(),
		Email:        email,
		PasswordHash: string(hashedBytes),
	}

	_, err = us.DB.Exec(`
			INSERT INTO users (id, email, password_hash)
			VALUES ($1, $2, $3)`, user.ID, user.Email, user.PasswordHash)

	var pgErr pgx.PgError
	if errors.As(err, &pgErr) {
		switch pgErr.ConstraintName {
		case "users_pkey":
			return nil, ErrUUIDCollision
		case "users_email_key":
			return nil, ErrEmailTaken
		default:
			return nil, err
		}
	}

	return &user, nil
}

func (us *UserService) Authenticate(email, password string) (*User, error) {
	email = strings.ToLower(email)

	row := us.DB.QueryRow(`
	SELECT * FROM users WHERE
	email = $1`, email)

	var user User
	if err := row.Scan(&user.ID, &user.Email, &user.PasswordHash); err != nil {
		return nil, fmt.Errorf("get user: %w", err)
	}

	err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return nil, fmt.Errorf("get user: %w", err)
	}

	return &user, nil

}

func (us *UserService) UpdatePassword(userID uuid.UUID, password string) error {
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("update password: %w", err)
	}

	_, err = us.DB.Exec(`
	UPDATE users
	SET password_hash = $1
	WHERE id = $2
	`, hashedBytes, userID)
	if err != nil {
		return fmt.Errorf("update password: %w", err)
	}
	// if rows, err := res.RowsAffected(); err != nil {
	// 	return fmt.Errorf("update password: %w", err)
	// } else if rows != 1 {
	// 	return fmt.Errorf("update password: userID %s not found", userID)
	// }
	return nil
}

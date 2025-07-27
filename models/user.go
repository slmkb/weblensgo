package models

import (
	"database/sql"
	"fmt"
	"strings"

	uuid "github.com/google/uuid"
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

	row := us.DB.QueryRow(`
	INSERT INTO users (id, email, password_hash)
	VALUES ($1, $2, $3) RETURNING id`, user.ID, user.Email, user.PasswordHash)

	var confirmUUID uuid.UUID
	if err := row.Scan(&confirmUUID); err != nil {
		return nil, fmt.Errorf("create user: %w", err)
	}

	if confirmUUID != user.ID {
		return nil, fmt.Errorf("create user: %w", err)
	}

	return &user, nil
}

func (us *UserService) GetUser(email, password string) (*User, error) {
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

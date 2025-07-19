package main

import (
	"errors"
	"fmt"
)

func Connect() error {
	return errors.New("connection failed")
}

func CreateUser() error {
	if err := Connect(); err != nil {
		return fmt.Errorf("create user: %w", err)
	}
	return nil
}

func CreateOrg() error {
	if err := CreateUser(); err != nil {
		return fmt.Errorf("create org: %w", err)
	}
	return nil
}

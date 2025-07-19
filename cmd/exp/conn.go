package main

import "errors"

func Connect() error {
	return errors.New("connection failed")
}

func CreateUser() error {
	if err := Connect(); err != nil {
		return err
	}
	return nil
}

func CreateOrg() error {
	if err := CreateUser(); err != nil {
		return err
	}
	return nil
}

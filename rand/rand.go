package rand

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
)

const (
	sessionTokenBytes = 32
	resetTokenBytes   = 64
	galleryIDBytes    = 8
)

func SessionToken() (string, error) {
	b := make([]byte, sessionTokenBytes)
	if _, err := rand.Read(b); err != nil {
		return "", fmt.Errorf("session token: %w", err)
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

func ResetToken() (string, error) {
	b := make([]byte, resetTokenBytes)
	if _, err := rand.Read(b); err != nil {
		return "", fmt.Errorf("reset token: %w", err)
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

func GalleryHash() (string, error) {
	b := make([]byte, galleryIDBytes)
	if _, err := rand.Read(b); err != nil {
		return "", fmt.Errorf("gallery hash: %w", err)
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}

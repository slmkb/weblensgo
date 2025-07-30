package rand

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
)

const sessionTokenBytes = 32

func SessionToken() (string, error) {
	b := make([]byte, sessionTokenBytes)
	if _, err := rand.Read(b); err != nil {
		return "", fmt.Errorf("session token: %w", err)
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

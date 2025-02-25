package encryption

import (
	"crypto/rand"
	"encoding/base64"
)

// Generate random base64 string
func GenerateRandomBase64String(bits int) (string, error) {
	randomBytes := make([]byte, bits)
	_, err := rand.Read(randomBytes)
	if err != nil {
		return "", err
	}

	state := base64.URLEncoding.EncodeToString(randomBytes)
	return state, nil
}

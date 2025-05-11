package util

import (
	"crypto/rand"
	"encoding/hex"
)

// GenerateSecureID creates a cryptographically secure random ID.
func GenerateSecureID() (string, error) {
	// Create a 32-byte random value (256-bit)
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}

	// Convert to a hex string
	return hex.EncodeToString(b), nil
}

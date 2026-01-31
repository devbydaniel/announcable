package random

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base32"
	"encoding/hex"
)

// CreateRandomToken generates a cryptographically random base32-encoded token.
func CreateRandomToken() string {
	bytes := make([]byte, 15)
	if _, err := rand.Read(bytes); err != nil {
		panic("failed to read random bytes: " + err.Error())
	}
	token := base32.StdEncoding.EncodeToString(bytes)
	return token
}

// EncodeToken returns the SHA-256 hex digest of the given token.
func EncodeToken(token string) string {
	byteID := sha256.Sum256([]byte(token))
	return hex.EncodeToString(byteID[:])
}

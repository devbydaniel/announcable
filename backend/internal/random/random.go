package random

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base32"
	"encoding/hex"
)

func CreateRandomToken() string {
	bytes := make([]byte, 15)
	rand.Read(bytes)
	token := base32.StdEncoding.EncodeToString(bytes)
	return token
}

func EncodeToken(token string) string {
	byteId := sha256.Sum256([]byte(token))
	return hex.EncodeToString(byteId[:])
}

package util

import (
	"crypto/sha256"
	"encoding/hex"
)

func HashPassword(username, password string) string {

	saltPass := append([]byte(password), []byte(username)...)
	hash := sha256.Sum256(saltPass)
	return hex.EncodeToString(hash[:])
}

func CheckPasswordHash(username, password, storedHash string) bool {
	// Generate hash using the same method as HashPassword
	generatedHash := HashPassword(username, password)

	// Compare the generated hash with stored hash
	return generatedHash == storedHash
}

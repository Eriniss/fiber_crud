package utils

import (
	"crypto/sha256"
	"encoding/hex"
)

// SHA-256 해싱 함수
func HashPassword(password string) string {
	hash := sha256.New()
	hash.Write([]byte(password))
	return hex.EncodeToString(hash.Sum(nil))
}

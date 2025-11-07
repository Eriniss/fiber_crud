package utils

import (
	"golang.org/x/crypto/bcrypt"
)

// bcrypt를 사용한 비밀번호 해싱 함수
func HashPassword(password string) (string, error) {
	// bcrypt.DefaultCost (10)를 사용하여 해싱
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedBytes), nil
}

// 비밀번호 검증 함수
func VerifyPassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

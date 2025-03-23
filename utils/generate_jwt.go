package utils

import (
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

// JWT 토큰 생성
func GenerateJWT(userId uint, email string, role uint) (string, error) {
	// 비밀 키 가져오기
	secretKey := os.Getenv("JWT_SECRET_KEY")

	if secretKey == "" {
		return "", errors.New("secret key not found")
	}

	// JWT claims (payload)
	claims := jwt.MapClaims{
		"id":    userId,
		"email": email,
		"role":  role,
		"exp":   time.Now().Add(24 * time.Hour).Unix(), // 24시간 후 만료 설정
	}

	// JWT 토큰 생성
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// 서명하여 토큰 생성
	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

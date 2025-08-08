package utils

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

// JWT 토큰 생성
func GenerateJWT(userId uint, email string, name string, point uint) (string, error) {
	// 비밀 키 가져오기
	secretKey := "mysecretkey"

	if secretKey == "" {
		return "", errors.New("secret key not found")
	}

	// JWT claims (payload)
	claims := jwt.MapClaims{
		"id":    userId,
		"email": email,
		"name":  name,
		"point": point,
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

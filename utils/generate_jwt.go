package utils

import (
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

// JWT 토큰 생성
func GenerateJWT(userId uint, email string, name string, point uint) (string, error) {
	// 환경변수에서 비밀 키 가져오기
	secretKey := os.Getenv("JWT_SECRET_KEY")

	if secretKey == "" {
		return "", errors.New("JWT_SECRET_KEY not found in environment variables")
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

// JWT 토큰 검증
func VerifyJWT(tokenString string) (*jwt.MapClaims, error) {
	// 환경변수에서 비밀 키 가져오기
	secretKey := os.Getenv("JWT_SECRET_KEY")

	if secretKey == "" {
		return nil, errors.New("JWT_SECRET_KEY not found in environment variables")
	}

	// 토큰 파싱
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// 서명 방법 검증
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(secretKey), nil
	})

	if err != nil {
		return nil, err
	}

	// claims 추출
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return &claims, nil
	}

	return nil, errors.New("invalid token")
}

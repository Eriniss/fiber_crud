package oidc

import (
	"crypto/rand"
	"encoding/base64"
	"os"

	"github.com/gofiber/fiber/v3"
)

// State 토큰 생성 (CSRF 방어)
func generateStateToken() (string, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

// OIDC 로그인 리다이렉트
// Logto 로그인 페이지로 리다이렉트
func LoginRedirect(c fiber.Ctx) error {
	// OIDC가 활성화되어 있는지 확인
	if !IsOIDCEnabled() {
		return c.Status(503).JSON(fiber.Map{
			"error": "OIDC is not configured",
			"note":  "Please configure LOGTO_ENDPOINT, LOGTO_APP_ID, and LOGTO_APP_SECRET in .env",
		})
	}

	// State 토큰 생성
	state, err := generateStateToken()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Failed to generate state token",
		})
	}

	// State를 쿠키에 저장 (CSRF 방어)
	c.Cookie(&fiber.Cookie{
		Name:     "oauth_state",
		Value:    state,
		HTTPOnly: true,
		Secure:   os.Getenv("ENVIRONMENT") == "production",
		SameSite: "Lax",
		MaxAge:   600, // 10분
	})

	// OAuth2 Authorization URL 생성
	authURL := GetOAuth2Config().AuthCodeURL(state)

	// 프론트엔드에 URL 반환
	return c.JSON(fiber.Map{
		"auth_url": authURL,
		"message":  "Redirect user to this URL for OIDC login",
	})
}

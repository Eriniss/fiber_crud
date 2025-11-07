package oidc

import (
	"os"

	"github.com/gofiber/fiber/v3"
)

// OIDC 로그인 리다이렉트
// Logto 로그인 페이지로 리다이렉트
func LoginRedirect(c fiber.Ctx) error {
	// Logto 설정 가져오기
	logtoEndpoint := os.Getenv("LOGTO_ENDPOINT")
	logtoAppId := os.Getenv("LOGTO_APP_ID")
	redirectUri := os.Getenv("LOGTO_REDIRECT_URI")

	if logtoEndpoint == "" || logtoAppId == "" || redirectUri == "" {
		return c.Status(500).JSON(fiber.Map{
			"error": "OIDC configuration is missing",
		})
	}

	// OIDC 인증 URL 생성
	// TODO: 실제 Logto SDK 또는 OIDC 라이브러리 사용 시 구현
	authURL := logtoEndpoint + "/oidc/auth?client_id=" + logtoAppId + "&redirect_uri=" + redirectUri + "&response_type=code&scope=openid profile email"

	return c.JSON(fiber.Map{
		"auth_url": authURL,
		"message":  "Redirect to this URL for OIDC login",
	})
}

package oidc

import (
	"fmt"
	"net/url"
	"os"

	"github.com/gofiber/fiber/v3"
)

// OIDC 로그아웃
func Logout(c fiber.Ctx) error {
	// OIDC가 활성화되어 있는지 확인
	if !IsOIDCEnabled() {
		return c.Status(503).JSON(fiber.Map{
			"error": "OIDC is not configured",
		})
	}

	// Logto 로그아웃 URL 생성
	logtoEndpoint := os.Getenv("LOGTO_ENDPOINT")
	postLogoutRedirect := os.Getenv("LOGTO_POST_LOGOUT_REDIRECT_URI")

	logoutURL := fmt.Sprintf("%s/oidc/session/end?post_logout_redirect_uri=%s",
		logtoEndpoint,
		url.QueryEscape(postLogoutRedirect),
	)

	return c.JSON(fiber.Map{
		"logout_url": logoutURL,
		"message":    "Redirect user to this URL to logout from Logto",
		"note":       "Client should also remove the JWT token from local storage",
	})
}

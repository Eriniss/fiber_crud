package routes

import (
	oidc "fiber_crud/handlers/oidc_handlers"

	"github.com/gofiber/fiber/v3"
)

func OIDCRoutes(app *fiber.App) {
	api := app.Group("/oidc")

	// OIDC 인증 관련 엔드포인트
	api.Get("/login", oidc.LoginRedirect)   // OIDC 로그인 시작
	api.Get("/callback", oidc.Callback)     // OIDC 콜백
	api.Post("/logout", oidc.Logout)        // OIDC 로그아웃
}

package middleware

import (
	"fiber_crud/utils"
	"strings"

	"github.com/gofiber/fiber/v3"
)

// JWT 인증 미들웨어
func AuthRequired() fiber.Handler {
	return func(c fiber.Ctx) error {
		// Authorization 헤더에서 토큰 가져오기
		authHeader := c.Get("Authorization")

		if authHeader == "" {
			return c.Status(401).JSON(fiber.Map{
				"error": "Missing authorization header",
			})
		}

		// "Bearer " 접두사 제거
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			return c.Status(401).JSON(fiber.Map{
				"error": "Invalid authorization format. Use: Bearer <token>",
			})
		}

		// JWT 토큰 검증
		claims, err := utils.VerifyJWT(tokenString)
		if err != nil {
			return c.Status(401).JSON(fiber.Map{
				"error": "Invalid or expired token",
			})
		}

		// claims를 context에 저장 (핸들러에서 사용 가능)
		c.Locals("user", claims)

		// 다음 핸들러로 진행
		return c.Next()
	}
}

package handlers

import (
	"fmt"

	"github.com/gofiber/fiber/v3"
)

// LoggerMiddleware - 모든 요청을 로그로 기록
func LoggerMiddleware(c fiber.Ctx) error {
	fmt.Printf("Request: %s %s\n", c.Method(), c.OriginalURL())
	return c.Next()
}

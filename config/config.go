package config

import (
	"log"

	"github.com/gofiber/fiber/v3"
)

// AppConfig 설정
func AppConfig() *fiber.App {
	app := fiber.New(fiber.Config{
		CaseSensitive: true,
		StrictRouting: false,
	})

	log.Println("Fiber server setting is done")
	return app
}

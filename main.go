package main

import (
	"log"

	"github.com/gofiber/fiber/v3"
)

func main() {
	// Fiber 앱 생성
	app := fiber.New()

	// 서버 실행
	log.Println("🚀 Server's hot in 3000 port!")
	log.Fatal(app.Listen(":3000"))
}

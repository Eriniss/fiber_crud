package main

import (
	"fiber_curd/database"
	"fiber_curd/routes"
	"log"

	"github.com/gofiber/fiber/v3"
)

func main() {
	// 초기 데이터베이스 생성 및 초기화
	database.InitDatabase()

	// Fiber 앱 생성
	app := fiber.New()

	// 라우트 설정
	routes.UserRoutes(app)

	// 서버 실행
	log.Println("🚀 Server's hot in 3000 port!")
	log.Fatal(app.Listen(":3000"))
}

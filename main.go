package main

import (
	"fiber_crud/database"
	"fiber_crud/routes"
	"log"
	"os"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
	"github.com/joho/godotenv"
)

func main() {
	// .env 로드
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	// Fiber 앱 생성
	app := fiber.New()

	// CORS 설정
	app.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"}, // 프론트 주소
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept"},
		AllowCredentials: false, // 세션/쿠키 쓸 경우 true
	}))

	port := os.Getenv("API_PORT")

	// DB 초기화
	database.InitDatabase()

	// 라우트 설정
	routes.UserRoutes(app)
	routes.BlogRoutes(app)

	// 서버 실행
	log.Printf("🚀 Server's hot in %s port!\n", port)
	log.Fatal(app.Listen(":" + port))
}

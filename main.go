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
	allowedOrigins := os.Getenv("ALLOWED_ORIGINS")
	if allowedOrigins == "" {
		allowedOrigins = "http://localhost:3000" // 기본값
	}

	app.Use(cors.New(cors.Config{
		AllowOrigins:     []string{allowedOrigins},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		AllowCredentials: true,
		MaxAge:           300, // Preflight 요청 캐시 시간 (초)
	}))

	port := os.Getenv("API_PORT")

	// DB 초기화
	database.InitDatabase()

	// 라우트 설정
	routes.UserRoutes(app)
	routes.OIDCRoutes(app) // Logto OIDC 인증

	// 서버 실행
	log.Printf("🚀 Server's hot in %s port!\n", port)
	log.Fatal(app.Listen(":" + port))
}

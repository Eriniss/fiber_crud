package main

import (
	"fiber_curd/database"
	"fiber_curd/routes"
	"log"
	"os"

	"github.com/gofiber/fiber/v3"
	"github.com/joho/godotenv"
)

func main() {
	// .env 파일 내의 변수는 godotenv 패키지가 있어야만 사용 가능
	// godotenv는 main.go(root)만 선언하면 os.Getenv()메서드로 어디에서든 사용 가능
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	port := os.Getenv("API_PORT")

	// 초기 데이터베이스 생성 및 초기화
	database.InitDatabase()

	// Fiber 앱 생성
	app := fiber.New()

	// 라우트 설정
	routes.UserRoutes(app)
	routes.BlogRoutes(app)

	// 서버 실행
	log.Println("🚀 Server's hot in 3000 port!")
	log.Fatal(app.Listen(":" + port))
}

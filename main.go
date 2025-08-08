package main

import (
	"fiber_curd/database"
	"fiber_curd/routes"
	"log"
	"os"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
)

func main() {
	app := fiber.New()

	app.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept"},
		AllowCredentials: false,
	}))

	database.InitDatabase()
	routes.UserRoutes(app)
	routes.BlogRoutes(app)

	port := os.Getenv("API_PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("🚀 HTTPS server on :%s", port)
	log.Fatal(app.Listen(":"+port, fiber.ListenConfig{
		CertFile:    "./localhost.crt",
		CertKeyFile: "./localhost.key",
	}))
}

package routes

import (
	"fiber_curd/handlers"

	"github.com/gofiber/fiber/v3"
)

func UserRoutes(app *fiber.App) {
	api := app.Group("/auth")

	api.Get("/users", handlers.GetAllUsers)
	api.Get("/user/:id", handlers.GetUser)
	api.Post("/user", handlers.CreateUser)
	api.Put("/user/:id", handlers.UpdateUser)
	api.Delete("/user/:id", handlers.DeleteUser)
}

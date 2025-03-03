package routes

import (
	user "fiber_curd/handlers/user_handlers"

	"github.com/gofiber/fiber/v3"
)

func UserRoutes(app *fiber.App) {
	api := app.Group("/auth")

	api.Get("/users", user.GetAllUsers)
	api.Get("/user/:id", user.GetUser)
	api.Post("/user", user.CreateUser)
	api.Put("/user/:id", user.UpdateUser)
	api.Delete("/user/:id", user.DeleteUser)
}

package routes

import (
	blog "fiber_curd/handlers/blog_handlers"

	"github.com/gofiber/fiber/v3"
)

func BlogRoutes(app *fiber.App) {
	api := app.Group("/blog")

	api.Post("/post", blog.CreatePost)
	api.Get("/posts", blog.GetAllPosts)
	api.Get("/post/:id", blog.GetPost)
	api.Put("/post/:id", blog.UpdatePost)
	api.Delete("/post/:id", blog.DeletePost)
}

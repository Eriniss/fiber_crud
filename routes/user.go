package routes

import (
	user "fiber_crud/handlers/user_handlers"
	"fiber_crud/middleware"

	"github.com/gofiber/fiber/v3"
)

func UserRoutes(app *fiber.App) {
	api := app.Group("/auth")

	// 인증 불필요 (Public) 엔드포인트
	api.Post("/user", user.CreateUser)      // 회원가입
	api.Post("/sign_in", user.SignInUser)   // 로그인

	// 인증 필요 (Protected) 엔드포인트
	protected := api.Use(middleware.AuthRequired())
	protected.Get("/users", user.GetAllUsers)
	protected.Get("/user/:id", user.GetUser)
	protected.Put("/user/:id", user.UpdateUser)
	protected.Delete("/user/:id", user.DeleteUser)
}

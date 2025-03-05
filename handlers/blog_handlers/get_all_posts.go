package blog

import (
	"fiber_curd/database"
	"fiber_curd/models"

	"github.com/gofiber/fiber/v3"
)

// 모든 Posts 조회
func GetAllPosts(c fiber.Ctx) error {
	var blogs []models.Blog
	database.DB.Find(&blogs)

	return c.JSON(blogs)
}

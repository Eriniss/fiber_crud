package blog

import (
	"fiber_crud/database"
	"fiber_crud/models"

	"github.com/gofiber/fiber/v3"
)

// 단일 Posts 조회
func GetPost(c fiber.Ctx) error {
	id := c.Params("id")
	var post models.Blog

	if err := database.DB.First(&post, id).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Post not found"})
	}

	return c.JSON(post)
}

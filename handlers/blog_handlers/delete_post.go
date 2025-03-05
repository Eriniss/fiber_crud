package blog

import (
	"fiber_curd/database"
	"fiber_curd/models"

	"github.com/gofiber/fiber/v3"
)

// 단일 Posts 삭제
func DeletePost(c fiber.Ctx) error {
	id := c.Params("id")
	var post models.Blog

	if err := database.DB.First(&post, id).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Post not found"})
	}

	if err := database.DB.Delete(&post).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to delete post"})
	}

	return c.JSON(fiber.Map{"message": "Post deleted", "post": post})
}

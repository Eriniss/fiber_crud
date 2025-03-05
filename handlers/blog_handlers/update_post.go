package blog

import (
	"fiber_curd/database"
	"fiber_curd/models"

	"github.com/gofiber/fiber/v3"
)

// 단일 Posts 수정
func UpdatePost(c fiber.Ctx) error {
	id := c.Params("id")
	var post models.Blog

	if err := database.DB.First(&post, id).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Post not found"})
	}

	if err := c.Bind().Body(&post); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid input"})
	}

	database.DB.Model(&post).Updates(post)
	return c.JSON(fiber.Map{"message": "Post updated", "post": post})
}

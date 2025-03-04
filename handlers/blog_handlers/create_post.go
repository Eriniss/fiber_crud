package blog

import (
	"fiber_curd/database"
	"fiber_curd/models"

	"github.com/gofiber/fiber/v3"
)

// 단일 Post 생성
func CreatePost(c fiber.Ctx) error {
	// blogs 콜렉션 매핑
	post := new(models.Blog)

	// 입력값을 바인딩
	if err := c.Bind().Body(post); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid input"})
	}

	// 블로그 데이터베이스에 저장
	if err := database.DB.Create(&post).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to create post"})
	}

	return c.JSON(fiber.Map{"message": "Post created", "post": post})
}

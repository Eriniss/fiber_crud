package user

import (
	"fiber_curd/database"
	"fiber_curd/models"

	"github.com/gofiber/fiber/v3"
)

// 모든 Users 조회
func GetAllUsers(c fiber.Ctx) error {
	var users []models.User
	// password 필드 제외하고 전체 조회
	if err := database.DB.Omit("password").Find(&users).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to retrieve users"})
	}

	return c.JSON(users)
}

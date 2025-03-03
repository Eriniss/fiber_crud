package user

import (
	"fiber_curd/database"
	"fiber_curd/models"

	"github.com/gofiber/fiber/v3"
)

// 단일 User 삭제
func DeleteUser(c fiber.Ctx) error {
	id := c.Params("id")
	var user models.User

	// URL 파라미터에서 ID 추출 및 변환
	if err := database.DB.Find(&user, id).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "User not found"})
	}

	// User 삭제
	// 이미 삭제된 User의 삭제를 시도하거나, 삭제가 올바르게 작동하지 않은 경우 에러 반환
	if err := database.DB.Delete(&user).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to delete user"})
	}

	return c.JSON(fiber.Map{"message": "User deleted", "user": user})
}

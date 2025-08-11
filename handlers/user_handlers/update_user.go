package user

import (
	"fiber_crud/database"
	"fiber_crud/models"

	"github.com/gofiber/fiber/v3"
)

// 단일 User 수정
func UpdateUser(c fiber.Ctx) error {
	id := c.Params("id")
	var user models.User

	// DB에서 기존 사용자 조회
	if err := database.DB.First(&user, id).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "User not found"})
	}

	// 요청값을 받을 구조체
	var input models.User
	if err := c.Bind().Body(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid input"})
	}

	// Point 값이 요청에 포함된 경우 연산 처리
	if input.Point != 0 {
		user.Point += input.Point // 기존 값에서 더하거나 빼기
	}

	// 변경 사항 저장
	if err := database.DB.Save(&user).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to update user"})
	}

	return c.JSON(fiber.Map{
		"message": "User updated",
		"user":    user,
	})
}

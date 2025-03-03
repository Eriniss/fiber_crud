package user

import (
	"fiber_curd/database"
	"fiber_curd/models"

	"github.com/gofiber/fiber/v3"
)

// 단일 User 수정
func UpdateUser(c fiber.Ctx) error {
	id := c.Params("id")
	var user models.User

	// URL 파라미터에서 ID 추출 및 변환
	if err := database.DB.First(&user, id).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "User not found"})
	}

	// 입력된 값 검증
	// 입력값이 올바르지 않을 경우 400 에러 반환
	if err := c.Bind().Body(&user); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid input"})
	}

	// 최종적으로 DB에 id값과 일치하는 user를 업데이트
	// Save() 대신 Update()를 사용하여 해당 ID가 없는 경우 새로 생성되는 현상을 방지
	database.DB.Model(&user).Updates(user)
	return c.JSON(fiber.Map{"message": "User updated", "user": user})
}

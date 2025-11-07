package user

import (
	"fiber_crud/database"
	"fiber_crud/models"
	"fiber_crud/utils"

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

	// 이메일은 변경 불가 (보안상 중요 필드)
	if input.Email != "" && input.Email != user.Email {
		return c.Status(403).JSON(fiber.Map{"error": "Email cannot be changed"})
	}

	// 비밀번호 변경 시 검증 및 bcrypt 해싱 적용
	if input.Password != "" {
		// 비밀번호 강도 검증
		if err := utils.ValidatePassword(input.Password); err != nil {
			return c.Status(400).JSON(fiber.Map{"error": err.Error()})
		}

		hashedPassword, err := utils.HashPassword(input.Password)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "Failed to hash password"})
		}
		user.Password = hashedPassword
	}

	// 다른 필드 업데이트
	if input.Name != "" {
		user.Name = input.Name
	}
	if input.Group != "" {
		user.Group = input.Group
	}
	if input.Gender != "" {
		user.Gender = input.Gender
	}
	if input.Point != 0 {
		user.Point += input.Point // 기존 값에서 더하거나 빼기
	}

	// 변경 사항 저장
	if err := database.DB.Save(&user).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to update user"})
	}

	// 응답에서 비밀번호 제거
	user.Password = ""
	return c.JSON(fiber.Map{
		"message": "User updated",
		"user":    user,
	})
}

package user

import (
	"fiber_curd/database"
	"fiber_curd/models"
	"fiber_curd/utils"

	"github.com/gofiber/fiber/v3"
)

// 단일 User 생성
func CreateUser(c fiber.Ctx) error {
	// users 콜렉션 매핑
	user := new(models.User)

	// 입력값을 바인딩
	if err := c.Bind().Body(user); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid input"})
	}

	// 이메일 중복 여부 확인
	var existingUser models.User
	if err := database.DB.Where("email = ?", user.Email).First(&existingUser).Error; err == nil {
		return c.Status(409).JSON(fiber.Map{"error": "Email already exists"})
	}

	// 비밀번호 SHA-256 해시 적용
	user.Password = utils.HashPassword(user.Password)

	// 사용자 데이터베이스에 저장
	if err := database.DB.Create(&user).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to create user"})
	}

	return c.JSON(fiber.Map{"message": "User created", "user": user})
}

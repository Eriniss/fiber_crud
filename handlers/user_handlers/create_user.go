package user

import (
	"fiber_crud/database"
	"fiber_crud/models"
	"fiber_crud/utils"

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

	// 이메일 형식 검증
	if err := utils.ValidateEmail(user.Email); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	// 비밀번호 강도 검증
	if err := utils.ValidatePassword(user.Password); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	// 이메일 중복 여부 확인
	var existingUser models.User
	if err := database.DB.Where("email = ?", user.Email).First(&existingUser).Error; err == nil {
		return c.Status(409).JSON(fiber.Map{"error": "Email already exists"})
	}

	// 비밀번호 bcrypt 해시 적용
	hashedPassword, err := utils.HashPassword(user.Password)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to hash password"})
	}
	user.Password = hashedPassword

	// 사용자 데이터베이스에 저장
	if err := database.DB.Create(&user).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to create user"})
	}

	// 응답에서 비밀번호 제거
	user.Password = ""
	return c.JSON(fiber.Map{"message": "User created", "user": user})
}

package user

import (
	"crypto/sha256"
	"encoding/hex"
	"fiber_curd/database"
	"fiber_curd/models"

	"github.com/gofiber/fiber/v3"
)

// 단일 User 생성
func CreateUser(c fiber.Ctx) error {
	// database.db 내 User 콜렉션 새롭게 생성
	// 이미 있는 경우 User 콜렉션 매핑
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
	hashedPassword := sha256.Sum256([]byte(user.Password)) // SHA-256 해시 생성
	user.Password = hex.EncodeToString(hashedPassword[:])  // 해시값을 hex 문자열로 변환

	// 사용자 데이터베이스에 저장
	if err := database.DB.Create(&user).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to create user"})
	}

	return c.JSON(fiber.Map{"message": "User created", "user": user})
}

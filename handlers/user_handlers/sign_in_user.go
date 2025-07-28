package user

import (
	"fiber_curd/database"
	"fiber_curd/models"
	"fiber_curd/utils"

	"github.com/gofiber/fiber/v3"
)

// 사용자 로그인
func SignInUser(c fiber.Ctx) error {
	// 입력값 바인딩
	loginData := new(models.User)
	if err := c.Bind().Body(loginData); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid input"})
	}

	var user models.User

	// 이메일로 사용자 검색
	if err := database.DB.Where("email = ?", loginData.Email).First(&user).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "User not found"})
	}

	// 입력된 비밀번호를 해시 후 비교
	hashedPassword := utils.HashPassword(loginData.Password)
	if hashedPassword != user.Password {
		return c.Status(401).JSON(fiber.Map{"error": "Invalid password"})
	}

	// 비밀번호를 제외하고 응답 데이터 준비
	user.Password = "" // 비밀번호 필드 삭제 후 반환

	// JWT 토큰 생성
	token, err := utils.GenerateJWT(user.ID, user.Email, user.Role, user.Point)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to generate token"})
	}

	return c.JSON(fiber.Map{
		"message": "Login successful",
		"user": fiber.Map{
			"id":    user.ID,
			"email": user.Email,
			"role":  user.Role,
		},
		"token": token, // JWT 토큰 반환
	})
}

package handlers

import (
	"fiber_curd/database"
	"fiber_curd/models"

	"github.com/gofiber/fiber/v3"
)

// 전체 User 조회
func GetAllUsers(c fiber.Ctx) error {
	// 단일 객체가 아닐때는 []를 앞에 표시하여 슬라이스 형태로 반환
	var users []models.User
	database.DB.Find(&users)

	return c.JSON(users)
}

// id값을 바탕으로 특정 User 조회
func GetUser(c fiber.Ctx) error {
	// 라우트에 설정된 /:id 값을 가져와 변수에 할당
	id := c.Params("id")
	var user models.User

	result := database.DB.First(&user, id)
	if result.Error != nil {
		return c.Status(404).JSON(fiber.Map{"error": "User not found"})
	}

	return c.JSON(user)
}

// 단일 User 생성
func CreateUser(c fiber.Ctx) error {
	// database.db 내 User 콜렉션 새롭게 생성
	// 이미 있는 경우 User 콜렉션 매핑
	user := new(models.User)

	if err := c.Bind().Body(user); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid input"})
	}

	// 이메일 중복 여부 확인
	var existingUser models.User
	if err := database.DB.Where("email = ?", user.Email).First(&existingUser).Error; err == nil {
		return c.Status(409).JSON(fiber.Map{"error": "Email already exists"})
	}

	database.DB.Create(&user)
	return c.JSON(fiber.Map{"message": "User created", "user": user})
}

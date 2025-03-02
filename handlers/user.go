package handlers

import (
	"fiber_curd/database"
	"fiber_curd/models"

	"github.com/gofiber/fiber/v3"
)

// 전체 User 조회
func GetAllUsers(c fiber.Ctx) error {
	var users []models.User
	// users의 포인터를 설정하여 해당 변수의 메모리를 직접 참조.
	database.DB.Find(&users)

	return c.JSON(users)
}

// id값을 바탕으로 특정 User 조회
func GetUser(c fiber.Ctx) error {
	id := c.Params("id")
	var user []models.User

	// 마찬가지로, user의 포인터를 설정하여 해당 변수의 메모리를 직접 참조.
	result := database.DB.First(&user, id)
	if result.Error != nil {
		return c.Status(404).JSON(fiber.Map{"error": "User not found"})
	}

	return c.JSON(user)
}

// 단일 User 생성
func CreateUser(c fiber.Ctx) error {
	user := new(models.User)

	if err := c.Bind().Body(user); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid input"})
	}

	database.DB.Create(&user)
	return c.Status(201).JSON(user, "User created")
}

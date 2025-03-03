package user

import (
	"fiber_curd/database"
	"fiber_curd/models"

	"github.com/gofiber/fiber/v3"
)

func GetAllUsers(c fiber.Ctx) error {
	// 단일 객체가 아닐때는 []를 앞에 표시하여 슬라이스 형태로 반환
	var users []models.User
	database.DB.Find(&users)

	return c.JSON(users)
}

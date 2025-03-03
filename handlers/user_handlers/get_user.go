package user

import (
	"fiber_curd/database"
	"fiber_curd/models"

	"github.com/gofiber/fiber/v3"
)

func GetUser(c fiber.Ctx) error {
	// 라우트에 설정된 /:id 값을 가져와 변수에 할당(파라미터에서 값 추출)
	id := c.Params("id")
	var user models.User

	// 추출된 id값을 바탕으로 DB 검색
	// 최초 검색된 데이터(First 메서드) result에 저장
	if err := database.DB.First(&user, id).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "User not found"})
	}

	return c.JSON(user)
}

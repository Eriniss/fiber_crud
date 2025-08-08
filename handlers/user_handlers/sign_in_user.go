package user

import (
	"fiber_curd/database"
	"fiber_curd/models"
	"fiber_curd/utils"
	"log"

	"github.com/gofiber/fiber/v3"
)

// 사용자 로그인
func SignInUser(c fiber.Ctx) error {
	log.Println("[SignInUser] 로그인 요청 수신")

	// 입력값 바인딩
	loginData := new(models.User)
	if err := c.Bind().Body(loginData); err != nil {
		log.Printf("[SignInUser] Body 바인딩 실패: %v\n", err)
		return c.Status(400).JSON(fiber.Map{"error": "Invalid input"})
	}
	log.Printf("[SignInUser] 파싱된 입력값: %+v\n", loginData)

	var user models.User

	// 이메일로 사용자 검색
	if err := database.DB.Where("email = ?", loginData.Email).First(&user).Error; err != nil {
		log.Printf("[SignInUser] 사용자 조회 실패 (email=%s): %v\n", loginData.Email, err)
		return c.Status(404).JSON(fiber.Map{"error": "User not found"})
	}
	log.Printf("[SignInUser] DB 조회 성공: %+v\n", user)

	// 입력된 비밀번호를 해시 후 비교
	hashedPassword := utils.HashPassword(loginData.Password)
	log.Printf("[SignInUser] 입력 비밀번호 해시값: %s\n", hashedPassword)
	log.Printf("[SignInUser] DB 저장된 비밀번호 해시값: %s\n", user.Password)

	if hashedPassword != user.Password {
		log.Println("[SignInUser] 비밀번호 불일치")
		return c.Status(401).JSON(fiber.Map{"error": "Invalid password"})
	}

	// 비밀번호를 제외하고 응답 데이터 준비
	user.Password = "" // 비밀번호 필드 삭제 후 반환

	// JWT 토큰 생성
	token, err := utils.GenerateJWT(user.ID, user.Email, user.Name, uint(user.Point))
	if err != nil {
		log.Printf("[SignInUser] JWT 생성 실패: %v\n", err)
		return c.Status(500).JSON(fiber.Map{"error": "Failed to generate token"})
	}

	log.Println("[SignInUser] 로그인 성공, 토큰 발급 완료")
	return c.JSON(fiber.Map{
		"message": "Login successful",
		"user": fiber.Map{
			"id":     user.ID,
			"email":  user.Email,
			"name":   user.Name,
			"group":  user.Group,
			"gender": user.Gender,
			"point":  user.Point,
		},
		"token": token, // JWT 토큰 반환
	})
}

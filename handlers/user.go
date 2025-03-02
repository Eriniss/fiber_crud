package handlers

import (
	"crypto/sha256"
	"encoding/hex"
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

// 단일 User 수정
func UpdateUser(c fiber.Ctx) error {
	id := c.Params("id")
	var user models.User

	// URL 파라미터에서 ID 추출 및 변환
	if err := database.DB.First(&user, id).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "User not found"})
	}

	// 입력된 값 검증
	// 입력값이 올바르지 않을 경우 400 에러 반환
	if err := c.Bind().Body(&user); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid input"})
	}

	// 최종적으로 DB에 id값과 일치하는 user를 업데이트
	// Save() 대신 Update()를 사용하여 해당 ID가 없는 경우 새로 생성되는 현상을 방지
	database.DB.Model(&user).Updates(user)
	return c.JSON(fiber.Map{"message": "User updated", "user": user})
}

// 단일 User 삭제
func DeleteUser(c fiber.Ctx) error {
	id := c.Params("id")
	var user models.User

	// URL 파라미터에서 ID 추출 및 변환
	if err := database.DB.Find(&user, id).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "User not found"})
	}

	// User 삭제
	// 이미 삭제된 User의 삭제를 시도하거나, 삭제가 올바르게 작동하지 않은 경우 에러 반환
	if err := database.DB.Delete(&user).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to delete user"})
	}

	return c.JSON(fiber.Map{"message": "User deleted", "user": user})
}

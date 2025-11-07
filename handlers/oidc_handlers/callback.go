package oidc

import (
	"fiber_crud/database"
	"fiber_crud/models"
	"fiber_crud/utils"

	"github.com/gofiber/fiber/v3"
)

// OIDC 콜백 핸들러
// Logto에서 인증 후 리다이렉트되는 엔드포인트
func Callback(c fiber.Ctx) error {
	// 인증 코드 가져오기
	code := c.Query("code")
	if code == "" {
		return c.Status(400).JSON(fiber.Map{
			"error": "Authorization code is missing",
		})
	}

	// TODO: 실제 Logto SDK 사용 시 구현
	// 1. 코드를 사용하여 토큰 교환
	// 2. 토큰을 사용하여 사용자 정보 가져오기
	// 3. 사용자 정보로 DB 조회 또는 생성

	// 예시: 사용자 정보 (실제로는 Logto에서 가져와야 함)
	// userInfo := map[string]interface{}{
	// 	"email": "user@example.com",
	// 	"name":  "John Doe",
	// 	"sub":   "logto_user_id_123",
	// }

	// 임시 응답 (실제 구현 시 제거)
	return c.JSON(fiber.Map{
		"message": "OIDC callback received",
		"code":    code,
		"note":    "This endpoint needs to be implemented with Logto SDK",
	})
}

// OIDC 사용자 정보로 로컬 사용자 생성 또는 조회
func getOrCreateOIDCUser(email, name, oidcSub string) (*models.User, error) {
	var user models.User

	// 이메일로 기존 사용자 검색
	if err := database.DB.Where("email = ?", email).First(&user).Error; err != nil {
		// 사용자가 없으면 새로 생성
		user = models.User{
			Email:  email,
			Name:   name,
			Group:  "user", // 기본 권한
			Point:  0,
		}

		if err := database.DB.Create(&user).Error; err != nil {
			return nil, err
		}
	}

	return &user, nil
}

// OIDC 로그아웃
func Logout(c fiber.Ctx) error {
	// TODO: Logto 로그아웃 URL로 리다이렉트
	// 클라이언트에서 토큰 제거 필요

	return c.JSON(fiber.Map{
		"message": "Logout successful",
		"note":    "Client should remove the JWT token",
	})
}

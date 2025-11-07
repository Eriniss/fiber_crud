package oidc

import (
	"context"
	"fiber_crud/database"
	"fiber_crud/models"
	"fiber_crud/utils"

	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/gofiber/fiber/v3"
)

// ID Token Claims 구조체
type LogtoClaims struct {
	Sub           string   `json:"sub"`
	Email         string   `json:"email"`
	EmailVerified bool     `json:"email_verified"`
	Name          string   `json:"name"`
	Picture       string   `json:"picture"`
	Roles         []string `json:"roles"`
}

// OIDC 콜백 핸들러
// Logto에서 인증 후 리다이렉트되는 엔드포인트
func Callback(c fiber.Ctx) error {
	ctx := context.Background()

	// OIDC가 활성화되어 있는지 확인
	if !IsOIDCEnabled() {
		return c.Status(503).JSON(fiber.Map{
			"error": "OIDC is not configured",
		})
	}

	// 1. State 검증 (CSRF 방어)
	state := c.Query("state")
	cookieState := c.Cookies("oauth_state")

	if state == "" || state != cookieState {
		return c.Status(400).JSON(fiber.Map{
			"error": "Invalid state parameter",
		})
	}

	// State 쿠키 삭제
	c.ClearCookie("oauth_state")

	// 2. Authorization Code 가져오기
	code := c.Query("code")
	if code == "" {
		return c.Status(400).JSON(fiber.Map{
			"error": "Authorization code is missing",
		})
	}

	// 3. 인증 코드를 토큰으로 교환
	oauth2Token, err := GetOAuth2Config().Exchange(ctx, code)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Failed to exchange token",
			"details": err.Error(),
		})
	}

	// 4. ID Token 추출 및 검증
	rawIDToken, ok := oauth2Token.Extra("id_token").(string)
	if !ok {
		return c.Status(500).JSON(fiber.Map{
			"error": "No id_token in response",
		})
	}

	verifier := GetProvider().Verifier(&oidc.Config{
		ClientID: GetOAuth2Config().ClientID,
	})

	idToken, err := verifier.Verify(ctx, rawIDToken)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Failed to verify ID token",
			"details": err.Error(),
		})
	}

	// 5. Claims 파싱
	var claims LogtoClaims
	if err := idToken.Claims(&claims); err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Failed to parse claims",
			"details": err.Error(),
		})
	}

	// 6. 사용자 조회 또는 생성
	user, err := getOrCreateOIDCUser(claims)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Failed to create user",
			"details": err.Error(),
		})
	}

	// 7. JWT 토큰 생성 (우리 시스템용)
	token, err := utils.GenerateJWT(user.ID, user.Email, user.Name, uint(user.Point))
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Failed to generate token",
		})
	}

	// 8. 사용자 정보 및 토큰 반환
	return c.JSON(fiber.Map{
		"message": "Login successful",
		"user": fiber.Map{
			"id":            user.ID,
			"email":         user.Email,
			"name":          user.Name,
			"group":         user.Group,
			"is_oidc_user":  user.IsOIDCUser,
		},
		"token": token,
	})
}

// OIDC 사용자 조회 또는 생성
func getOrCreateOIDCUser(claims LogtoClaims) (*models.User, error) {
	var user models.User

	// OIDC Subject로 기존 사용자 검색
	err := database.DB.Where("oidc_subject = ? AND oidc_provider = ?",
		claims.Sub, "logto").First(&user).Error

	if err == nil {
		// 기존 사용자: 정보 업데이트
		user.Name = claims.Name
		user.Group = mapLogtoRolesToGroup(claims.Roles)
		database.DB.Save(&user)
		return &user, nil
	}

	// 이메일로 기존 사용자 검색 (연동 가능)
	err = database.DB.Where("email = ?", claims.Email).First(&user).Error
	if err == nil {
		// 기존 이메일 사용자를 OIDC로 연동
		user.OIDCSubject = claims.Sub
		user.OIDCProvider = "logto"
		user.IsOIDCUser = true
		user.Group = mapLogtoRolesToGroup(claims.Roles)
		database.DB.Save(&user)
		return &user, nil
	}

	// 새 사용자 생성
	user = models.User{
		Email:        claims.Email,
		Name:         claims.Name,
		Group:        mapLogtoRolesToGroup(claims.Roles),
		OIDCSubject:  claims.Sub,
		OIDCProvider: "logto",
		IsOIDCUser:   true,
		Point:        0,
	}

	if err := database.DB.Create(&user).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

// Logto Role을 로컬 Group으로 매핑
func mapLogtoRolesToGroup(roles []string) string {
	// Admin Role이 있으면 admin
	for _, role := range roles {
		if role == "admin" {
			return "admin"
		}
	}

	// 기본값은 user
	return "user"
}

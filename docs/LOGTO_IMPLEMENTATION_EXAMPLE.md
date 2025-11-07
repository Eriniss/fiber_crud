# Logto 구현 예시 코드

이 문서는 Logto OIDC 통합을 위한 실제 구현 코드 예시를 제공합니다.

## 1. User 모델 확장

```go
// models/user.go
package models

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model

	// 기존 필드
	Email    string `gorm:"unique" json:"email"`
	Password string `json:"-"`                    // JSON 응답에서 제외
	Name     string `json:"name"`
	Group    string `json:"group"`                // admin, user
	Gender   string `json:"gender"`
	Point    int    `json:"point"`

	// OIDC 필드 추가
	OIDCSubject  string `gorm:"index" json:"-"`   // Logto user ID (sub claim)
	OIDCProvider string `json:"-"`                 // "logto"
	IsOIDCUser   bool   `json:"is_oidc_user"`     // OIDC 로그인 여부
}

// OIDC 사용자 여부 확인
func (u *User) IsFromOIDC() bool {
	return u.IsOIDCUser && u.OIDCSubject != ""
}

// 비밀번호 인증이 필요한지 확인
func (u *User) RequiresPassword() bool {
	return !u.IsOIDCUser
}
```

## 2. OIDC Config 설정

```go
// handlers/oidc_handlers/config.go
package oidc

import (
	"context"
	"os"

	"github.com/coreos/go-oidc/v3/oidc"
	"golang.org/x/oauth2"
)

var (
	provider     *oidc.Provider
	oauth2Config *oauth2.Config
)

// OIDC Provider 초기화
func InitOIDCProvider() error {
	ctx := context.Background()

	// Logto OIDC Provider 설정
	logtoEndpoint := os.Getenv("LOGTO_ENDPOINT")
	var err error
	provider, err = oidc.NewProvider(ctx, logtoEndpoint)
	if err != nil {
		return err
	}

	// OAuth2 Config 설정
	oauth2Config = &oauth2.Config{
		ClientID:     os.Getenv("LOGTO_APP_ID"),
		ClientSecret: os.Getenv("LOGTO_APP_SECRET"),
		RedirectURL:  os.Getenv("LOGTO_REDIRECT_URI"),
		Endpoint:     provider.Endpoint(),
		Scopes:       []string{oidc.ScopeOpenID, "profile", "email", "roles"},
	}

	return nil
}

// OAuth2 Config 반환
func GetOAuth2Config() *oauth2.Config {
	return oauth2Config
}

// OIDC Provider 반환
func GetProvider() *oidc.Provider {
	return provider
}
```

## 3. 로그인 핸들러 구현

```go
// handlers/oidc_handlers/login.go
package oidc

import (
	"crypto/rand"
	"encoding/base64"
	"os"

	"github.com/gofiber/fiber/v3"
)

// State 토큰 생성 (CSRF 방어)
func generateStateToken() (string, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

// OIDC 로그인 리다이렉트
func LoginRedirect(c fiber.Ctx) error {
	// State 토큰 생성
	state, err := generateStateToken()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Failed to generate state token",
		})
	}

	// State를 세션 또는 쿠키에 저장 (CSRF 방어)
	c.Cookie(&fiber.Cookie{
		Name:     "oauth_state",
		Value:    state,
		HTTPOnly: true,
		Secure:   os.Getenv("ENVIRONMENT") == "production",
		SameSite: "Lax",
		MaxAge:   600, // 10분
	})

	// OAuth2 Authorization URL 생성
	authURL := GetOAuth2Config().AuthCodeURL(state)

	// 프론트엔드가 있는 경우 URL 반환
	return c.JSON(fiber.Map{
		"auth_url": authURL,
		"message":  "Redirect user to this URL",
	})

	// 또는 직접 리다이렉트
	// return c.Redirect(authURL)
}
```

## 4. 콜백 핸들러 구현

```go
// handlers/oidc_handlers/callback.go
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
func Callback(c fiber.Ctx) error {
	ctx := context.Background()

	// 1. State 검증 (CSRF 방어)
	state := c.Query("code")
	cookieState := c.Cookies("oauth_state")

	if state != cookieState {
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
		})
	}

	// 5. Claims 파싱
	var claims LogtoClaims
	if err := idToken.Claims(&claims); err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Failed to parse claims",
		})
	}

	// 6. 사용자 조회 또는 생성
	user, err := getOrCreateOIDCUser(claims)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Failed to create user",
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
```

## 5. 로그아웃 핸들러

```go
// handlers/oidc_handlers/logout.go
package oidc

import (
	"fmt"
	"net/url"
	"os"

	"github.com/gofiber/fiber/v3"
)

// OIDC 로그아웃
func Logout(c fiber.Ctx) error {
	// Logto 로그아웃 URL 생성
	logtoEndpoint := os.Getenv("LOGTO_ENDPOINT")
	postLogoutRedirect := os.Getenv("LOGTO_POST_LOGOUT_REDIRECT_URI")

	logoutURL := fmt.Sprintf("%s/oidc/session/end?post_logout_redirect_uri=%s",
		logtoEndpoint,
		url.QueryEscape(postLogoutRedirect),
	)

	return c.JSON(fiber.Map{
		"logout_url": logoutURL,
		"message":    "Redirect user to this URL to logout from Logto",
	})
}
```

## 6. main.go 초기화

```go
// main.go
package main

import (
	"fiber_crud/database"
	oidc "fiber_crud/handlers/oidc_handlers"
	"fiber_crud/routes"
	"log"
	"os"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
	"github.com/joho/godotenv"
)

func main() {
	// .env 로드
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	// Fiber 앱 생성
	app := fiber.New()

	// CORS 설정
	allowedOrigins := os.Getenv("ALLOWED_ORIGINS")
	if allowedOrigins == "" {
		allowedOrigins = "http://localhost:3000"
	}

	app.Use(cors.New(cors.Config{
		AllowOrigins:     []string{allowedOrigins},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	port := os.Getenv("API_PORT")

	// DB 초기화
	database.InitDatabase()

	// OIDC Provider 초기화
	if err := oidc.InitOIDCProvider(); err != nil {
		log.Printf("⚠️  Warning: OIDC initialization failed: %v", err)
		log.Println("OIDC features will be disabled")
	} else {
		log.Println("✅ OIDC Provider initialized")
	}

	// 라우트 설정
	routes.UserRoutes(app)
	routes.OIDCRoutes(app)

	// 서버 실행
	log.Printf("🚀 Server's hot in %s port!\n", port)
	log.Fatal(app.Listen(":" + port))
}
```

## 7. 데이터베이스 마이그레이션

기존 데이터베이스에 OIDC 필드를 추가하려면:

```bash
# 서버를 한 번 실행하면 GORM이 자동으로 새 필드를 추가합니다
go run main.go

# 또는 마이그레이션 스크립트 작성
```

## 8. 테스트

### REST Client 테스트 파일 업데이트

```http
# http/oidc.http

### 1. OIDC 로그인 URL 얻기
GET http://localhost:8080/oidc/login

### 2. 브라우저에서 로그인 후 콜백 테스트
# (실제로는 브라우저에서 진행)
GET http://localhost:8080/oidc/callback?code=<auth_code>&state=<state>

### 3. 로그아웃 URL 얻기
GET http://localhost:8080/oidc/logout
```

## 9. 환경 변수 예시

```env
# .env
API_PORT=8080
DATABASE_PATH=./database.db
JWT_SECRET_KEY=your-very-secure-secret-key-here
ALLOWED_ORIGINS=http://localhost:3000

# Logto OIDC
LOGTO_ENDPOINT=https://beh25r.logto.app
LOGTO_APP_ID=abc123xyz
LOGTO_APP_SECRET=secret_abc123xyz
LOGTO_REDIRECT_URI=http://localhost:8080/oidc/callback
LOGTO_POST_LOGOUT_REDIRECT_URI=http://localhost:3000

ENVIRONMENT=development
```

## 10. 프론트엔드 통합 예시

```javascript
// React/Vue/Next.js 예시
async function loginWithLogto() {
  // 1. 백엔드에서 로그인 URL 받기
  const response = await fetch('http://localhost:8080/oidc/login');
  const { auth_url } = await response.json();

  // 2. Logto 로그인 페이지로 리다이렉트
  window.location.href = auth_url;

  // 3. 콜백 페이지에서 토큰 받기 (자동 처리됨)
}

// 콜백 페이지 처리
// /callback 라우트에서
const urlParams = new URLSearchParams(window.location.search);
const token = urlParams.get('token'); // 백엔드가 쿼리로 전달

// 토큰 저장
localStorage.setItem('jwt_token', token);

// 홈으로 리다이렉트
window.location.href = '/';
```

## 다음 단계

1. 위 코드를 프로젝트에 통합
2. 환경 변수 설정
3. 서버 실행 및 테스트
4. 에러 처리 개선
5. 프론트엔드 통합

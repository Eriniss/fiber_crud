package models

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model

	Email    string `gorm:"unique" json:"email"`  // 이메일
	Password string `json:"-"`                    // Password는 json 응답에서 제외
	Name     string `json:"name"`                 // 이름
	Group    string `json:"group"`                // 역할. admin, user
	Gender   string `json:"gender"`               // 성별. male, female
	Point    int    `json:"point"`                // 점수

	// OIDC 필드
	OIDCSubject  string `gorm:"index" json:"-"`   // Logto user ID (sub claim)
	OIDCProvider string `json:"-"`                // OIDC provider 이름 (예: "logto")
	IsOIDCUser   bool   `json:"is_oidc_user"`     // OIDC로 로그인한 사용자인지 여부
}

// OIDC 사용자인지 확인
func (u *User) IsFromOIDC() bool {
	return u.IsOIDCUser && u.OIDCSubject != ""
}

// 비밀번호 인증이 필요한지 확인
func (u *User) RequiresPassword() bool {
	return !u.IsOIDCUser
}

package models

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model

	Email    string `gorm:"unique" json:"email"` // 이메일
	Password string `json:"password"`            // Password는 json 객체 반환에서
	Role     string `json:"role"`                // 역할. admin, user
	Gender   string `json:"gender"`              // 성별. male, female
	Point    uint   `json:"point"`               // 점수.
}

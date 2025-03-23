package models

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model

	Email    string `gorm:"unique" json:"email"`
	Password string `json:"password"` // Password는 json 객체 반환에서
	Role     uint   `json:"role"`
}

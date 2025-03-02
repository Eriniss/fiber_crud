package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Email    string `gorm:"unique" json:"email"` // 정책상 email 중복 불허
	Password string `json:"password"`
	Role     int    `json:"role"`
}

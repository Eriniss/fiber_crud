package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model

	Email     string     `gorm:"unique" json:"email"`
	Password  string     `json:"-"` // Password는 json 객체 반환에서 제외
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
}

package models

import (
	"gorm.io/gorm"
)

// Reply (댓글) 모델
type Reply struct {
	gorm.Model

	PostID   uint   `json:"post_id"` // Blog와 관계를 형성하는 외래 키
	Email    string `json:"email"`
	Contents string `json:"contents"`
}

package models

import (
	"gorm.io/gorm"
)

// Post (게시글) 모델
type Blog struct {
	gorm.Model

	Email    string   `json:"email"`
	Title    string   `json:"title"`
	Contents string   `json:"contents"`
	Tags     []string `gorm:"serializer:json" json:"tags"`
	Replies  []Reply  `gorm:"foreignKey:PostID" json:"replies"`
	Private  bool     `json:"private"`
}

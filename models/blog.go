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

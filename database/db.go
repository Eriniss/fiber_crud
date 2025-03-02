package database

import (
	"fiber_curd/models"
	"log"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDatabase() {
	var err error
	DB, err = gorm.Open(sqlite.Open("database.db"), &gorm.Config{})
	if err != nil {
		log.Fatal("📌 Failed to connect to the database:", err)
	}

	// 데이터베이스 마이그레이션
	DB.AutoMigrate(&models.User{})
	log.Println("📌 Database migration completed!")
}

package database

import (
	"fiber_curd/models"
	"log"
	"os"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDatabase() {
	// 데이터베이스 url 설정
	// 추후 MariaDB 또는 ProstgreSQL 마이그레이션을 고려하여 변수명 변경
	dbPath := os.Getenv("DATABASE_PATH")

	var err error
	DB, err = gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		log.Fatal("📌 Failed to connect to the database:", err)
	}

	// 데이터베이스 마이그레이션
	// 데이터베이스가 없을 경우 자동으로 추가
	DB.AutoMigrate(&models.User{})
	DB.AutoMigrate(&models.Blog{})

	log.Println("📌 Database migration completed!")
}

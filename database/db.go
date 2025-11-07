package database

import (
	"fiber_crud/models"
	"log"
	"os"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDatabase() {
	// 환경변수에서 데이터베이스 경로 읽기 (없으면 기본값 사용)
	dbPath := os.Getenv("DATABASE_PATH")
	log.Printf("[DB] Connecting to: %s", dbPath)

	var err error
	DB, err = gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		log.Fatal("📌 Failed to connect to the database:", err)
	}

	// 마이그레이션
	if err := DB.AutoMigrate(&models.User{}); err != nil {
		log.Fatal("📌 Database migration failed:", err)
	}

	log.Println("📌 Database migration completed!")
}

package models

import (
	"log"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Todo struct {
	ID          uint `gorm:"primaryKey"`
	Title       string
	Description string
	Completed   bool
}

func InitDB() *gorm.DB {

	db, err := gorm.Open(sqlite.Open("todo.db"), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}
	if err := db.AutoMigrate(&Todo{}); err != nil {
		log.Fatalf("failed to migrate database: %v", err)
	}

	return db
}

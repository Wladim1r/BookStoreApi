package database

import (
	"bookstore-api/internal/models"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db *gorm.DB

func InitDB() *gorm.DB {
	dsn := "host=localhost user=postgres password=password dbname=postgres port=5432 sslmode=disable"
	var err error

	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Could not connect to DB")
	}

	if err := db.AutoMigrate(&models.User{}, &models.Book{}); err != nil {
		log.Fatal("Error when creation DB")
	}

	return db
}

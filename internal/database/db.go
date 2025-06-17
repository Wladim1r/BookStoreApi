package database

import (
	"bookstore-api/internal/models"
	"fmt"
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var db *gorm.DB

func InitDB() *gorm.DB {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"),
	)

	cfg := &gorm.Config{}

	if os.Getenv("DEBUG_MODE") == "true" {
		cfg.Logger = logger.Default.LogMode(logger.Info)
	} else {
		cfg.Logger = logger.Default.LogMode(logger.Silent)
	}

	var err error

	db, err = gorm.Open(postgres.Open(dsn), cfg)
	if err != nil {
		log.Fatal("Could not connect to DB")
	}

	if err := db.AutoMigrate(&models.User{}, &models.Book{}); err != nil {
		log.Fatal("Error when creation DB")
	}

	return db
}

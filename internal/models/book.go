package models

import "gorm.io/gorm"

type Book struct {
	gorm.Model
	Title  string `json:"title" binding:"required"`
	Author string `json:"author" binding:"required"`
	Price  uint   `json:"price"`
	UserID uint   `json:"user_id"`
	User   User   `gorm:"foreignKey:UserID" json:"-"`
}

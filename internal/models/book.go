package models

type Book struct {
	ID     uint    `gorm:"primaryKey" json:"id" binding:"required"`
	UserID uint    `json:"user_id"`
	Title  string  `json:"title" binding:"required"`
	Author string  `json:"author" binding:"required"`
	Email  string  `json:"email" binding:"required,email"`
	Price  float32 `json:"price"`
	//User   User    `gorm:"foreignKey:UserID" json:"-"`
}

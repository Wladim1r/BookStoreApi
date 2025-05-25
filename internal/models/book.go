package models

// @Description Book model to show it contains
// @Example {"id":1,"title":"Война и мир","author":"Л. Н. Толстой","price":1300,"user_id":1}
type Book struct {
	ID     uint   `json:"id"      gorm:"primarykey"        example:"1"`
	Title  string `json:"title"                            example:"Война и мир"   binding:"required"`
	Author string `json:"author"                           example:"Л. Н. Толстой" binding:"required"`
	Price  uint   `json:"price"                            example:"1300"          binding:"required"`
	UserID uint   `json:"user_id"                          example:"1"`
	User   User   `json:"-"       gorm:"foreignKey:UserID"`
}

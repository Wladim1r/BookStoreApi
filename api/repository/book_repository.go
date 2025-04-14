package repository

import (
	"bookstore-api/internal/models"

	"gorm.io/gorm"
)

type BookRepository interface {
	GetAllBooks() ([]models.Book, error)
	GetUserBook(uint, uint) (models.Book, error)
	GetUserBooks(uint) ([]models.Book, error)
	PostBook(models.Book) (*models.Book, error)
	UpdateBook(uint, uint, models.Book) error
	DeleteBook(uint, uint) error
}

type bookRepository struct {
	db *gorm.DB
}

func NewBookRepository(db *gorm.DB) BookRepository {
	return &bookRepository{db: db}
}

func (r *bookRepository) GetAllBooks() ([]models.Book, error) {
	var books []models.Book
	err := r.db.Preload("User").Find(&books).Error

	return books, err
}

func (r *bookRepository) GetUserBook(userID, bookID uint) (models.Book, error) {
	var book models.Book
	err := r.db.Where("id = ? AND user_id = ?", bookID, userID).First(&book).Error

	return book, err
}

func (r *bookRepository) GetUserBooks(userID uint) ([]models.Book, error) {
	var book []models.Book
	err := r.db.Where("user_id = ?", userID).Find(&book).Error

	return book, err
}

func (r *bookRepository) PostBook(book models.Book) (*models.Book, error) {
	err := r.db.Create(&book).Error
	return &book, err
}

func (r *bookRepository) UpdateBook(userID, bookID uint, book models.Book) error {
	return r.db.Model(&models.Book{}).Where("id = ? AND user_id = ?", bookID, userID).
		Select("Title", "Author", "Price").
		Updates(book).Error
}

func (r *bookRepository) DeleteBook(userID, bookID uint) error {
	return r.db.Unscoped().Where("id = ? AND user_id = ?", bookID, userID).Delete(&models.Book{}).Error
}

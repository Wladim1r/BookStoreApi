package service

import (
	"bookstore-api/internal/models"

	"gorm.io/gorm"
)

type BookRepository interface {
	GetBook(uint) (models.Book, error)
	GetBooks() ([]models.Book, error)
	PostBook(models.Book) error
	UpdateBook(string, models.Book) error
	DeleteBook(uint) error
}

type bookRepository struct {
	db *gorm.DB
}

func NewBookRepository(db *gorm.DB) BookRepository {
	return &bookRepository{db: db}
}

func (r *bookRepository) GetBook(id uint) (models.Book, error) {
	var book models.Book
	err := r.db.First(&book, id).Error

	return book, err
}

func (r *bookRepository) GetBooks() ([]models.Book, error) {
	var books []models.Book
	err := r.db.Find(&books).Error

	return books, err
}

func (r *bookRepository) PostBook(book models.Book) error {
	return r.db.Create(&book).Error
}

func (r *bookRepository) UpdateBook(id string, book models.Book) error {
	return r.db.Model(&models.Book{}).Where("id = ?", id).Updates(book).Error
}

func (r *bookRepository) DeleteBook(id uint) error {
	return r.db.Delete(&models.Book{}, id).Error
}

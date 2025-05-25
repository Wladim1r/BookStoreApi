package repository

import (
	"bookstore-api/internal/lib/errs"
	"bookstore-api/internal/models"
	"fmt"

	"gorm.io/gorm"
)

type BookRepository interface {
	GetAllBooks() ([]models.Book, error)
	GetUserBooks(uint, string, string, int) ([]models.Book, error)
	PostBook(models.Book) error
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
	result := r.db.Preload("User").Find(&books)

	if result.Error != nil {
		return nil, fmt.Errorf("%w: %v", errs.ErrDBOperation, result.Error)
	}
	if result.RowsAffected == 0 {
		return nil, errs.ErrNotFound
	}

	return books, nil
}

func (r *bookRepository) GetUserBooks(
	userID uint,
	author, title string,
	limit int,
) ([]models.Book, error) {
	var books []models.Book

	query := r.db.Model(&models.Book{}).
		Select("id, title, author, price").
		Where("user_id = ?", userID)

	if author != "" {
		query = query.Where("author ILIKE ?", "%"+author+"%")
	}

	if title != "" {
		query = query.Where("title ILIKE ?", "%"+title+"%")
	}

	if limit > 0 {
		query = query.Limit(limit)
	}

	result := query.Find(&books)
	if result.Error != nil {
		return nil, fmt.Errorf("%w: %v", errs.ErrDBOperation, result.Error)
	}

	if result.RowsAffected == 0 {
		return nil, errs.ErrNotFound
	}

	return books, nil
}

func (r *bookRepository) PostBook(book models.Book) error {
	result := r.db.Create(&book)

	if result.Error != nil {
		return fmt.Errorf("%w: could not create book %v", errs.ErrDBOperation, result.Error)
	}

	return nil
}

func (r *bookRepository) UpdateBook(userID, bookID uint, book models.Book) error {
	err := r.db.Transaction(func(tx *gorm.DB) error {
		result := tx.Model(&models.Book{}).
			Where(&models.Book{ID: bookID, UserID: userID}).
			Select("Title", "Author", "Price").
			Updates(book)

		if result.Error != nil {
			return fmt.Errorf("%w: cound not update book %v", errs.ErrDBOperation, result.Error)
		}
		if result.RowsAffected == 0 {
			return errs.ErrNotFound
		}

		return nil
	})

	return err
}

func (r *bookRepository) DeleteBook(userID, bookID uint) error {
	err := r.db.Transaction(func(tx *gorm.DB) error {
		result := tx.Unscoped().
			Where(&models.Book{ID: bookID, UserID: userID}).
			Delete(&models.Book{})

		if result.Error != nil {
			return fmt.Errorf("%w: could not delete book %v", errs.ErrDBOperation, result.Error)
		}
		if result.RowsAffected == 0 {
			return errs.ErrNotFound
		}

		return nil
	})

	return err
}

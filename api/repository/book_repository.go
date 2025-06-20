package repository

import (
	"bookstore-api/internal/lib/errs"
	"bookstore-api/internal/models"
	"fmt"

	"gorm.io/gorm"
)

type BookRepository interface {
	GetAllBooks() ([]models.Book, error)
	GetUserBooks(
		userID uint,
		author string,
		title string,
		limit int,
	) ([]models.Book, models.KafkaError)
	PostBook(book models.Book) models.KafkaError
	UpdateBook(userID uint, bookID uint, newBook models.Book) models.KafkaError
	DeleteBook(userID uint, bookID uint) models.KafkaError
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
) ([]models.Book, models.KafkaError) {
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
		return nil, models.KafkaError{
			Error:   errs.ErrDBOperation.Error(),
			Message: result.Error.Error(),
		}
	}

	if result.RowsAffected == 0 {
		return nil, models.KafkaError{
			Error: errs.ErrNotFound.Error(),
		}
	}

	return books, models.KafkaError{}
}

func (r *bookRepository) PostBook(book models.Book) models.KafkaError {
	result := r.db.Create(&book)

	if result.Error != nil {
		return models.KafkaError{
			Error:   errs.ErrDBOperation.Error(),
			Message: fmt.Sprintf("could not create book %v", result.Error),
		}
	}

	return models.KafkaError{}
}

func (r *bookRepository) UpdateBook(userID, bookID uint, book models.Book) models.KafkaError {
	var Error, Message string

	err := r.db.Transaction(func(tx *gorm.DB) error {
		result := tx.Model(&models.Book{}).
			Where(&models.Book{ID: bookID, UserID: userID}).
			Select("Title", "Author", "Price").
			Updates(book)

		if result.Error != nil {
			Error = errs.ErrDBOperation.Error()
			Message = "could not update book " + result.Error.Error()
			return fmt.Errorf("%w: cound not update book %v", errs.ErrDBOperation, result.Error)
		}
		if result.RowsAffected == 0 {
			Error = errs.ErrNotFound.Error()
			return errs.ErrNotFound
		}

		return nil
	})

	if err != nil {
		return models.KafkaError{
			Error:   Error,
			Message: Message,
		}
	}

	return models.KafkaError{}
}

func (r *bookRepository) DeleteBook(userID, bookID uint) models.KafkaError {
	var Error, Message string

	err := r.db.Transaction(func(tx *gorm.DB) error {
		result := tx.Unscoped().
			Where(&models.Book{ID: bookID, UserID: userID}).
			Delete(&models.Book{})

		if result.Error != nil {
			Error = errs.ErrDBOperation.Error()
			Message = "could not update book " + result.Error.Error()
			return fmt.Errorf("%w: cound not update book %v", errs.ErrDBOperation, result.Error)
		}
		if result.RowsAffected == 0 {
			Error = errs.ErrNotFound.Error()
			return errs.ErrNotFound
		}

		return nil
	})

	if err != nil {
		return models.KafkaError{
			Error:   Error,
			Message: Message,
		}
	}

	return models.KafkaError{}
}

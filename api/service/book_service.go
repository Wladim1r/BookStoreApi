package service

import (
	"bookstore-api/api/repository"
	"bookstore-api/internal/models"
	"errors"
)

type BookService interface {
	GetAllBooks() ([]models.Book, error)
	GetUserBook(uint, uint) (models.Book, error)
	GetUserBooks(uint) ([]models.Book, error)
	PostBook(models.Book) (*models.Book, error)
	UpdateBook(uint, uint, models.Book) error
	DeleteBook(uint, uint) error
}

type bookService struct {
	repo repository.BookRepository
}

func NewBookService(r repository.BookRepository) BookService {
	return &bookService{repo: r}
}

func (s *bookService) GetAllBooks() ([]models.Book, error) {
	return s.repo.GetAllBooks()
}

func (s *bookService) GetUserBook(userID, bookID uint) (models.Book, error) {
	return s.repo.GetUserBook(userID, bookID)
}

func (s *bookService) GetUserBooks(userID uint) ([]models.Book, error) {
	return s.repo.GetUserBooks(userID)
}

func (s *bookService) PostBook(book models.Book) (*models.Book, error) {
	return s.repo.PostBook(book)

}

func (s *bookService) UpdateBook(userID, bookID uint, book models.Book) error {
	if _, err := s.GetUserBook(userID, bookID); err != nil {
		return errors.New("Book does not exist")
	}

	return s.repo.UpdateBook(userID, bookID, book)
}

func (s *bookService) DeleteBook(userID, bookID uint) error {
	if _, err := s.GetUserBook(userID, bookID); err != nil {
		return errors.New("Book does not exist")
	}

	return s.repo.DeleteBook(userID, bookID)
}

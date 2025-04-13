package service

import (
	"bookstore-api/internal/models"
	"errors"
	"strconv"
)

type BookService interface {
	GetBook(string) (models.Book, error)
	GetBooks() ([]models.Book, error)
	PostBook(models.Book) error
	UpdateBook(string, models.Book) error
	DeleteBook(string) error
}

type bookService struct {
	repo BookRepository
}

func NewBookService(r BookRepository) BookService {
	return &bookService{repo: r}
}

func (s *bookService) GetBook(id_str string) (models.Book, error) {
	id, err := strconv.Atoi(id_str)
	if err != nil {
		return models.Book{}, err
	}

	if id <= 0 {
		return models.Book{}, errors.New("Invalid id")
	}

	return s.repo.GetBook(uint(id))
}

func (s *bookService) GetBooks() ([]models.Book, error) {
	return s.repo.GetBooks()
}

func (s *bookService) PostBook(book models.Book) error {
	return s.repo.PostBook(book)
}

func (s *bookService) UpdateBook(id string, book models.Book) error {
	if _, err := s.GetBook(id); err != nil {
		return errors.New("Invalid id")
	}

	return s.repo.UpdateBook(id, book)
}

func (s *bookService) DeleteBook(id_str string) error {
	id, err := strconv.Atoi(id_str)
	if err != nil {
		return err
	}

	if id <= 0 {
		return errors.New("Invalid id")
	}

	return s.repo.DeleteBook(uint(id))
}

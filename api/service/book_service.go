package service

import (
	"bookstore-api/api/repository"
	"bookstore-api/internal/lib/errs"
	"bookstore-api/internal/models"
	"fmt"
	"strconv"
)

type BookService interface {
	GetAllBooks() ([]models.UserBooksResponse, error)
	GetUserBooks(interface{}, string, string, string) ([]models.Book, uint, error)
	PostBook(interface{}, models.BookRequest) error
	UpdateBook(interface{}, string, models.BookRequest) error
	DeleteBook(interface{}, string) error
}

type bookService struct {
	repo repository.BookRepository
}

func NewBookService(r repository.BookRepository) BookService {
	return &bookService{repo: r}
}

func (s *bookService) GetAllBooks() ([]models.UserBooksResponse, error) {
	books, err := s.repo.GetAllBooks()
	if err != nil {
		return nil, err
	}

	usersMap := make(map[string]*models.UserBooksResponse)

	for _, book := range books {
		username := book.User.Username

		if _, exists := usersMap[username]; !exists {
			usersMap[username] = &models.UserBooksResponse{
				Username:   username,
				TotalBooks: 0,
				Books:      []models.BookResponse{},
			}
		}

		usersMap[username].Books = append(usersMap[username].Books, models.BookResponse{
			ID:     book.ID,
			Title:  book.Title,
			Author: book.Author,
			Price:  book.Price,
		})
		usersMap[username].TotalBooks++
	}

	result := make([]models.UserBooksResponse, 0, len(usersMap))
	for _, userBooks := range usersMap {
		result = append(result, *userBooks)
	}

	return result, nil
}

func (s *bookService) GetUserBooks(
	userID_iface interface{},
	author, title, limitStr string,
) ([]models.Book, uint, error) {
	userID := interface_into_uint(userID_iface)

	var limit int

	if limitStr == "" {
		limit = 0
	} else {
		var err error
		limit, err = strconv.Atoi(limitStr)
		if err != nil || limit <= 0 {
			return nil, 0, fmt.Errorf("%w: %v", errs.ErrInvalidParam, err)
		}
	}

	books, err := s.repo.GetUserBooks(userID, author, title, limit)
	if err != nil {
		return nil, 0, err
	}

	return books, userID, nil
}

func (s *bookService) PostBook(
	userID_iface interface{},
	input models.BookRequest,
) error {
	userID := interface_into_uint(userID_iface)

	book := models.Book{
		Title:  input.Title,
		Author: input.Author,
		Price:  input.Price,
		UserID: userID,
	}

	return s.repo.PostBook(book)

}

func (s *bookService) UpdateBook(
	userID_iface interface{},
	bookIDStr string,
	input models.BookRequest,
) error {
	userID := interface_into_uint(userID_iface)

	bookID, err := strconv.Atoi(bookIDStr)
	if err != nil || bookID <= 0 {
		return fmt.Errorf("%w: invalid ID in request %v", errs.ErrInvalidID, err)
	}

	book := models.Book{
		Title:  input.Title,
		Author: input.Author,
		Price:  input.Price,
		UserID: userID,
	}

	return s.repo.UpdateBook(userID, uint(bookID), book)
}

func (s *bookService) DeleteBook(userID_iface interface{}, bookIDStr string) error {
	userID := interface_into_uint(userID_iface)

	bookID, err := strconv.Atoi(bookIDStr)
	if err != nil || bookID <= 0 {
		return fmt.Errorf("%w: invalid ID in request %v", errs.ErrInvalidID, err)
	}

	return s.repo.DeleteBook(userID, uint(bookID))
}

func interface_into_uint(userID_iface interface{}) uint {
	var userID uint
	switch v := userID_iface.(type) {
	case float64:
		userID = uint(v)
	case int:
		userID = uint(v)
	case uint:
		userID = v
	}
	return userID
}

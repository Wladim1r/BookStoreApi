package service

import (
	"bookstore-api/api/repository"
	"bookstore-api/internal/lib/errs"
	"bookstore-api/internal/models"
	"encoding/json"
	"fmt"
	"strconv"
	"sync"
	"time"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/google/uuid"
)

type BookService interface {
	GetAllBooks() ([]models.UserBooksResponse, error)
	GetUserBooks(interface{}, string, string, string) ([]models.Book, uint, error)
	PostBook(interface{}, models.BookRequest) error
	UpdateBook(interface{}, string, models.BookRequest) error
	DeleteBook(interface{}, string) error
}

type bookService struct {
	repo      repository.BookRepository
	producer  *kafka.Producer
	consumer  *kafka.Consumer
	topic     string
	responses sync.Map
}

func NewBookService(
	r repository.BookRepository,
	p *kafka.Producer,
	c *kafka.Consumer,
	t string,
) BookService {
	s := &bookService{
		repo:     r,
		producer: p,
		consumer: c,
		topic:    t,
	}

	go s.consumptionMessage()

	return s
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

	relID := uuid.New().String()
	ch := make(chan models.KafkaBookResponse)
	s.responses.Store(relID, ch)

	defer func() {
		s.responses.Delete(relID)
		close(ch)
	}()

	payload := models.GetUserBooksRequest{
		UserID: userID,
		Author: author,
		Title:  title,
		Limit:  limit,
	}

	rawMes, err := json.Marshal(payload)
	if err != nil {
		return nil, 0, fmt.Errorf("%w: %v", errs.ErrInternal, err)
	}

	request := models.KafkaBookRequest{
		Method:     getUserBooksMethod,
		Type:       requestType,
		RelationID: relID,
		Payload:    json.RawMessage(rawMes),
	}
	requestBytes, err := json.Marshal(request)
	if err != nil {
		return nil, 0, fmt.Errorf("%w: %v", errs.ErrInternal, err)
	}

	if err := s.sendKafkaRequest(requestBytes); err != nil {
		return nil, 0, err
	}

	select {
	case resp := <-ch:
		if resp.Error.Error != "" {
			if resp.Error.Error == errs.ErrDBOperation.Error() {
				return nil, 0, fmt.Errorf("%w: %v", errs.ErrDBOperation, resp.Error.Message)
			}
			if resp.Error.Error == errs.ErrNotFound.Error() {
				return nil, 0, fmt.Errorf("%w", errs.ErrNotFound)
			}
		}

		var r models.GetUserBooksResponse
		if err := json.Unmarshal(resp.Result, &r); err != nil {
			return nil, 0, fmt.Errorf("%w: %v", errs.ErrInternal, err)
		}

		return r.Books, userID, nil
	case <-time.After(30 * time.Second):
		return nil, 0, errs.ErrTimeout
	}
}

func (s *bookService) PostBook(
	userID_iface interface{},
	input models.BookRequest,
) error {
	userID := interface_into_uint(userID_iface)

	relID := uuid.New().String()
	ch := make(chan models.KafkaBookResponse)
	s.responses.Store(relID, ch)

	defer func() {
		close(ch)
		s.responses.Delete(relID)
	}()

	book := models.Book{
		Title:  input.Title,
		Author: input.Author,
		Price:  input.Price,
		UserID: userID,
	}
	rawMes, err := json.Marshal(book)
	if err != nil {
		return fmt.Errorf("%w: %v", errs.ErrInternal, err)
	}

	request := models.KafkaBookRequest{
		Method:     postBookMethod,
		Type:       requestType,
		RelationID: relID,
		Payload:    json.RawMessage(rawMes),
	}
	requestBytes, err := json.Marshal(request)
	if err != nil {
		return fmt.Errorf("%w: %v", errs.ErrInternal, err)
	}

	if err := s.sendKafkaRequest(requestBytes); err != nil {
		return err
	}

	select {
	case resp := <-ch:
		if resp.Error.Error != "" {
			if resp.Error.Error == errs.ErrDBOperation.Error() {
				return fmt.Errorf("%w: %v", errs.ErrDBOperation, resp.Error.Message)
			}
		}
		return nil
	case <-time.After(30 * time.Second):
		return errs.ErrTimeout
	}
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

	relID := uuid.New().String()
	ch := make(chan models.KafkaBookResponse)
	s.responses.Store(relID, ch)

	defer func() {
		s.responses.Delete(relID)
		close(ch)
	}()

	book := models.Book{
		ID:     uint(bookID),
		Title:  input.Title,
		Author: input.Author,
		Price:  input.Price,
		UserID: userID,
	}
	rawMes, err := json.Marshal(book)
	if err != nil {
		return fmt.Errorf("%w: %v", errs.ErrInternal, err)
	}

	request := models.KafkaBookRequest{
		Method:     updateBookMethod,
		Type:       requestType,
		RelationID: relID,
		Payload:    json.RawMessage(rawMes),
	}
	requestBytes, err := json.Marshal(request)
	if err != nil {
		return fmt.Errorf("%w: %v", errs.ErrInternal, err)
	}

	if err := s.sendKafkaRequest(requestBytes); err != nil {
		return err
	}

	select {
	case resp := <-ch:
		if resp.Error.Error != "" {
			if resp.Error.Error == errs.ErrDBOperation.Error() {
				return fmt.Errorf("%w: %v", errs.ErrDBOperation, resp.Error.Message)
			}
			if resp.Error.Error == errs.ErrNotFound.Error() {
				return fmt.Errorf("%w", errs.ErrNotFound)
			}
		}
		return nil
	case <-time.After(30 * time.Second):
		return errs.ErrTimeout
	}

}

func (s *bookService) DeleteBook(userID_iface interface{}, bookIDStr string) error {
	userID := interface_into_uint(userID_iface)

	bookID, err := strconv.Atoi(bookIDStr)
	if err != nil || bookID <= 0 {
		return fmt.Errorf("%w: invalid ID in request %v", errs.ErrInvalidID, err)
	}

	relID := uuid.New().String()
	ch := make(chan models.KafkaBookResponse)
	s.responses.Store(relID, ch)

	defer func() {
		s.responses.Delete(relID)
		close(ch)
	}()

	req := models.DeleteBook{
		ID:     uint(bookID),
		UserID: userID,
	}
	rawMes, err := json.Marshal(req)
	if err != nil {
		return fmt.Errorf("%w: %v", errs.ErrInternal, err)
	}

	request := models.KafkaBookRequest{
		Method:     deleteBookMethod,
		Type:       requestType,
		RelationID: relID,
		Payload:    json.RawMessage(rawMes),
	}
	requestBytes, err := json.Marshal(request)
	if err != nil {
		return fmt.Errorf("%w: %v", errs.ErrInternal, err)
	}

	if err := s.sendKafkaRequest(requestBytes); err != nil {
		return err
	}

	select {
	case resp := <-ch:
		if resp.Error.Error != "" {
			if resp.Error.Error == errs.ErrDBOperation.Error() {
				return fmt.Errorf("%w: %v", errs.ErrDBOperation, resp.Error.Message)
			}
			if resp.Error.Error == errs.ErrNotFound.Error() {
				return fmt.Errorf("%w", errs.ErrNotFound)
			}
		}
		return nil
	case <-time.After(30 * time.Second):
		return errs.ErrTimeout
	}
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

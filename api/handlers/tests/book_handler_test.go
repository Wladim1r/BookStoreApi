package tests

import (
	"bookstore-api/api/handlers"
	"bookstore-api/internal/models"
	"bookstore-api/internal/utils"
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockBookService struct {
	mock.Mock
}

func (m *MockBookService) GetAllBooks() ([]models.Book, error) {
	args := m.Called()
	return args.Get(0).([]models.Book), args.Error(1)
}

func (m *MockBookService) GetUserBook(userID, bookID uint) (models.Book, error) {
	args := m.Called(userID, bookID)
	return args.Get(0).(models.Book), args.Error(1)
}

func (m *MockBookService) GetUserBooks(userID uint) ([]models.Book, error) {
	args := m.Called(userID)
	return args.Get(0).([]models.Book), args.Error(1)
}

func (m *MockBookService) PostBook(book models.Book) (*models.Book, error) {
	args := m.Called(book)
	return args.Get(0).(*models.Book), args.Error(1)
}

func (m *MockBookService) UpdateBook(userID, bookID uint, book models.Book) error {
	args := m.Called(userID, bookID, book)
	return args.Error(0)
}

func (m *MockBookService) DeleteBook(userID, bookID uint) error {
	args := m.Called(userID, bookID)
	return args.Error(0)
}

func TestBookHandler_GetUserBook(t *testing.T) {
	mockService := new(MockBookService)
	handler := &handlers.BookHandler{Service: mockService}

	t.Run("Successful test", func(t *testing.T) {
		testUserID := uint(13)
		testBookID := uint(3)
		testBook := models.Book{Title: "test", Author: "name", Price: 123}

		mockService.On("GetUserBook", testUserID, testBookID).Return(testBook, nil).Once()

		responseRecorder := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(responseRecorder)
		c.Set("userID", float64(testUserID))

		c.Request = httptest.NewRequest(http.MethodGet, "/api/books/3", nil)
		c.Params = gin.Params{gin.Param{Key: "id", Value: "3"}}

		handler.GetUserBook(c)

		assert.Equal(t, http.StatusOK, responseRecorder.Code)
		utils.Success_GetUserBook(t,
			map[string]interface{}{"title": "test", "author": "name", "price": 123},
			responseRecorder.Body.String())

		mockService.AssertExpectations(t)
	})

	t.Run("Failed UserID", func(t *testing.T) {
		responseRecorder := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(responseRecorder)

		handler.GetUserBooks(c)

		assert.Equal(t, http.StatusUnauthorized, responseRecorder.Code)
		utils.WrongID(t, responseRecorder.Body.String())
	})

	t.Run("Missing ID in params", func(t *testing.T) {
		testUserID := uint(8)

		responseRecorder := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(responseRecorder)
		c.Set("userID", float64(testUserID))

		c.Params = gin.Params{gin.Param{Key: "id", Value: ""}}

		handler.GetUserBook(c)

		assert.Equal(t, http.StatusBadRequest, responseRecorder.Code)
		utils.WrongParamID(t, responseRecorder.Body.String())
	})

	t.Run("Negative ID in params", func(t *testing.T) {
		testUserID := uint(9)

		responseRecorder := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(responseRecorder)
		c.Set("userID", float64(testUserID))

		c.Params = gin.Params{gin.Param{Key: "id", Value: "-1"}}

		handler.GetUserBook(c)

		assert.Equal(t, http.StatusBadRequest, responseRecorder.Code)
		utils.WrongParamID(t, responseRecorder.Body.String())
	})

	t.Run("Incorrect ID in params", func(t *testing.T) {
		testUserID := uint(9)

		responseRecorder := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(responseRecorder)
		c.Set("userID", float64(testUserID))

		c.Params = gin.Params{gin.Param{Key: "id", Value: "abc"}}

		handler.GetUserBook(c)

		assert.Equal(t, http.StatusBadRequest, responseRecorder.Code)
		utils.WrongParamID(t, responseRecorder.Body.String())
	})

	t.Run("Not found book", func(t *testing.T) {
		testUserID := uint(4)
		nonExistentBookID := uint(999)

		mockService.On("GetUserBook", testUserID, nonExistentBookID).
			Return(models.Book{}, errors.New("Book with entred ID does not exist")).Once()

		responseRecorder := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(responseRecorder)
		c.Set("userID", float64(testUserID))

		c.Params = gin.Params{gin.Param{Key: "id", Value: "999"}}

		handler.GetUserBook(c)

		assert.Equal(t, http.StatusNotFound, responseRecorder.Code)
		utils.NotFound(t, responseRecorder.Body.String())
	})
}

func TestBookHandler_GetUserBooks(t *testing.T) {
	mockService := new(MockBookService)
	handler := &handlers.BookHandler{Service: mockService}

	t.Run("Successful test", func(t *testing.T) {
		testUserID := uint(13)
		testBooks := []models.Book{
			{Title: "test 1", Author: "name 1", Price: 111},
			{Title: "test 2", Author: "name 2", Price: 222},
		}

		mockService.On("GetUserBooks", testUserID).Return(testBooks, nil).Once()

		responseRecorder := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(responseRecorder)
		c.Set("userID", float64(testUserID))

		handler.GetUserBooks(c)

		utils.Success_GetUserBooks(t, map[string]interface{}{
			"data": []interface{}{
				map[string]interface{}{"title": "test 1", "author": "name 1", "price": 111},
				map[string]interface{}{"title": "test 2", "author": "name 2", "price": 222},
			},
			"meta": map[string]interface{}{
				"total":   len(testBooks),
				"user_id": testUserID,
			},
		}, responseRecorder.Body.String())

		mockService.AssertExpectations(t)
	})

	t.Run("Failed UserID", func(t *testing.T) {
		responseRecorder := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(responseRecorder)

		handler.GetUserBooks(c)

		assert.Equal(t, http.StatusUnauthorized, responseRecorder.Code)
		utils.WrongID(t, responseRecorder.Body.String())
	})

	t.Run("Service Error", func(t *testing.T) {
		UserID := uint(2)

		mockService.On("GetUserBooks", UserID).Return([]models.Book{},
			errors.New("Database connection failed")).Once()

		responseRecorder := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(responseRecorder)
		c.Set("userID", float64(UserID))

		handler.GetUserBooks(c)

		assert.Equal(t, http.StatusInternalServerError, responseRecorder.Code)
		utils.ServiceError(t, responseRecorder.Body.String())
	})
}

func TestBookHandler_PostBook(t *testing.T) {
	mockService := new(MockBookService)
	handler := handlers.BookHandler{Service: mockService}

	t.Run("Successful test", func(t *testing.T) {
		testUserID := uint(4)
		testBook := models.Book{
			Title:  "test",
			Author: "name",
			Price:  369,
		}

		bookJson, err := json.Marshal(testBook)
		assert.NoError(t, err)

		expectedBook := models.Book{
			Title:  testBook.Title,
			Author: testBook.Author,
			Price:  testBook.Price,
			UserID: testUserID,
		}

		mockService.On("PostBook", expectedBook).Return(&expectedBook, nil).Once()

		responseRecorder := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(responseRecorder)
		c.Set("userID", testUserID)

		c.Request = httptest.NewRequest(http.MethodPost, "/api/books", bytes.NewBuffer(bookJson))
		c.Request.Header.Set("Content-Type", "application/json")

		handler.PostBook(c)

		assert.Equal(t, http.StatusCreated, responseRecorder.Code)
		utils.Success_PostBook(t, map[string]interface{}{
			"title":  "test",
			"author": "name",
			"price":  369,
		}, responseRecorder.Body.String())

		mockService.AssertExpectations(t)
	})

	t.Run("Failed UserID", func(t *testing.T) {
		responseRecorder := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(responseRecorder)

		handler.PostBook(c)

		assert.Equal(t, http.StatusUnauthorized, responseRecorder.Code)
		utils.WrongID(t, responseRecorder.Body.String())

		mockService.AssertNotCalled(t, "PostBook")
	})

	t.Run("Invalid body request", func(t *testing.T) {
		testUserID := uint(10)
		testBook := models.Book{
			Title:  "empty field",
			Author: "",
			Price:  10,
		}
		bookJson, err := json.Marshal(testBook)
		assert.NoError(t, err)

		responseRecorder := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(responseRecorder)
		c.Set("userID", testUserID)

		c.Request = httptest.NewRequest(http.MethodPost, "/api/books", bytes.NewBuffer(bookJson))
		c.Request.Header.Set("Content-Type", "application/json")

		handler.PostBook(c)

		assert.Equal(t, http.StatusBadRequest, responseRecorder.Code)
		utils.InvalidBodyRequest(t, responseRecorder.Body.String())

		mockService.AssertNotCalled(t, "PostBook")
	})

	t.Run("Service Error", func(t *testing.T) {
		testUserID := uint(6)
		testBook := models.Book{
			Title:  "who",
			Author: "ho",
			Price:  101,
		}

		jsonData, err := json.Marshal(testBook)
		assert.NoError(t, err)

		expectedBook := models.Book{
			Title:  testBook.Title,
			Author: testBook.Author,
			Price:  testBook.Price,
			UserID: testUserID,
		}

		mockService.On("PostBook", expectedBook).Return(&models.Book{},
			errors.New("Database connection failed")).Once()

		responseRecorder := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(responseRecorder)
		c.Set("userID", testUserID)

		c.Request = httptest.NewRequest(http.MethodPost, "/api/books", bytes.NewBuffer(jsonData))
		c.Request.Header.Set("Content-Type", "application/json")

		handler.PostBook(c)

		assert.Equal(t, http.StatusInternalServerError, responseRecorder.Code)
		utils.ServiceError(t, responseRecorder.Body.String())

		mockService.AssertExpectations(t)
	})
}

func TestBookHandler_UpdateBook(t *testing.T) {
	mockService := new(MockBookService)
	handler := handlers.BookHandler{Service: mockService}

	t.Run("Successful UpdateBook", func(t *testing.T) {
		testUserID := uint(8)
		testBookID := uint(2)
		testBook := models.Book{
			Title:  "altered test",
			Author: "altered name",
			Price:  78,
		}

		bookJson, err := json.Marshal(testBook)
		assert.NoError(t, err)

		expectedBook := models.Book{
			Title:  testBook.Title,
			Author: testBook.Author,
			Price:  testBook.Price,
			UserID: testUserID,
		}

		mockService.On("UpdateBook", testUserID, testBookID, expectedBook).Return(nil).Once()

		responseRecorder := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(responseRecorder)

		c.Set("userID", testUserID)
		c.Params = gin.Params{gin.Param{Key: "id", Value: "2"}}
		c.Request = httptest.NewRequest(http.MethodPatch, "/api/books/2", bytes.NewBuffer(bookJson))
		c.Request.Header.Set("Content-Type", "application/json")

		handler.UpdateBook(c)

		assert.Equal(t, http.StatusOK, responseRecorder.Code)
		utils.Success_UpdateBook(t, responseRecorder.Body.String())
	})

	t.Run("Failed UserID", func(t *testing.T) {
		responseRecorder := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(responseRecorder)

		handler.UpdateBook(c)

		assert.Equal(t, http.StatusUnauthorized, responseRecorder.Code)
		utils.WrongID(t, responseRecorder.Body.String())
	})

	t.Run("Missing ID in params", func(t *testing.T) {
		testUserID := uint(8)

		responseRecorder := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(responseRecorder)
		c.Set("userID", testUserID)

		c.Params = gin.Params{gin.Param{Key: "id", Value: ""}}

		handler.UpdateBook(c)

		assert.Equal(t, http.StatusBadRequest, responseRecorder.Code)
		utils.WrongParamID(t, responseRecorder.Body.String())
	})

	t.Run("Negative ID in params", func(t *testing.T) {
		testUserID := uint(9)

		responseRecorder := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(responseRecorder)
		c.Set("userID", float64(testUserID))

		c.Params = gin.Params{gin.Param{Key: "id", Value: "-1"}}

		handler.UpdateBook(c)

		assert.Equal(t, http.StatusBadRequest, responseRecorder.Code)
		utils.WrongParamID(t, responseRecorder.Body.String())
	})

	t.Run("Incorrect ID in params", func(t *testing.T) {
		testUserID := uint(9)

		responseRecorder := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(responseRecorder)
		c.Set("userID", float64(testUserID))

		c.Params = gin.Params{gin.Param{Key: "id", Value: "abc"}}

		handler.UpdateBook(c)

		assert.Equal(t, http.StatusBadRequest, responseRecorder.Code)
		utils.WrongParamID(t, responseRecorder.Body.String())
	})

	t.Run("Invalid body request", func(t *testing.T) {
		testUserID := uint(11)
		testBook := struct {
			Title  string  `json:"title"`
			Author int     `json:"author"`
			Price  float32 `json:"price"`
		}{
			Title:  "wrong fields",
			Author: 64,
			Price:  123.321,
		}

		jsonBook, err := json.Marshal(testBook)
		assert.NoError(t, err)

		responseRecorder := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(responseRecorder)

		c.Set("userID", testUserID)
		c.Params = gin.Params{gin.Param{Key: "id", Value: "11"}}
		c.Request = httptest.NewRequest(http.MethodPatch, "/api/books/11", bytes.NewBuffer(jsonBook))
		c.Request.Header.Set("Content-Type", "application/json")

		handler.PostBook(c)

		assert.Equal(t, http.StatusBadRequest, responseRecorder.Code)
		utils.InvalidBodyRequest(t, responseRecorder.Body.String())
	})

	t.Run("Service Error", func(t *testing.T) {
		testUserID := uint(2)
		testBookID := uint(12)
		testBook := models.Book{
			Title:  "test",
			Author: "name",
			Price:  14,
		}

		jsonBook, err := json.Marshal(testBook)
		assert.NoError(t, err)

		expectedBook := models.Book{
			Title:  testBook.Title,
			Author: testBook.Author,
			Price:  testBook.Price,
			UserID: testUserID,
		}

		mockService.On("UpdateBook", testUserID, testBookID, expectedBook).
			Return(errors.New("Database connection failed")).Once()

		responseRecorder := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(responseRecorder)

		c.Set("userID", testUserID)
		c.Params = gin.Params{gin.Param{Key: "id", Value: "12"}}
		c.Request = httptest.NewRequest(http.MethodPatch, "/api/books/12", bytes.NewBuffer(jsonBook))
		c.Request.Header.Set("Content-Type", "application/json")

		handler.UpdateBook(c)

		assert.Equal(t, http.StatusInternalServerError, responseRecorder.Code)
		utils.ServiceError(t, responseRecorder.Body.String())
	})
}

func TestBookHandler_DeleteBook(t *testing.T) {
	mockService := new(MockBookService)
	handler := handlers.BookHandler{Service: mockService}

	t.Run("Successful DeleteBook", func(t *testing.T) {
		testUserID := uint(15)
		testBookID := uint(5)

		mockService.On("DeleteBook", testUserID, testBookID).Return(nil).Once()

		responseRecorder := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(responseRecorder)

		c.Set("userID", testUserID)
		c.Params = gin.Params{gin.Param{Key: "id", Value: "5"}}

		handler.DeleteBook(c)

		assert.Equal(t, http.StatusOK, responseRecorder.Code)
		utils.Success_DeleteBook(t, responseRecorder.Body.String())

		mockService.AssertExpectations(t)
	})

	t.Run("Failed UserID", func(t *testing.T) {
		responseRecorder := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(responseRecorder)

		handler.DeleteBook(c)

		assert.Equal(t, http.StatusUnauthorized, responseRecorder.Code)
		utils.WrongID(t, responseRecorder.Body.String())
	})

	t.Run("Missing ID in params", func(t *testing.T) {
		testUserID := uint(8)

		responseRecorder := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(responseRecorder)
		c.Set("userID", testUserID)

		c.Params = gin.Params{gin.Param{Key: "id", Value: ""}}

		handler.DeleteBook(c)

		assert.Equal(t, http.StatusBadRequest, responseRecorder.Code)
		utils.WrongParamID(t, responseRecorder.Body.String())
	})

	t.Run("Negative ID in params", func(t *testing.T) {
		testUserID := uint(9)

		responseRecorder := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(responseRecorder)
		c.Set("userID", float64(testUserID))

		c.Params = gin.Params{gin.Param{Key: "id", Value: "-1"}}

		handler.DeleteBook(c)

		assert.Equal(t, http.StatusBadRequest, responseRecorder.Code)
		utils.WrongParamID(t, responseRecorder.Body.String())
	})

	t.Run("Incorrect ID in params", func(t *testing.T) {
		testUserID := uint(9)

		responseRecorder := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(responseRecorder)
		c.Set("userID", float64(testUserID))

		c.Params = gin.Params{gin.Param{Key: "id", Value: "abc"}}

		handler.DeleteBook(c)

		assert.Equal(t, http.StatusBadRequest, responseRecorder.Code)
		utils.WrongParamID(t, responseRecorder.Body.String())
	})

	t.Run("Service Error", func(t *testing.T) {
		testUserID := uint(1)
		testBookID := uint(1)

		mockService.On("DeleteBook", testUserID, testBookID).
			Return(errors.New("Database connection failed")).Once()

		responseRecorder := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(responseRecorder)

		c.Set("userID", testUserID)
		c.Params = gin.Params{gin.Param{Key: "id", Value: "1"}}

		handler.DeleteBook(c)

		assert.Equal(t, http.StatusInternalServerError, responseRecorder.Code)
		utils.ServiceError(t, responseRecorder.Body.String())
	})
}

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
	"gorm.io/gorm"
)

type MockUserService struct {
	mock.Mock
}

func (m *MockUserService) CreateUser(user models.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockUserService) GetAllUsers() ([]models.User, error) {
	args := m.Called()
	return args.Get(0).([]models.User), args.Error(1)
}

func (m *MockUserService) GetByUsername(username string) (models.User, error) {
	args := m.Called(username)
	return args.Get(0).(models.User), args.Error(1)
}

func (m *MockUserService) DeleteByUsername(username string) error {
	args := m.Called(username)
	return args.Error(0)
}

func TestUserHandler_Register(t *testing.T) {
	mockService := new(MockUserService)
	handler := &handlers.UserHandler{Service: mockService}

	t.Run("Successful register", func(t *testing.T) {
		testRequest := models.Request{
			Username: "test name",
			Password: "qwerty",
		}

		jsonBody, err := json.Marshal(testRequest)
		assert.NoError(t, err)

		mockService.On("CreateUser", mock.MatchedBy(func(u models.User) bool {
			return u.Username == testRequest.Username &&
				u.Password != testRequest.Password &&
				len(u.Password) > 50
		})).Return(nil).Once()

		responseRecorder := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(responseRecorder)

		c.Request = httptest.NewRequest(http.MethodPost, "/api/register", bytes.NewBuffer(jsonBody))

		handler.Register(c)

		assert.Equal(t, http.StatusCreated, responseRecorder.Code)
		utils.Success_Register(t, responseRecorder.Body.String())

		mockService.AssertExpectations(t)
	})

	t.Run("Invalid body request", func(t *testing.T) {
		testRequest := models.Request{
			Username: "",
			Password: "empty field",
		}

		jsonBody, err := json.Marshal(testRequest)
		assert.NoError(t, err)

		responseRecorder := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(responseRecorder)

		c.Request = httptest.NewRequest(http.MethodPost, "/api/register", bytes.NewBuffer(jsonBody))

		handler.Register(c)

		assert.Equal(t, http.StatusBadRequest, responseRecorder.Code)
		utils.InvalidBodyRequest(t, responseRecorder.Body.String())
	})

	t.Run("Failed to hash password", func(t *testing.T) {
		testRequest := models.Request{
			Username: "too long password",
			Password: "qwertyuiopasdfghjklzxcvbnmqwertyuiopasdfghjklzxcvbnm12345678901234567890098765432211234567890098765432112345678900987654321",
		}

		jsonBody, err := json.Marshal(testRequest)
		assert.NoError(t, err)

		responseRecorder := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(responseRecorder)

		c.Request = httptest.NewRequest(http.MethodPost, "/api/register", bytes.NewBuffer(jsonBody))

		handler.Register(c)

		assert.Equal(t, http.StatusInternalServerError, responseRecorder.Code)
		utils.HashPasswordError(t, responseRecorder.Body.String())
	})

	t.Run("Database connection failed", func(t *testing.T) {
		testRequest := models.Request{
			Username: "name",
			Password: "1234",
		}

		jsonBody, err := json.Marshal(testRequest)
		assert.NoError(t, err)

		mockService.On("CreateUser", mock.MatchedBy(func(u models.User) bool {
			return u.Username == testRequest.Username &&
				u.Password != testRequest.Password &&
				len(u.Password) > 50
		})).Return(errors.New("Database connection failed")).Once()

		responseRecorder := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(responseRecorder)

		c.Request = httptest.NewRequest(http.MethodPost, "/api/register", bytes.NewBuffer(jsonBody))

		handler.Register(c)

		assert.Equal(t, http.StatusInternalServerError, responseRecorder.Code)
		utils.ServiceError(t, responseRecorder.Body.String())
	})
}

func TestUserHandler_Login(t *testing.T) {
	mockService := new(MockUserService)
	handler := &handlers.UserHandler{Service: mockService}

	t.Run("Successful login", func(t *testing.T) {
		testRequest := models.Request{
			Username: "test",
			Password: "correct",
		}

		jsonBody, err := json.Marshal(testRequest)
		assert.NoError(t, err)

		testUser := models.User{
			Model:    gorm.Model{ID: 1},
			Username: testRequest.Username,
			Password: testRequest.Password,
		}

		err = testUser.HashPassword()
		assert.NoError(t, err)

		mockService.On("GetByUsername", testRequest.Username).Return(testUser, nil).Once()

		responseRecorder := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(responseRecorder)

		c.Request = httptest.NewRequest(http.MethodPost, "/api/login", bytes.NewBuffer(jsonBody))

		handler.Login(c)

		assert.Equal(t, http.StatusOK, responseRecorder.Code)
		utils.Success_Login(t, responseRecorder.Body.String())

		mockService.AssertExpectations(t)
	})

	t.Run("Invalid body request", func(t *testing.T) {
		testRequest := models.Request{
			Username: "empty field",
			Password: "",
		}

		jsonBody, err := json.Marshal(testRequest)
		assert.NoError(t, err)

		responseRecorder := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(responseRecorder)

		c.Request = httptest.NewRequest(http.MethodPost, "/api/login", bytes.NewBuffer(jsonBody))

		handler.Login(c)

		assert.Equal(t, http.StatusBadRequest, responseRecorder.Code)
		utils.InvalidBodyRequest(t, responseRecorder.Body.String())
	})

	t.Run("Not registered yet", func(t *testing.T) {
		testRequest := models.Request{
			Username: "test",
			Password: "pppp",
		}

		jsonBody, err := json.Marshal(testRequest)
		assert.NoError(t, err)

		mockService.On("GetByUsername", testRequest.Username).
			Return(models.User{}, errors.New("You have not registered yet")).Once()

		responseRecorder := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(responseRecorder)
		c.Request = httptest.NewRequest(http.MethodPost, "/api/login", bytes.NewBuffer(jsonBody))

		handler.Login(c)

		assert.Equal(t, http.StatusUnauthorized, responseRecorder.Code)
		utils.NotRegistred(t, responseRecorder.Body.String())
	})

	t.Run("Incorrect password", func(t *testing.T) {
		testRequest := models.Request{
			Username: "name",
			Password: "test",
		}

		jsonBody, err := json.Marshal(testRequest)
		assert.NoError(t, err)

		testUser := models.User{
			Model:    gorm.Model{ID: 1},
			Username: "name",
			Password: "no test",
		}
		err = testUser.HashPassword()
		assert.NoError(t, err)

		mockService.On("GetByUsername", testRequest.Username).
			Return(testUser, nil).Once()

		responseRecorder := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(responseRecorder)
		c.Request = httptest.NewRequest(http.MethodPost, "/api/login", bytes.NewBuffer(jsonBody))

		handler.Login(c)

		assert.Equal(t, http.StatusUnauthorized, responseRecorder.Code)
		utils.IncorrectPassword(t, responseRecorder.Body.String())
	})
}

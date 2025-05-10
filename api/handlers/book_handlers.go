package handlers

import (
	"bookstore-api/api/service"
	"bookstore-api/internal/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type BookHandler struct {
	Service service.BookService
}

func NewBookHandler(service service.BookService) *BookHandler {
	return &BookHandler{Service: service}
}

func (b *BookHandler) GetAllBooks(c *gin.Context) {
	books, err := b.Service.GetAllBooks()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid request",
		})
		return
	}

	type bookResponse struct {
		ID     uint   `json:"id"`
		Title  string `json:"title"`
		Author string `json:"author"`
		Price  uint   `json:"price"`
	}

	type userBooksResponse struct {
		Username   string         `json:"username"`
		TotalBooks int            `json:"total_books"`
		Books      []bookResponse `json:"books"`
	}

	usersMap := make(map[string]*userBooksResponse)

	for _, book := range books {
		username := book.User.Username

		if _, exists := usersMap[username]; !exists {
			usersMap[username] = &userBooksResponse{
				Username:   username,
				TotalBooks: 0,
				Books:      []bookResponse{},
			}
		}

		usersMap[username].Books = append(usersMap[username].Books, bookResponse{
			ID:     book.ID,
			Title:  book.Title,
			Author: book.Author,
			Price:  book.Price,
		})
		usersMap[username].TotalBooks++
	}

	result := make([]userBooksResponse, 0, len(usersMap))
	for _, userBooks := range usersMap {
		result = append(result, *userBooks)
	}

	c.JSON(http.StatusOK, result)
}

func (b *BookHandler) GetUserBook(c *gin.Context) {
	userID_iface, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Authentication required",
		})
		return
	}
	userID := interface_into_uint(userID_iface)

	bookIDStr := c.Param("id")
	bookID, err := strconv.Atoi(bookIDStr)
	if err != nil || bookID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid ID in request",
		})
		return
	}

	book, err := b.Service.GetUserBook(userID, uint(bookID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Book with entred ID does not exist",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": book,
	})
}

func (b *BookHandler) GetUserBooks(c *gin.Context) {
	userID_iface, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Authentication required",
		})
		return
	}
	userID := interface_into_uint(userID_iface)

	books, err := b.Service.GetUserBooks(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Database connection failed",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": books,
		"meta": gin.H{
			"total":   len(books),
			"user_id": userID,
		},
	})
}

func (b *BookHandler) PostBook(c *gin.Context) {
	userID_iface, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Authentication required",
		})
		return
	}
	userID := interface_into_uint(userID_iface)

	var input struct {
		Title  string `json:"title" binding:"required"`
		Author string `json:"author" binding:"required"`
		Price  uint   `json:"price" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid body request",
		})
		return
	}

	book := models.Book{
		Title:  input.Title,
		Author: input.Author,
		Price:  input.Price,
		UserID: userID,
	}

	createdBook, err := b.Service.PostBook(book)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Database connection failed",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"Created Book": *createdBook,
	})
}

func (b *BookHandler) UpdateBook(c *gin.Context) {
	userID_iface, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Authentication required",
		})
		return
	}
	userID := interface_into_uint(userID_iface)

	bookIDStr := c.Param("id")
	bookID, err := strconv.Atoi(bookIDStr)
	if err != nil || bookID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid ID in request",
		})
		return
	}

	var input struct {
		Title  string `json:"title" binding:"required"`
		Author string `json:"author" binding:"required"`
		Price  uint   `json:"price" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid body request",
		})
		return
	}

	book := models.Book{
		Title:  input.Title,
		Author: input.Author,
		Price:  input.Price,
		UserID: userID,
	}

	err = b.Service.UpdateBook(userID, uint(bookID), book)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Database connection failed",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Alterations have been done",
	})
}

func (b *BookHandler) DeleteBook(c *gin.Context) {
	userID_iface, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Authentication required",
		})
		return
	}
	userID := interface_into_uint(userID_iface)

	bookIDStr := c.Param("id")
	bookID, err := strconv.Atoi(bookIDStr)
	if err != nil || bookID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid ID in request",
		})
		return
	}

	err = b.Service.DeleteBook(userID, uint(bookID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Database connection failed",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Book was successfully deleted",
	})
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

package api

import (
	"bookstore-api/internal/models"
	"bookstore-api/internal/service"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type BookHandler struct {
	service service.BookService
}

func NewBookHandler(service service.BookService) *BookHandler {
	return &BookHandler{service: service}
}

func (b *BookHandler) GetBooks(c *gin.Context) {
	books, err := b.service.GetBooks()

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid request",
		})
		return
	}

	c.JSON(http.StatusOK, books)
}

func (b *BookHandler) GetBook(c *gin.Context) {
	book, err := b.service.GetBook(c.Param("id"))

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"message": fmt.Sprintf("Book with entred ID does not exist"),
		})
		return
	}

	c.JSON(http.StatusOK, book)
}

func (b *BookHandler) PostBook(c *gin.Context) {
	var book models.Book

	if err := c.ShouldBindJSON(&book); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid body request",
		})
		return
	}

	err := b.service.PostBook(book)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Book with entred ID already exist",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"Created Book": book,
	})
}

func (b *BookHandler) UpdateBook(c *gin.Context) {
	var book models.Book

	if err := c.ShouldBindJSON(&book); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid body request",
		})
		return
	}

	err := b.service.UpdateBook(c.Param("id"), book)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "You do not have permission to make changes Book`s ID",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Alterations have been done",
	})
}

func (b *BookHandler) DeleteBook(c *gin.Context) {

	if _, err := b.service.GetBook(c.Param("id")); err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"message": fmt.Sprintf("Book with entred ID does not exist"),
		})
		return

	}

	err := b.service.DeleteBook(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Could not delete Book",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Book was successfully deleted",
	})
}

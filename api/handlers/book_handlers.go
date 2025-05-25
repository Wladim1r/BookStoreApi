package handlers

import (
	"bookstore-api/api/service"
	"bookstore-api/internal/lib/errs"
	"bookstore-api/internal/models"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

type BookHandler struct {
	Service service.BookService
}

func NewBookHandler(service service.BookService) *BookHandler {
	return &BookHandler{Service: service}
}

// @Summary Get books of all users
// @Description Get books collection for all registered users
// @Tags Admin
// @ID get-all-books
// @Security BasicAuth
// @Accept json
// @Produce json
// @Success 200 {object} models.UsersBooksResponse "List of books of all users"
// @Failure 401 {object} models.ErrorResponse "User unauthorized"
// @Failure 404 {object} models.ErrorResponse "Records not found"
// @Failure 500 {object} models.ErrorResponse "Database or Server error"
// @Router /admin/books [get]
func (b *BookHandler) GetAllBooks(c *gin.Context) {
	books, err := b.Service.GetAllBooks()
	if err != nil {

		switch {
		case errors.Is(err, errs.ErrNotFound):
			c.JSON(http.StatusNotFound, models.ErrorResponse{
				Error: "Could not found records",
			})
		case errors.Is(err, errs.ErrDBOperation):
			c.JSON(http.StatusInternalServerError, models.ErrorResponse{
				Error: "Database operation failed",
			})
		default:
			c.JSON(http.StatusInternalServerError, models.ErrorResponse{
				Error: "Internal server error",
			})
		}

		return
	}

	c.JSON(http.StatusOK, models.UsersBooksResponse{
		Data: books,
	})
}

// @Summary Get books of user
// @Description Get user books collection with parameters
// @Tags Books
// @ID get-user-books
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param author query string false "Filter by author name" example("Пушкин")
// @Param title query string false "Filter by title" example("Я вас любил")
// @Param limit query int false "Limit number of records" minimum(1) example(10)
// @Success 200 {object} models.GetBooks "List of books of user"
// @Failure 400 {object} models.ErrorResponse "Invalid query body"
// @Failure 401 {object} models.ErrorResponse "User unauthorized"
// @Failure 404 {object} models.ErrorResponse "Records not found"
// @Failure 500 {object} models.ErrorResponse "Database or Server error"
// @Router /api/books [get]
func (b *BookHandler) GetUserBooks(c *gin.Context) {
	userID_iface, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, models.ErrorResponse{
			Error: "Authentication required",
		})
		return
	}

	author := c.Query("author")
	title := c.Query("title")
	limitStr := c.Query("limit")

	books, userID, err := b.Service.GetUserBooks(userID_iface, author, title, limitStr)
	if err != nil {

		switch {
		case errors.Is(err, errs.ErrInvalidParam):
			c.JSON(http.StatusBadRequest, models.ErrorResponse{
				Error: "Invalid limit value",
			})
		case errors.Is(err, errs.ErrNotFound):
			c.JSON(http.StatusNotFound, models.ErrorResponse{
				Error: "Could not found records",
			})
		case errors.Is(err, errs.ErrDBOperation):
			c.JSON(http.StatusInternalServerError, models.ErrorResponse{
				Error: "Database operation failed",
			})
		default:
			c.JSON(http.StatusInternalServerError, models.ErrorResponse{
				Error: "Internal server error",
			})
		}

		return
	}

	c.JSON(http.StatusOK, models.GetBooks{
		Data: books,
		Meta: models.MetaBook{
			Total:  len(books),
			UserID: userID,
		},
	})
}

// @Summary Create book
// @Description Create book with users parameters
// @Tags Books
// @ID post-user-book
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param request body models.BookRequest true "Data for create book"
// @Success 201 {object} models.SuccessResponse "Message about successfully creating"
// @Failure 400 {object} models.ErrorResponse "Invalid request body"
// @Failure 401 {object} models.ErrorResponse "User unauthorized"
// @Failure 500 {object} models.ErrorResponse "Database or Server error"
// @Router /api/books [post]
func (b *BookHandler) PostBook(c *gin.Context) {
	userID_iface, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, models.ErrorResponse{
			Error: "Authentication required",
		})
		return
	}

	var input models.BookRequest

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error: "Invalid body request",
		})
		return
	}

	err := b.Service.PostBook(userID_iface, input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error: "Database operation failed",
		})
		return
	}

	c.JSON(http.StatusCreated, models.SuccessResponse{
		Message: "Book was successfully created",
	})
}

// @Summary Update book
// @Description Change book to new parameters
// @Tags Books
// @ID update-user-book
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param id path int true "ID of the book to change" minimum(1) example(13)
// @Param request body models.BookRequest true "New data for change existing data"
// @Success 200 {object} models.SuccessResponse "Message about successfully updating"
// @Failure 400 {object} models.ErrorResponse "Invalid request body"
// @Failure 401 {object} models.ErrorResponse "User unauthorized"
// @Failure 404 {object} models.ErrorResponse "Record not found"
// @Failure 500 {object} models.ErrorResponse "Database or Server error"
// @Router /api/books/{id} [patch]
func (b *BookHandler) UpdateBook(c *gin.Context) {
	userID_iface, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, models.ErrorResponse{
			Error: "Authentication required",
		})
		return
	}

	bookIDStr := c.Param("id")

	var input models.BookRequest

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error: "Invalid body request",
		})
		return
	}

	err := b.Service.UpdateBook(userID_iface, bookIDStr, input)
	if err != nil {

		switch {
		case errors.Is(err, errs.ErrNotFound):
			c.JSON(http.StatusNotFound, models.ErrorResponse{
				Error: "Could not found record",
			})
		case errors.Is(err, errs.ErrDBOperation):
			c.JSON(http.StatusInternalServerError, models.ErrorResponse{
				Error: "Database operation failed",
			})
		default:
			c.JSON(http.StatusInternalServerError, models.ErrorResponse{
				Error: "Internal server error",
			})
		}

		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse{
		Message: "Alterations have been done",
	})
}

// @Summary Delete book
// @Description Permanently delete book
// @Tags Books
// @ID delete-user-book
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param id path int true "ID of the book to delete" minimum(1) example(3)
// @Success 200 {object} models.SuccessResponse "Message about successfully deleting"
// @Failure 400 {object} models.ErrorResponse "Invalid request body"
// @Failure 401 {object} models.ErrorResponse "User unauthorized"
// @Failure 404 {object} models.ErrorResponse "Record not found"
// @Failure 500 {object} models.ErrorResponse "Database or Server error"
// @Router /api/books/{id} [delete]
func (b *BookHandler) DeleteBook(c *gin.Context) {
	userID_iface, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, models.ErrorResponse{
			Error: "Authentication required",
		})
		return
	}

	bookIDStr := c.Param("id")

	err := b.Service.DeleteBook(userID_iface, bookIDStr)
	if err != nil {

		switch {
		case errors.Is(err, errs.ErrInvalidID):
			c.JSON(http.StatusBadRequest, models.ErrorResponse{
				Error: "Invalid ID",
			})
		case errors.Is(err, errs.ErrNotFound):
			c.JSON(http.StatusNotFound, models.ErrorResponse{
				Error: "Could not found record",
			})
		case errors.Is(err, errs.ErrDBOperation):
			c.JSON(http.StatusInternalServerError, models.ErrorResponse{
				Error: "Database operation failed",
			})
		default:
			c.JSON(http.StatusInternalServerError, models.ErrorResponse{
				Error: "Internal server error",
			})
		}

		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse{
		Message: "Book was successfully deleted",
	})
}

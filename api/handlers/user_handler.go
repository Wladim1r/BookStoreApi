package handlers

import (
	"bookstore-api/api/service"
	"bookstore-api/internal/models"
	"bookstore-api/internal/utils"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	Service service.UserService
}

func NewUserHandler(s service.UserService) *UserHandler {
	return &UserHandler{Service: s}
}

func (u *UserHandler) Register(c *gin.Context) {
	var creds models.Request

	if err := c.ShouldBindJSON(&creds); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid body request",
		})
		return
	}

	user := models.User{
		Username: creds.Username,
		Password: creds.Password,
	}

	if err := user.HashPassword(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to hash password",
		})
		return
	}

	if err := u.Service.CreateUser(user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Database connection failed",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "User created",
	})
}

func (u *UserHandler) Login(c *gin.Context) {
	var creds models.Request

	if err := c.ShouldBindJSON(&creds); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid body request",
		})
		return
	}

	user, err := u.Service.GetByUsername(creds.Username)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "You have not registered yet",
		})
		return
	}

	if err := user.CheckPassword(creds.Password); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Entred incorrect password",
		})
		return
	}

	token, err := utils.GenerateToken(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Could not create JWT token",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": fmt.Sprintf("Your token: %s", token),
	})
}

func (u *UserHandler) GetAllUsers(c *gin.Context) {
	users, err := u.Service.GetAllUsers()

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Could not get Users list",
		})
		return
	}

	c.JSON(http.StatusOK, users)
}

func (u *UserHandler) DeleteByUsername(c *gin.Context) {
	username := c.Param("username")

	if _, err := u.Service.GetByUsername(username); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "User is not exists",
		})
		return
	}

	if err := u.Service.DeleteByUsername(username); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Database connection failed",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "User account was successfully deleted",
	})
}

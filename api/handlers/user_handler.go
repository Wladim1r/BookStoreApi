package api

import (
	"bookstore-api/api/service"
	"bookstore-api/internal/models"
	"bookstore-api/internal/utils"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	service service.UserService
}

func NewUserHandler(s service.UserService) *UserHandler {
	return &UserHandler{service: s}
}

func (u *UserHandler) Register(c *gin.Context) {
	var creds models.RegisterRequest

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

	if err := u.service.CreateUser(user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "User created",
	})
}

func (u *UserHandler) Login(c *gin.Context) {
	var creds models.LoginRequest

	if err := c.ShouldBindJSON(&creds); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid body request",
		})
		return
	}

	user, err := u.service.GetByUsername(creds.Username)
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
	users, err := u.service.GetAllUsers()

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

	if _, err := u.service.GetByUsername(username); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "User is not exists",
		})
		return
	}

	if err := u.service.DeleteByUsername(username); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Could not delete User",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "User account was successfully deleted",
	})
}

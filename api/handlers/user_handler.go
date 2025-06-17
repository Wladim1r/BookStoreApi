package handlers

import (
	"bookstore-api/api/service"
	"bookstore-api/internal/lib/errs"
	"bookstore-api/internal/lib/sl"
	"bookstore-api/internal/models"
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	Service service.UserService
}

func NewUserHandler(s service.UserService) *UserHandler {
	return &UserHandler{Service: s}
}

// @Summary Register new user
// @Description Register user with his username and password
// @Tags Authorization
// @ID register-user
// @Accept json
// @Produce json
// @Param request body models.Request true "Credentials for create user account"
// @Success 200 {object} models.SuccessResponse "User created successfully"
// @Failure 400 {object} models.ErrorResponse "Invalid body request"
// @Failure 401 {object} models.ErrorResponse "User unauthorized"
// @Failure 500 {object} models.ErrorResponse "Internal Server Error"
// @Router /auth/register [post]
func (u *UserHandler) Register(c *gin.Context) {
	var creds models.Request

	if err := c.ShouldBindJSON(&creds); err != nil {
		slog.Error("handlers.Register", sl.Error(err))

		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error: "Invalid body request",
		})
		return
	}

	slog.Info("credentials", "creds", creds)

	err := u.Service.CreateUser(creds)
	if err != nil {
		slog.Error("handlers.Register", sl.Error(err))

		switch {
		case errors.Is(err, errs.ErrInternal):
			c.JSON(http.StatusInternalServerError, models.ErrorResponse{
				Error: "Generate password failed",
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

	c.JSON(http.StatusCreated, models.SuccessResponse{
		Message: "User created",
	})
}

// @Summary Login
// @Description Login in created account with yourself credentials
// @Tags Authorization
// @ID login-user
// @Accept json
// @Produce json
// @Param request body models.Request true "Credentials for login in created accound"
// @Success 200 {object} models.SuccessResponse "Give Token after successfully authorization"
// @Failure 400 {object} models.ErrorResponse "Invalid body request"
// @Failure 401 {object} models.ErrorResponse "User do not registred"
// @Failure 403 {object} models.ErrorResponse "Incorrect password"
// @Failure 500 {object} models.ErrorResponse "Internal Server Error"
// @Router /auth/login [post]
func (u *UserHandler) Login(c *gin.Context) {
	var creds models.Request

	if err := c.ShouldBindJSON(&creds); err != nil {
		slog.Error("handlers.Login 94", sl.Error(err))

		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error: "Invalid body request",
		})
		return
	}

	slog.Info("credentials", "creds", creds)

	token, err := u.Service.GetUserToken(creds)
	if err != nil {
		slog.Error("handlers.Login 106", sl.Error(err))

		switch {
		case errors.Is(err, errs.ErrNotRegistred):
			c.JSON(http.StatusForbidden, models.ErrorResponse{
				Error: "Incorrect password",
			})
		case errors.Is(err, errs.ErrNotAuthorized):
			c.JSON(http.StatusUnauthorized, models.ErrorResponse{
				Error: "You are not authorized",
			})
		case errors.Is(err, errs.ErrDBOperation):
			c.JSON(http.StatusInternalServerError, models.ErrorResponse{
				Error: "Database operation failed",
			})
		case errors.Is(err, errs.ErrInternal):
			c.JSON(http.StatusInternalServerError, models.ErrorResponse{
				Error: "Generate JWT token failed",
			})
		default:
			c.JSON(http.StatusInternalServerError, models.ErrorResponse{
				Error: "Internal server error",
			})
		}

		return
	}

	slog.Debug("users jwt token", "token", token)

	c.JSON(http.StatusOK, models.SuccessResponse{
		Message: fmt.Sprintf("Your token: %s", token),
	})
}

// @Summary Get all users
// @Description Get credentials of all users
// @Tags Admin
// @ID get-all-users
// @Security BasicAuth
// @Accept json
// @Produce json
// @Success 200 {object} models.UsersResponse "List of all users"
// @Failure 401 {object} models.ErrorResponse "User unauthorized"
// @Failure 404 {object} models.ErrorResponse "Records not found"
// @Failure 500 {object} models.ErrorResponse "Database or Server error"
// @Router /admin/users [get]
func (u *UserHandler) GetAllUsers(c *gin.Context) {
	users, err := u.Service.GetAllUsers()

	if err != nil {
		slog.Error("handlers.GetAllUsers 157", sl.Error(err))

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

	slog.Debug("users quantity", "number", len(users))

	c.JSON(http.StatusOK, models.UsersResponse{
		Users: users,
	})
}

// @Summary Delete user
// @Description Delete user by username
// @Tags Admin
// @ID delete-user
// @Security BasicAuth
// @Accept json
// @Produce json
// @Param username path string true "Username to delete"
// @Success 200 {object} models.UsersResponse "Message about successfully deleting"
// @Failure 401 {object} models.ErrorResponse "User unauthorized"
// @Failure 404 {object} models.ErrorResponse "Record not found"
// @Failure 500 {object} models.ErrorResponse "Database or Server error"
// @Router /admin/users/{username} [delete]
func (u *UserHandler) DeleteByUsername(c *gin.Context) {
	username := c.Param("username")

	slog.Info("username to delete", "username", username)

	err := u.Service.DeleteByUsername(username)
	if err != nil {
		slog.Error("handlers.DeleteByUsername 204", sl.Error(err))

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
		Message: "User account was successfully deleted",
	})
}

package service

import (
	"bookstore-api/api/repository"
	"bookstore-api/internal/lib/errs"
	"bookstore-api/internal/models"
	"bookstore-api/internal/utils"
	"fmt"
)

type UserService interface {
	CreateUser(models.Request) error
	GetAllUsers() ([]models.UserResponse, error)
	GetUserToken(models.Request) (string, error)
	DeleteByUsername(string) error
}

type userService struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) UserService {
	return &userService{repo: repo}
}

func (s *userService) CreateUser(creds models.Request) error {
	user := models.User{
		Username: creds.Username,
		Password: creds.Password,
	}

	if err := user.HashPassword(); err != nil {
		return fmt.Errorf("%w: %v", errs.ErrInternal, err)
	}

	return s.repo.CreateUser(user)
}

func (s *userService) GetAllUsers() ([]models.UserResponse, error) {
	results, err := s.repo.GetAllUsers()
	if err != nil {
		return nil, err
	}

	users := make([]models.UserResponse, len(results))
	for i, res := range results {
		users[i] = models.UserResponse{
			ID:       res.ID,
			Username: res.Username,
			Total:    len(res.Books),
		}
	}

	return users, nil

}

func (s *userService) GetUserToken(creds models.Request) (string, error) {
	user, err := s.repo.GetByUsername(creds.Username)
	if err != nil {
		return "", err
	}

	err = user.CheckPassword(creds.Password)
	if err != nil {
		return "", fmt.Errorf("%w: %v", errs.ErrNotRegistred, err)
	}

	token, err := utils.GenerateToken(user.ID)
	if err != nil {
		return "", fmt.Errorf("%w: %v", errs.ErrInternal, err)
	}

	return token, nil
}

func (s *userService) DeleteByUsername(username string) error {
	return s.repo.DeleteByUsername(username)
}

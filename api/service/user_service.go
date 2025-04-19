package service

import (
	"bookstore-api/api/repository"
	"bookstore-api/internal/models"
)

type UserService interface {
	CreateUser(models.User) error
	GetAllUsers() ([]models.User, error)
	GetByUsername(string) (models.User, error)
	DeleteByUsername(string) error
}

type userService struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) UserService {
	return &userService{repo: repo}
}

func (s *userService) CreateUser(user models.User) error {
	return s.repo.CreateUser(user)
}

func (s *userService) GetAllUsers() ([]models.User, error) {
	return s.repo.GetAllUsers()
}

func (s *userService) GetByUsername(username string) (models.User, error) {
	return s.repo.GetByUsername(username)
}

func (s *userService) DeleteByUsername(username string) error {
	return s.repo.DeleteByUsername(username)
}

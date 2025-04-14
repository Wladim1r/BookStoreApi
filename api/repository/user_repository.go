package repository

import (
	"bookstore-api/internal/models"

	"gorm.io/gorm"
)

type UserRepository interface {
	CreateUser(models.User) error
	GetAllUsers() ([]models.User, error)
	GetByUsername(string) (models.User, error)
	DeleteByUsername(string) error
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) CreateUser(user models.User) error {
	return r.db.Create(&user).Error
}

func (r *userRepository) GetAllUsers() ([]models.User, error) {
	var users []models.User
	err := r.db.Find(&users).Error

	return users, err
}

func (r *userRepository) GetByUsername(username string) (models.User, error) {
	var user models.User
	err := r.db.Where("username = ?", username).First(&user).Error

	return user, err
}

func (r *userRepository) DeleteByUsername(username string) error {
	return r.db.Unscoped().Where("username = ?", username).Delete(&models.User{}).Error
}

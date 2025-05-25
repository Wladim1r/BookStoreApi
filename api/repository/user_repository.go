package repository

import (
	"bookstore-api/internal/lib/errs"
	"bookstore-api/internal/models"
	"fmt"

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
	result := r.db.Create(&user)

	if result.Error != nil {
		return fmt.Errorf("%w: %v", errs.ErrDBOperation, result.Error)
	}

	return nil
}

func (r *userRepository) GetAllUsers() ([]models.User, error) {
	var users []models.User
	result := r.db.Preload("Books").Find(&users)

	if result.Error != nil {
		return nil, fmt.Errorf("%w: %v", errs.ErrDBOperation, result.Error)
	}

	if result.RowsAffected == 0 {
		return nil, errs.ErrNotFound
	}

	return users, nil
}

func (r *userRepository) GetByUsername(username string) (models.User, error) {
	var user models.User
	result := r.db.Where("username = ?", username).First(&user)

	if result.Error != nil {
		return models.User{}, fmt.Errorf("%w: %v", errs.ErrDBOperation, result.Error)
	}

	if result.RowsAffected == 0 {
		return models.User{}, errs.ErrNotAuthorized
	}

	return user, nil
}

func (r *userRepository) DeleteByUsername(username string) error {
	result := r.db.Unscoped().Where("username = ?", username).Delete(&models.User{})

	if result.Error != nil {
		return fmt.Errorf("%w: %v", errs.ErrDBOperation, result.Error)
	}

	if result.RowsAffected == 0 {
		return errs.ErrNotFound
	}

	return nil
}

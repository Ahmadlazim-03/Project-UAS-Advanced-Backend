package repository

import (
	"github.com/Ahmadlazim-03/Project-UAS-Advanced-Backend/database"
	"github.com/Ahmadlazim-03/Project-UAS-Advanced-Backend/models"
	"gorm.io/gorm"
)

type UserRepository interface {
	FindByUsername(username string) (*models.User, error)
	CreateUser(user *models.User) error
	FindRoleByName(name string) (*models.Role, error)
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository() UserRepository {
	return &userRepository{
		db: database.DB,
	}
}

func (r *userRepository) FindByUsername(username string) (*models.User, error) {
	var user models.User
	err := r.db.Preload("Role.Permissions").Where("username = ?", username).First(&user).Error
	return &user, err
}

func (r *userRepository) CreateUser(user *models.User) error {
	return r.db.Create(user).Error
}

func (r *userRepository) FindRoleByName(name string) (*models.Role, error) {
	var role models.Role
	err := r.db.Where("name = ?", name).First(&role).Error
	return &role, err
}

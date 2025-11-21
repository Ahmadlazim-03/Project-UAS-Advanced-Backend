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
	FindAllUsers() ([]models.User, error)
	FindUserByID(id string) (*models.User, error)
	UpdateUser(user *models.User) error
	DeleteUser(id string) error
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

func (r *userRepository) FindAllUsers() ([]models.User, error) {
	var users []models.User
	err := r.db.Preload("Role").Find(&users).Error
	return users, err
}

func (r *userRepository) FindUserByID(id string) (*models.User, error) {
	var user models.User
	err := r.db.Preload("Role").Where("id = ?", id).First(&user).Error
	return &user, err
}

func (r *userRepository) UpdateUser(user *models.User) error {
	return r.db.Save(user).Error
}

func (r *userRepository) DeleteUser(id string) error {
	return r.db.Delete(&models.User{}, "id = ?", id).Error
}

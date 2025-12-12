package repository

import (
	"student-achievement-system/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserRepository interface {
	FindByUsername(username string) (*models.User, error)
	FindByEmail(email string) (*models.User, error)
	FindByUsernameOrEmail(identifier string) (*models.User, error)
	FindByID(id uuid.UUID) (*models.User, error)
	Create(user *models.User) error
	Update(user *models.User) error
	Delete(id uuid.UUID) error
	FindAll(offset, limit int) ([]models.User, int64, error)
	GetUserPermissions(roleID uuid.UUID) ([]string, error)
	FindDeleted(offset, limit int) ([]models.User, int64, error)
	Restore(id uuid.UUID) error
	HardDelete(id uuid.UUID) error
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) FindByUsername(username string) (*models.User, error) {
	var user models.User
	err := r.db.Where("username = ?", username).Preload("Role").First(&user).Error
	return &user, err
}

func (r *userRepository) FindByEmail(email string) (*models.User, error) {
	var user models.User
	err := r.db.Where("email = ?", email).Preload("Role").First(&user).Error
	return &user, err
}

func (r *userRepository) FindByUsernameOrEmail(identifier string) (*models.User, error) {
	var user models.User
	err := r.db.Where("username = ? OR email = ?", identifier, identifier).
		Preload("Role").First(&user).Error
	return &user, err
}

func (r *userRepository) FindByID(id uuid.UUID) (*models.User, error) {
	var user models.User
	err := r.db.Where("id = ?", id).Preload("Role").First(&user).Error
	return &user, err
}

func (r *userRepository) Create(user *models.User) error {
	return r.db.Create(user).Error
}

func (r *userRepository) Update(user *models.User) error {
	return r.db.Save(user).Error
}

func (r *userRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&models.User{}, id).Error
}

func (r *userRepository) FindAll(offset, limit int) ([]models.User, int64, error) {
	var users []models.User
	var total int64

	r.db.Model(&models.User{}).Count(&total)
	err := r.db.Offset(offset).Limit(limit).Preload("Role").Find(&users).Error

	return users, total, err
}

func (r *userRepository) GetUserPermissions(roleID uuid.UUID) ([]string, error) {
	var rolePermissions []models.RolePermission
	if err := r.db.Where("role_id = ?", roleID).
		Preload("Permission").Find(&rolePermissions).Error; err != nil {
		return nil, err
	}

	permissions := make([]string, len(rolePermissions))
	for i, rp := range rolePermissions {
		permissions[i] = rp.Permission.Name
	}

	return permissions, nil
}

func (r *userRepository) FindDeleted(offset, limit int) ([]models.User, int64, error) {
	var users []models.User
	var total int64

	r.db.Unscoped().Where("deleted_at IS NOT NULL").Model(&models.User{}).Count(&total)
	err := r.db.Unscoped().Where("deleted_at IS NOT NULL").
		Offset(offset).Limit(limit).
		Preload("Role").Find(&users).Error

	return users, total, err
}

func (r *userRepository) Restore(id uuid.UUID) error {
	return r.db.Unscoped().Model(&models.User{}).Where("id = ?", id).Update("deleted_at", nil).Error
}

func (r *userRepository) HardDelete(id uuid.UUID) error {
	return r.db.Unscoped().Delete(&models.User{}, id).Error
}

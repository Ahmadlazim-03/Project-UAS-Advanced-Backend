package repository

import (
	"student-achievement-system/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type RoleRepository interface {
	FindByID(id uuid.UUID) (*models.Role, error)
	FindByName(name string) (*models.Role, error)
	FindAll() ([]models.Role, error)
	Create(role *models.Role) error
	Update(role *models.Role) error
	Delete(id uuid.UUID) error
}

type roleRepository struct {
	db *gorm.DB
}

func NewRoleRepository(db *gorm.DB) RoleRepository {
	return &roleRepository{db: db}
}

func (r *roleRepository) FindByID(id uuid.UUID) (*models.Role, error) {
	var role models.Role
	query := `SELECT * FROM roles WHERE id = ? LIMIT 1`
	err := r.db.Raw(query, id).Scan(&role).Error
	return &role, err
}

func (r *roleRepository) FindByName(name string) (*models.Role, error) {
	var role models.Role
	query := `SELECT * FROM roles WHERE name = ? LIMIT 1`
	err := r.db.Raw(query, name).Scan(&role).Error
	return &role, err
}

func (r *roleRepository) FindAll() ([]models.Role, error) {
	var roles []models.Role
	query := `SELECT * FROM roles ORDER BY name`
	err := r.db.Raw(query).Scan(&roles).Error
	return roles, err
}

func (r *roleRepository) Create(role *models.Role) error {
	query := `
		INSERT INTO roles (id, name, description, created_at)
		VALUES (?, ?, ?, ?)
	`
	return r.db.Exec(query, role.ID, role.Name, role.Description, role.CreatedAt).Error
}

func (r *roleRepository) Update(role *models.Role) error {
	query := `
		UPDATE roles 
		SET name = ?, description = ?
		WHERE id = ?
	`
	return r.db.Exec(query, role.Name, role.Description, role.ID).Error
}

func (r *roleRepository) Delete(id uuid.UUID) error {
	query := `DELETE FROM roles WHERE id = ?`
	return r.db.Exec(query, id).Error
}

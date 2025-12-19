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
	query := `
		SELECT u.*, r.id as "Role__id", r.name as "Role__name", r.description as "Role__description", r.created_at as "Role__created_at"
		FROM users u
		LEFT JOIN roles r ON u.role_id = r.id
		WHERE u.username = ? AND u.deleted_at IS NULL
		LIMIT 1
	`
	err := r.db.Raw(query, username).Scan(&user).Error
	if err != nil {
		return nil, err
	}
	// Load role separately if needed
	if user.RoleID != uuid.Nil {
		r.db.Raw("SELECT * FROM roles WHERE id = ?", user.RoleID).Scan(&user.Role)
	}
	return &user, nil
}

func (r *userRepository) FindByEmail(email string) (*models.User, error) {
	var user models.User
	query := `
		SELECT u.*, r.id as "Role__id", r.name as "Role__name", r.description as "Role__description", r.created_at as "Role__created_at"
		FROM users u
		LEFT JOIN roles r ON u.role_id = r.id
		WHERE u.email = ? AND u.deleted_at IS NULL
		LIMIT 1
	`
	err := r.db.Raw(query, email).Scan(&user).Error
	if err != nil {
		return nil, err
	}
	if user.RoleID != uuid.Nil {
		r.db.Raw("SELECT * FROM roles WHERE id = ?", user.RoleID).Scan(&user.Role)
	}
	return &user, nil
}

func (r *userRepository) FindByUsernameOrEmail(identifier string) (*models.User, error) {
	var user models.User
	query := `
		SELECT * FROM users 
		WHERE (username = ? OR email = ?) AND deleted_at IS NULL 
		LIMIT 1
	`
	err := r.db.Raw(query, identifier, identifier).Scan(&user).Error
	if err != nil {
		return nil, err
	}
	if user.RoleID != uuid.Nil {
		r.db.Raw("SELECT * FROM roles WHERE id = ?", user.RoleID).Scan(&user.Role)
	}
	return &user, nil
}

func (r *userRepository) FindByID(id uuid.UUID) (*models.User, error) {
	var user models.User
	query := `
		SELECT * FROM users 
		WHERE id = ? AND deleted_at IS NULL 
		LIMIT 1
	`
	err := r.db.Raw(query, id).Scan(&user).Error
	if err != nil {
		return nil, err
	}
	if user.RoleID != uuid.Nil {
		r.db.Raw("SELECT * FROM roles WHERE id = ?", user.RoleID).Scan(&user.Role)
	}
	return &user, nil
}

func (r *userRepository) Create(user *models.User) error {
	// Generate UUID if not set
	if user.ID == uuid.Nil {
		user.ID = uuid.New()
	}
	
	query := `
		INSERT INTO users (id, username, email, password_hash, full_name, role_id, is_active, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, NOW(), NOW())
	`
	return r.db.Exec(query, 
		user.ID, user.Username, user.Email, user.PasswordHash, 
		user.FullName, user.RoleID, user.IsActive,
	).Error
}

func (r *userRepository) Update(user *models.User) error {
	query := `
		UPDATE users 
		SET username = ?, email = ?, password_hash = ?, full_name = ?, role_id = ?, updated_at = ?
		WHERE id = ? AND deleted_at IS NULL
	`
	return r.db.Exec(query,
		user.Username, user.Email, user.PasswordHash, user.FullName, 
		user.RoleID, user.UpdatedAt, user.ID,
	).Error
}

func (r *userRepository) Delete(id uuid.UUID) error {
	query := `UPDATE users SET deleted_at = NOW() WHERE id = ?`
	return r.db.Exec(query, id).Error
}

func (r *userRepository) FindAll(offset, limit int) ([]models.User, int64, error) {
	var users []models.User
	var total int64

	// Count total
	countQuery := `SELECT COUNT(*) FROM users WHERE deleted_at IS NULL`
	r.db.Raw(countQuery).Scan(&total)

	// Get users
	query := `
		SELECT * FROM users 
		WHERE deleted_at IS NULL 
		ORDER BY created_at DESC 
		LIMIT ? OFFSET ?
	`
	err := r.db.Raw(query, limit, offset).Scan(&users).Error
	
	// Load roles for each user
	for i := range users {
		if users[i].RoleID != uuid.Nil {
			r.db.Raw("SELECT * FROM roles WHERE id = ?", users[i].RoleID).Scan(&users[i].Role)
		}
	}

	return users, total, err
}

func (r *userRepository) GetUserPermissions(roleID uuid.UUID) ([]string, error) {
	var permissions []string
	query := `
		SELECT p.name 
		FROM role_permissions rp
		INNER JOIN permissions p ON rp.permission_id = p.id
		WHERE rp.role_id = ?
	`
	err := r.db.Raw(query, roleID).Pluck("name", &permissions).Error
	return permissions, err
}

func (r *userRepository) FindDeleted(offset, limit int) ([]models.User, int64, error) {
	var users []models.User
	var total int64

	// Count deleted users
	countQuery := `SELECT COUNT(*) FROM users WHERE deleted_at IS NOT NULL`
	r.db.Raw(countQuery).Scan(&total)

	// Get deleted users
	query := `
		SELECT * FROM users 
		WHERE deleted_at IS NOT NULL 
		ORDER BY deleted_at DESC 
		LIMIT ? OFFSET ?
	`
	err := r.db.Raw(query, limit, offset).Scan(&users).Error

	// Load roles
	for i := range users {
		if users[i].RoleID != uuid.Nil {
			r.db.Raw("SELECT * FROM roles WHERE id = ?", users[i].RoleID).Scan(&users[i].Role)
		}
	}

	return users, total, err
}

func (r *userRepository) Restore(id uuid.UUID) error {
	query := `UPDATE users SET deleted_at = NULL WHERE id = ?`
	return r.db.Exec(query, id).Error
}

func (r *userRepository) HardDelete(id uuid.UUID) error {
	query := `DELETE FROM users WHERE id = ?`
	return r.db.Exec(query, id).Error
}

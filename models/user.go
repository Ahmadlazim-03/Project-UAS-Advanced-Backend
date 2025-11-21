package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID           uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	Username     string    `gorm:"unique;not null;index"`
	Email        string    `gorm:"unique;not null;index"`
	PasswordHash string    `gorm:"not null"`
	FullName     string    `gorm:"not null"`
	RoleID       uuid.UUID `gorm:"type:uuid;not null;index"`
	Role         Role      `gorm:"foreignKey:RoleID"`
	IsActive     bool      `gorm:"default:true;index"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type Role struct {
	ID          uuid.UUID    `gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	Name        string       `gorm:"unique;not null"`
	Description string
	Permissions []Permission `gorm:"many2many:role_permissions;"`
	CreatedAt   time.Time
}

type Permission struct {
	ID          uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	Name        string    `gorm:"unique;not null"`
	Resource    string    `gorm:"not null"`
	Action      string    `gorm:"not null"`
	Description string
}

type RolePermission struct {
	RoleID       uuid.UUID `gorm:"type:uuid;primaryKey"`
	PermissionID uuid.UUID `gorm:"type:uuid;primaryKey"`
}

func MigrateUsers(db *gorm.DB) error {
	return db.AutoMigrate(&User{}, &Role{}, &Permission{}, &RolePermission{})
}

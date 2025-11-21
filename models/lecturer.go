package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Lecturer struct {
	ID         uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	UserID     uuid.UUID `gorm:"type:uuid;not null;unique;index"`
	User       User      `gorm:"foreignKey:UserID"`
	LecturerID string    `gorm:"unique;not null;index"`
	Department string    `gorm:"not null"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

func MigrateLecturers(db *gorm.DB) error {
	if db.Migrator().HasTable(&Lecturer{}) {
		return nil
	}
	return db.AutoMigrate(&Lecturer{})
}

package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Lecturer struct {
	ID          uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	UserID      uuid.UUID `gorm:"type:uuid;not null;unique;index" json:"user_id"`
	User        User      `gorm:"foreignKey:UserID" json:"user"`
	LecturerID  string    `gorm:"unique;not null;index" json:"nip"`
	Department  string    `gorm:"not null" json:"department"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func MigrateLecturers(db *gorm.DB) error {
	if db.Migrator().HasTable(&Lecturer{}) {
		return nil
	}
	return db.AutoMigrate(&Lecturer{})
}

package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Student struct {
	ID           uuid.UUID  `gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	UserID       uuid.UUID  `gorm:"type:uuid;not null;unique;index"`
	User         User       `gorm:"foreignKey:UserID"`
	StudentID    string     `gorm:"unique;not null;index"`
	ProgramStudy string     `gorm:"not null"`
	AcademicYear string     `gorm:"not null"`
	AdvisorID    *uuid.UUID `gorm:"type:uuid;index"`
	Advisor      *Lecturer  `gorm:"foreignKey:AdvisorID"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

func MigrateStudents(db *gorm.DB) error {
	if db.Migrator().HasTable(&Student{}) {
		return nil
	}
	return db.AutoMigrate(&Student{})
}

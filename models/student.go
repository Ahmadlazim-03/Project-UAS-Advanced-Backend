package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Student struct {
	ID           uuid.UUID  `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	UserID       uuid.UUID  `gorm:"type:uuid;not null;unique;index" json:"user_id"`
	User         User       `gorm:"foreignKey:UserID" json:"user"`
	StudentID    string     `gorm:"unique;not null;index" json:"nim"`
	ProgramStudy string     `gorm:"not null" json:"program_study"`
	AcademicYear string     `gorm:"not null" json:"academic_year"`
	AdvisorID    *uuid.UUID `gorm:"type:uuid;index" json:"advisor_id,omitempty"`
	Advisor      *Lecturer  `gorm:"foreignKey:AdvisorID" json:"advisor,omitempty"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
}

func MigrateStudents(db *gorm.DB) error {
	if db.Migrator().HasTable(&Student{}) {
		return nil
	}
	return db.AutoMigrate(&Student{})
}

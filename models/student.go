package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Student represents a student in the system
type Student struct {
	ID            uuid.UUID  `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	UserID        uuid.UUID  `gorm:"type:uuid;not null;unique" json:"user_id"`
	User          User       `gorm:"foreignKey:UserID" json:"user,omitempty"`
	StudentID     string     `gorm:"type:varchar(20);unique;not null" json:"student_id"`
	ProgramStudy  string     `gorm:"type:varchar(100)" json:"program_study"`
	AcademicYear  string     `gorm:"type:varchar(10)" json:"academic_year"`
	AdvisorID     *uuid.UUID `gorm:"type:uuid" json:"advisor_id"`
	Advisor       *Lecturer  `gorm:"foreignKey:AdvisorID" json:"advisor,omitempty"`
	CreatedAt     time.Time  `json:"created_at"`
}

// BeforeCreate hook for Student
func (s *Student) BeforeCreate(tx *gorm.DB) error {
	if s.ID == uuid.Nil {
		s.ID = uuid.New()
	}
	return nil
}

// Lecturer represents a lecturer in the system
type Lecturer struct {
	ID         uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	UserID     uuid.UUID `gorm:"type:uuid;not null;unique" json:"user_id"`
	User       User      `gorm:"foreignKey:UserID" json:"user,omitempty"`
	LecturerID string    `gorm:"type:varchar(20);unique;not null" json:"lecturer_id"`
	Department string    `gorm:"type:varchar(100)" json:"department"`
	CreatedAt  time.Time `json:"created_at"`
}

// BeforeCreate hook for Lecturer
func (l *Lecturer) BeforeCreate(tx *gorm.DB) error {
	if l.ID == uuid.Nil {
		l.ID = uuid.New()
	}
	return nil
}

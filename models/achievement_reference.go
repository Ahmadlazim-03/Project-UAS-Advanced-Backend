package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AchievementReference struct {
	ID                 uuid.UUID      `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	StudentID          uuid.UUID      `gorm:"type:uuid;not null;index" json:"student_id"`
	MongoAchievementID string         `gorm:"not null;index" json:"mongo_achievement_id"`
	Status             string         `gorm:"default:'draft';index" json:"status"` // draft, submitted, verified, rejected
	SubmittedAt        *time.Time     `json:"submitted_at,omitempty"`
	VerifiedAt         *time.Time     `json:"verified_at,omitempty"`
	VerifiedBy         *uuid.UUID     `gorm:"type:uuid;index" json:"verified_by,omitempty"`
	RejectionNote      string         `json:"rejection_note,omitempty"`
	CreatedAt          time.Time      `json:"created_at"`
	UpdatedAt          time.Time      `json:"updated_at"`
	DeletedAt          gorm.DeletedAt `gorm:"index" json:"-"` // Soft delete
	
	// Relations
	Student            Student        `gorm:"foreignKey:StudentID" json:"student,omitempty"`
}

func MigrateAchievements(db *gorm.DB) error {
	if db.Migrator().HasTable(&AchievementReference{}) {
		return nil
	}
	return db.AutoMigrate(&AchievementReference{})
}

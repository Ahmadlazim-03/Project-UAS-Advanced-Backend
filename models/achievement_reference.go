package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AchievementReference struct {
	ID                 uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	StudentID          uuid.UUID `gorm:"type:uuid;not null;index"`
	MongoAchievementID string    `gorm:"not null;index"`
	Status             string    `gorm:"default:'draft';index"` // draft, submitted, verified, rejected
	SubmittedAt        *time.Time
	VerifiedAt         *time.Time
	VerifiedBy         *uuid.UUID `gorm:"type:uuid;index"`
	RejectionNote      string
	CreatedAt          time.Time
	UpdatedAt          time.Time
}

func MigrateAchievements(db *gorm.DB) error {
	return db.AutoMigrate(&AchievementReference{})
}

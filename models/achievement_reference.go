package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// AchievementStatus represents the status of an achievement
type AchievementStatus string

const (
	StatusDraft     AchievementStatus = "draft"
	StatusSubmitted AchievementStatus = "submitted"
	StatusVerified  AchievementStatus = "verified"
	StatusRejected  AchievementStatus = "rejected"
	StatusDeleted   AchievementStatus = "deleted"
)

// AchievementReference represents the reference to achievement data in MongoDB
type AchievementReference struct {
	ID                 uuid.UUID         `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	StudentID          uuid.UUID         `gorm:"type:uuid;not null" json:"student_id"`
	MongoAchievementID string            `gorm:"type:varchar(24);not null" json:"mongo_achievement_id"`
	Status             AchievementStatus `gorm:"type:varchar(20);default:'draft'" json:"status"`
	SubmittedAt        *time.Time        `json:"submitted_at,omitempty"`
	VerifiedAt         *time.Time        `json:"verified_at,omitempty"`
	VerifiedBy         *uuid.UUID        `gorm:"type:uuid" json:"verified_by,omitempty"`
	RejectionNote      string            `gorm:"type:text" json:"rejection_note,omitempty"`
	CreatedAt          time.Time         `json:"created_at"`
	UpdatedAt          time.Time         `json:"updated_at"`
}

// BeforeCreate hook for AchievementReference
func (a *AchievementReference) BeforeCreate(tx *gorm.DB) error {
	if a.ID == uuid.Nil {
		a.ID = uuid.New()
	}
	return nil
}

// TableName specifies the table name for AchievementReference
func (AchievementReference) TableName() string {
	return "achievement_references"
}

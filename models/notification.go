package models

import (
	"time"

	"github.com/google/uuid"
)

type NotificationType string

const (
	NotificationTypeAchievementSubmitted NotificationType = "achievement_submitted"
	NotificationTypeAchievementVerified  NotificationType = "achievement_verified"
	NotificationTypeAchievementRejected  NotificationType = "achievement_rejected"
	NotificationTypeAdvisorAssigned      NotificationType = "advisor_assigned"
)

type Notification struct {
	ID        uuid.UUID        `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	UserID    uuid.UUID        `json:"user_id" gorm:"type:uuid;not null;index"`
	Type      NotificationType `json:"type" gorm:"type:varchar(50);not null"`
	Title     string           `json:"title" gorm:"type:varchar(255);not null"`
	Message   string           `json:"message" gorm:"type:text;not null"`
	Data      string           `json:"data" gorm:"type:jsonb"` // Additional data as JSON
	IsRead    bool             `json:"is_read" gorm:"default:false"`
	ReadAt    *time.Time       `json:"read_at"`
	CreatedAt time.Time        `json:"created_at" gorm:"default:CURRENT_TIMESTAMP"`

	// Relations
	User User `json:"user,omitempty" gorm:"foreignKey:UserID;references:ID"`
}

func (Notification) TableName() string {
	return "notifications"
}

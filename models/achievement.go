package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Achievement struct {
	ID              primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	StudentID       string             `bson:"studentId" json:"studentId"`
	AchievementType string             `bson:"achievementType" json:"achievementType"`
	Title           string             `bson:"title" json:"title"`
	Description     string             `bson:"description" json:"description"`
	Category        string             `bson:"category" json:"category"`
	Organizer       string             `bson:"organizer" json:"organizer"`
	AchievementDate time.Time          `bson:"achievementDate" json:"achievement_date"`
	CertificateNo   string             `bson:"certificateNumber" json:"certificate_number"`
	Details         map[string]interface{} `bson:"details" json:"details"`
	Attachments     []Attachment       `bson:"attachments" json:"attachments"`
	Tags            []string           `bson:"tags" json:"tags"`
	Points          int                `bson:"points" json:"points"`
	IsDeleted       bool               `bson:"isDeleted" json:"is_deleted"`
	DeletedAt       *time.Time         `bson:"deletedAt,omitempty" json:"deleted_at,omitempty"`
	CreatedAt       time.Time          `bson:"createdAt" json:"createdAt"`
	UpdatedAt       time.Time          `bson:"updatedAt" json:"updatedAt"`
}

type Attachment struct {
	FileName   string    `bson:"fileName" json:"fileName"`
	FileURL    string    `bson:"fileUrl" json:"fileUrl"`
	FileType   string    `bson:"fileType" json:"fileType"`
	UploadedAt time.Time `bson:"uploadedAt" json:"uploadedAt"`
}

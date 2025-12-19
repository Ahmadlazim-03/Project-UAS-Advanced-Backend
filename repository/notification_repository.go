package repository

import (
	"student-achievement-system/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type NotificationRepository interface {
	Create(notification *models.Notification) error
	FindByUserID(userID uuid.UUID, offset, limit int) ([]models.Notification, int64, error)
	FindUnreadByUserID(userID uuid.UUID) ([]models.Notification, error)
	MarkAsRead(notificationID uuid.UUID) error
	MarkAllAsRead(userID uuid.UUID) error
	CountUnread(userID uuid.UUID) (int64, error)
}

type notificationRepository struct {
	db *gorm.DB
}

func NewNotificationRepository(db *gorm.DB) NotificationRepository {
	return &notificationRepository{db: db}
}

func (r *notificationRepository) Create(notification *models.Notification) error {
	return r.db.Create(notification).Error
}

func (r *notificationRepository) FindByUserID(userID uuid.UUID, offset, limit int) ([]models.Notification, int64, error) {
	var notifications []models.Notification
	var total int64

	query := r.db.Model(&models.Notification{}).Where("user_id = ?", userID)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := query.Order("created_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&notifications).Error

	return notifications, total, err
}

func (r *notificationRepository) FindUnreadByUserID(userID uuid.UUID) ([]models.Notification, error) {
	var notifications []models.Notification
	err := r.db.Where("user_id = ? AND is_read = ?", userID, false).
		Order("created_at DESC").
		Find(&notifications).Error
	return notifications, err
}

func (r *notificationRepository) MarkAsRead(notificationID uuid.UUID) error {
	now := gorm.Expr("CURRENT_TIMESTAMP")
	return r.db.Model(&models.Notification{}).
		Where("id = ?", notificationID).
		Updates(map[string]interface{}{
			"is_read": true,
			"read_at": now,
		}).Error
}

func (r *notificationRepository) MarkAllAsRead(userID uuid.UUID) error {
	now := gorm.Expr("CURRENT_TIMESTAMP")
	return r.db.Model(&models.Notification{}).
		Where("user_id = ? AND is_read = ?", userID, false).
		Updates(map[string]interface{}{
			"is_read": true,
			"read_at": now,
		}).Error
}

func (r *notificationRepository) CountUnread(userID uuid.UUID) (int64, error) {
	var count int64
	err := r.db.Model(&models.Notification{}).
		Where("user_id = ? AND is_read = ?", userID, false).
		Count(&count).Error
	return count, err
}

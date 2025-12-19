package service

import (
	"encoding/json"
	"fmt"
	"student-achievement-system/middleware"
	"student-achievement-system/models"
	"student-achievement-system/repository"
	"student-achievement-system/utils"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type NotificationService interface {
	GetMyNotifications(c *fiber.Ctx) error
	GetUnreadNotifications(c *fiber.Ctx) error
	MarkAsRead(c *fiber.Ctx) error
	MarkAllAsRead(c *fiber.Ctx) error
	GetUnreadCount(c *fiber.Ctx) error
}

type notificationService struct {
	notificationRepo repository.NotificationRepository
	userRepo         repository.UserRepository
}

func NewNotificationService(
	notificationRepo repository.NotificationRepository,
	userRepo repository.UserRepository,
) NotificationService {
	return &notificationService{
		notificationRepo: notificationRepo,
		userRepo:         userRepo,
	}
}

// Helper function to create notification (used by other services)
func CreateNotification(
	notificationRepo repository.NotificationRepository,
	userID uuid.UUID,
	notifType models.NotificationType,
	title string,
	message string,
	data interface{},
) error {
	dataJSON, _ := json.Marshal(data)

	notification := &models.Notification{
		UserID:    userID,
		Type:      notifType,
		Title:     title,
		Message:   message,
		Data:      string(dataJSON),
		IsRead:    false,
		CreatedAt: time.Now(),
	}

	return notificationRepo.Create(notification)
}

// GetMyNotifications godoc
// @Summary      Get my notifications
// @Description  Get paginated list of current user's notifications
// @Tags         Notifications
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        page   query    int  false  "Page number (default 1)"
// @Param        limit  query    int  false  "Items per page (default 10, max 100)"
// @Success      200 {object} map[string]interface{} "List of notifications with pagination"
// @Failure      401 {object} map[string]interface{} "Unauthorized"
// @Failure      500 {object} map[string]interface{} "Internal server error"
// @Router       /notifications [get]
func (s *notificationService) GetMyNotifications(c *fiber.Ctx) error {
	claims := middleware.GetUserFromContext(c)
	if claims == nil {
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "User not authenticated")
	}

	pagination := utils.GetPaginationParams(c)

	notifications, total, err := s.notificationRepo.FindByUserID(
		claims.UserID,
		pagination.Offset,
		pagination.Limit,
	)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to fetch notifications")
	}

	return utils.PaginatedResponse(c, fiber.Map{
		"notifications": notifications,
	}, total, pagination.Page, pagination.Limit)
}

// GetUnreadNotifications godoc
// @Summary      Get unread notifications
// @Description  Get all unread notifications for current user
// @Tags         Notifications
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      200 {object} map[string]interface{} "List of unread notifications"
// @Failure      401 {object} map[string]interface{} "Unauthorized"
// @Failure      500 {object} map[string]interface{} "Internal server error"
// @Router       /notifications/unread [get]
func (s *notificationService) GetUnreadNotifications(c *fiber.Ctx) error {
	claims := middleware.GetUserFromContext(c)
	if claims == nil {
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "User not authenticated")
	}

	notifications, err := s.notificationRepo.FindUnreadByUserID(claims.UserID)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to fetch notifications")
	}

	return utils.SuccessResponse(c, "Unread notifications retrieved successfully", fiber.Map{
		"notifications": notifications,
		"count":         len(notifications),
	})
}

// MarkAsRead godoc
// @Summary      Mark notification as read
// @Description  Mark a specific notification as read
// @Tags         Notifications
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id  path     string  true  "Notification ID (UUID)"
// @Success      200 {object} map[string]interface{} "Notification marked as read"
// @Failure      400 {object} map[string]interface{} "Invalid notification ID"
// @Failure      401 {object} map[string]interface{} "Unauthorized"
// @Failure      500 {object} map[string]interface{} "Internal server error"
// @Router       /notifications/{id}/read [put]
func (s *notificationService) MarkAsRead(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid notification ID")
	}

	if err := s.notificationRepo.MarkAsRead(id); err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to mark notification as read")
	}

	return utils.SuccessResponse(c, "Notification marked as read", fiber.Map{
		"id": id,
	})
}

// MarkAllAsRead godoc
// @Summary      Mark all notifications as read
// @Description  Mark all user's notifications as read
// @Tags         Notifications
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      200 {object} map[string]interface{} "All notifications marked as read"
// @Failure      401 {object} map[string]interface{} "Unauthorized"
// @Failure      500 {object} map[string]interface{} "Internal server error"
// @Router       /notifications/read-all [put]
func (s *notificationService) MarkAllAsRead(c *fiber.Ctx) error {
	claims := middleware.GetUserFromContext(c)
	if claims == nil {
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "User not authenticated")
	}

	if err := s.notificationRepo.MarkAllAsRead(claims.UserID); err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to mark notifications as read")
	}

	return utils.SuccessResponse(c, "All notifications marked as read", nil)
}

// GetUnreadCount godoc
// @Summary      Get unread notification count
// @Description  Get count of unread notifications for current user
// @Tags         Notifications
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      200 {object} map[string]interface{} "Unread count"
// @Failure      401 {object} map[string]interface{} "Unauthorized"
// @Failure      500 {object} map[string]interface{} "Internal server error"
// @Router       /notifications/unread/count [get]
func (s *notificationService) GetUnreadCount(c *fiber.Ctx) error {
	claims := middleware.GetUserFromContext(c)
	if claims == nil {
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "User not authenticated")
	}

	count, err := s.notificationRepo.CountUnread(claims.UserID)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to get unread count")
	}

	return utils.SuccessResponse(c, fmt.Sprintf("You have %d unread notifications", count), fiber.Map{
		"count": count,
	})
}

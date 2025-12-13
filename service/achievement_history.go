package service

import (
	"context"
	"student-achievement-system/database"
	"student-achievement-system/middleware"
	"student-achievement-system/models"
	"student-achievement-system/utils"
	"time"

	"github.com/gofiber/fiber/v2"
)

// GetAchievementHistory godoc
// @Summary      Get achievement status history
// @Description  Get the status change history of an achievement
// @Tags         Achievements
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id  path     string  true  "Achievement ID (MongoDB ObjectID)"
// @Success      200 {object} map[string]interface{} "Achievement status history"
// @Failure      401 {object} map[string]interface{} "Unauthorized"
// @Failure      404 {object} map[string]interface{} "Achievement not found"
// @Failure      500 {object} map[string]interface{} "Internal server error"
// @Router       /achievements/{id}/history [get]
func (s *achievementService) GetAchievementHistory(c *fiber.Ctx) error {
	id := c.Params("id")

	// Get achievement reference
	achievementRef, err := s.achievementRefRepo.FindByMongoID(id)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusNotFound, "Achievement not found")
	}

	// Get status history from database
	type HistoryWithUser struct {
		models.AchievementStatusHistory
		ChangedByName  string `json:"changed_by_name"`
		ChangedByEmail string `json:"changed_by_email"`
	}
	
	query := `
		SELECT ash.*, u.full_name as changed_by_name, u.email as changed_by_email
		FROM achievement_status_history ash
		LEFT JOIN users u ON ash.changed_by = u.id
		WHERE ash.achievement_ref_id = ?
		ORDER BY ash.created_at ASC
	`
	
	var historyWithUsers []HistoryWithUser
	if err := database.PostgresDB.Raw(query, achievementRef.ID).Scan(&historyWithUsers).Error; err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to get achievement history")
	}

	// Format response
	historyResponse := make([]fiber.Map, 0)
	for _, h := range historyWithUsers {
		historyResponse = append(historyResponse, fiber.Map{
			"id":          h.ID,
			"old_status":  h.OldStatus,
			"new_status":  h.NewStatus,
			"changed_by": fiber.Map{
				"id":    h.ChangedBy,
				"name":  h.ChangedByName,
				"email": h.ChangedByEmail,
			},
			"notes":      h.Notes,
			"created_at": h.CreatedAt,
		})
	}

	return utils.SuccessResponse(c, "Achievement history retrieved successfully", fiber.Map{
		"achievement_id": id,
		"current_status": achievementRef.Status,
		"history":        historyResponse,
	})
}

// UploadAttachment godoc
// @Summary      Upload attachment to achievement
// @Description  Upload and attach a file to an existing achievement
// @Tags         Achievements
// @Accept       multipart/form-data
// @Produce      json
// @Security     BearerAuth
// @Param        id    path      string  true   "Achievement ID (MongoDB ObjectID)"
// @Param        file  formData  file    true   "File to upload"
// @Success      200 {object} map[string]interface{} "Attachment uploaded successfully"
// @Failure      400 {object} map[string]interface{} "Invalid file or achievement ID"
// @Failure      401 {object} map[string]interface{} "Unauthorized"
// @Failure      403 {object} map[string]interface{} "Forbidden"
// @Failure      404 {object} map[string]interface{} "Achievement not found"
// @Failure      500 {object} map[string]interface{} "Internal server error"
// @Router       /achievements/{id}/attachments [post]
func (s *achievementService) UploadAttachment(c *fiber.Ctx) error {
	id := c.Params("id")
	claims := middleware.GetUserFromContext(c)
	if claims == nil {
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "User not authenticated")
	}

	// Get achievement reference to check ownership
	achievementRef, err := s.achievementRefRepo.FindByMongoID(id)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusNotFound, "Achievement not found")
	}

	// Get student info to verify ownership
	student, err := s.studentRepo.FindByID(achievementRef.StudentID)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusNotFound, "Student not found")
	}

	// SECURITY CHECK: Only the owner or admin can upload attachments
	if claims.RoleName != "Admin" && student.UserID != claims.UserID {
		return utils.ErrorResponse(c, fiber.StatusForbidden, "You can only upload attachments to your own achievements")
	}

	// Upload file using middleware
	filename, err := middleware.UploadFile(c, "file")
	if err != nil {
		utils.GlobalLogger.Error("File upload failed", err, map[string]interface{}{
			"user_id":        claims.UserID,
			"achievement_id": id,
		})
		return utils.ErrorResponse(c, fiber.StatusBadRequest, err.Error())
	}

	// Get achievement from MongoDB
	achievement, err := s.achievementRepo.FindByID(context.Background(), id)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusNotFound, "Achievement not found in database")
	}

	// Add attachment to achievement
	newAttachment := models.Attachment{
		FileName:   filename,
		FileURL:    "/uploads/achievements/" + filename,
		FileType:   c.Get("Content-Type"),
		UploadedAt: time.Now(),
	}

	achievement.Attachments = append(achievement.Attachments, newAttachment)
	achievement.UpdatedAt = time.Now()

	// Update in MongoDB
	if err := s.achievementRepo.Update(context.Background(), id, achievement); err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to update achievement")
	}

	utils.GlobalLogger.Info("Attachment uploaded successfully", map[string]interface{}{
		"user_id":        claims.UserID,
		"achievement_id": id,
		"filename":       filename,
	})

	return utils.SuccessResponse(c, "Attachment uploaded successfully", fiber.Map{
		"filename":          filename,
		"file_url":          newAttachment.FileURL,
		"file_type":         newAttachment.FileType,
		"uploaded_at":       newAttachment.UploadedAt,
		"total_attachments": len(achievement.Attachments),
	})
}

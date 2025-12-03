package service

import (
	"context"
	"strconv"
	"student-achievement-system/middleware"
	"student-achievement-system/models"
	"student-achievement-system/repository"
	"student-achievement-system/utils"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AchievementService interface {
	ListAchievements(c *fiber.Ctx) error
	GetAchievement(c *fiber.Ctx) error
	CreateAchievement(c *fiber.Ctx) error
	UpdateAchievement(c *fiber.Ctx) error
	DeleteAchievement(c *fiber.Ctx) error
}

type VerificationService interface {
	SubmitForVerification(c *fiber.Ctx) error
	VerifyAchievement(c *fiber.Ctx) error
	RejectAchievement(c *fiber.Ctx) error
	GetAdviseeAchievements(c *fiber.Ctx) error
}

type CreateAchievementRequest struct {
	ReferenceID  string                 `json:"reference_id,omitempty"`
	Title        string                 `json:"title" validate:"required"`
	Description  string                 `json:"description"`
	AchievedDate string                 `json:"achieved_date" validate:"required"`
	Data         map[string]interface{} `json:"data"`
	Attachments  []string               `json:"attachments"`
}

type UpdateAchievementRequest struct {
	Title        string                 `json:"title,omitempty"`
	Description  string                 `json:"description,omitempty"`
	AchievedDate string                 `json:"achieved_date,omitempty"`
	Data         map[string]interface{} `json:"data,omitempty"`
	Attachments  []string               `json:"attachments,omitempty"`
}

type VerifyRequest struct {
	Comments string `json:"comments,omitempty"`
}

type RejectRequest struct {
	Reason string `json:"reason" validate:"required"`
}

type achievementService struct {
	achievementRepo    repository.AchievementRepository
	achievementRefRepo repository.AchievementReferenceRepository
	studentRepo        repository.StudentRepository
	lecturerRepo       repository.LecturerRepository
}

type verificationService struct {
	achievementRefRepo repository.AchievementReferenceRepository
	studentRepo        repository.StudentRepository
	lecturerRepo       repository.LecturerRepository
}

func NewAchievementService(
	achievementRepo repository.AchievementRepository,
	achievementRefRepo repository.AchievementReferenceRepository,
	studentRepo repository.StudentRepository,
	lecturerRepo repository.LecturerRepository,
) AchievementService {
	return &achievementService{
		achievementRepo:    achievementRepo,
		achievementRefRepo: achievementRefRepo,
		studentRepo:        studentRepo,
		lecturerRepo:       lecturerRepo,
	}
}

func NewVerificationService(
	achievementRefRepo repository.AchievementReferenceRepository,
	studentRepo repository.StudentRepository,
	lecturerRepo repository.LecturerRepository,
) VerificationService {
	return &verificationService{
		achievementRefRepo: achievementRefRepo,
		studentRepo:        studentRepo,
		lecturerRepo:       lecturerRepo,
	}
}

func (s *achievementService) ListAchievements(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	// For now, return a simple response
	// Implement pagination with MongoDB aggregation if needed
	return utils.SuccessResponse(c, "Achievements retrieved successfully", fiber.Map{
		"achievements": []interface{}{},
		"pagination": fiber.Map{
			"page":  page,
			"limit": limit,
			"total": 0,
		},
	})
}

func (s *achievementService) GetAchievement(c *fiber.Ctx) error {
	id := c.Params("id")

	achievement, err := s.achievementRepo.FindByID(context.Background(), id)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusNotFound, "Achievement not found")
	}

	return utils.SuccessResponse(c, "Achievement retrieved successfully", achievement)
}

func (s *achievementService) CreateAchievement(c *fiber.Ctx) error {
	var req CreateAchievementRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body")
	}

	claims := middleware.GetUserFromContext(c)
	if claims == nil {
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "User not authenticated")
	}

	// Find student record for this user
	student, err := s.studentRepo.FindByUserID(claims.UserID)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Only students can create achievements")
	}

	// Parse achievement details from req.Data
	achievementDetails := models.AchievementDetails{}
	if req.Data != nil {
		// Map fields from req.Data to achievementDetails
		if competitionName, ok := req.Data["competition_name"].(string); ok {
			achievementDetails.CompetitionName = competitionName
		}
		if competitionLevel, ok := req.Data["competition_level"].(string); ok {
			achievementDetails.CompetitionLevel = competitionLevel
		}
		if medalType, ok := req.Data["medal_type"].(string); ok {
			achievementDetails.MedalType = medalType
		}
		if rank, ok := req.Data["rank"].(float64); ok {
			rankInt := int(rank)
			achievementDetails.Rank = &rankInt
		}
	}

	achievement := &models.Achievement{
		StudentID:       student.ID.String(),
		AchievementType: models.TypeCompetition, // Default to competition
		Title:           req.Title,
		Description:     req.Description,
		Details:         achievementDetails,
		Attachments:     []models.Attachment{},
		Tags:            []string{},
		Points:          0,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	id, err := s.achievementRepo.Create(context.Background(), achievement)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to create achievement")
	}

	achievement.ID, _ = primitive.ObjectIDFromHex(id)

	// Create achievement reference in PostgreSQL
	achievementRef := &models.AchievementReference{
		StudentID:          student.ID,
		MongoAchievementID: id,
		Status:             models.StatusDraft,
		CreatedAt:          time.Now(),
		UpdatedAt:          time.Now(),
	}

	if err := s.achievementRefRepo.Create(achievementRef); err != nil {
		// If reference creation fails, we should probably delete the MongoDB record
		// But for now, just log the error and continue
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to create achievement reference")
	}

	return utils.SuccessResponse(c, "Achievement created successfully", achievement)
}

func (s *achievementService) UpdateAchievement(c *fiber.Ctx) error {
	id := c.Params("id")

	var req UpdateAchievementRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body")
	}

	// Get existing achievement
	achievement, err := s.achievementRepo.FindByID(context.Background(), id)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusNotFound, "Achievement not found")
	}

	// Update fields
	if req.Title != "" {
		achievement.Title = req.Title
	}
	if req.Description != "" {
		achievement.Description = req.Description
	}
	achievement.UpdatedAt = time.Now()

	if err := s.achievementRepo.Update(context.Background(), id, achievement); err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to update achievement")
	}

	return utils.SuccessResponse(c, "Achievement updated successfully", achievement)
}

func (s *achievementService) DeleteAchievement(c *fiber.Ctx) error {
	id := c.Params("id")

	if err := s.achievementRepo.Delete(context.Background(), id); err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to delete achievement")
	}

	return utils.SuccessResponse(c, "Achievement deleted successfully", nil)
}

// Verification Service Methods
func (s *verificationService) SubmitForVerification(c *fiber.Ctx) error {
	id := c.Params("id")

	// Get achievement reference
	achievementRef, err := s.achievementRefRepo.FindByMongoID(id)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusNotFound, "Achievement not found")
	}

	// Update status to submitted
	achievementRef.Status = models.StatusSubmitted
	now := time.Now()
	achievementRef.SubmittedAt = &now
	achievementRef.UpdatedAt = now

	if err := s.achievementRefRepo.Update(achievementRef); err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to submit achievement")
	}

	return utils.SuccessResponse(c, "Achievement submitted for verification", fiber.Map{
		"id":           id,
		"status":       "pending_verification",
		"submitted_at": now,
	})
}

func (s *verificationService) VerifyAchievement(c *fiber.Ctx) error {
	id := c.Params("id")
	claims := middleware.GetUserFromContext(c)
	if claims == nil {
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "User not authenticated")
	}

	var req VerifyRequest
	c.BodyParser(&req)

	// Get achievement reference to check status and get student info
	achievementRef, err := s.achievementRefRepo.FindByMongoID(id)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusNotFound, "Achievement not found")
	}

	// Get student info to check advisor
	student, err := s.studentRepo.FindByID(achievementRef.StudentID)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusNotFound, "Student not found")
	}

	// SECURITY CHECK: Verify lecturer is the advisor of this student
	lecturer, err := s.lecturerRepo.FindByUserID(claims.UserID)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusForbidden, "Only lecturers can verify achievements")
	}

	if student.AdvisorID == nil || *student.AdvisorID != lecturer.ID {
		return utils.ErrorResponse(c, fiber.StatusForbidden, "You can only verify achievements from your advisees")
	}

	// Update achievement reference status
	achievementRef.Status = models.StatusVerified
	now := time.Now()
	achievementRef.VerifiedAt = &now
	achievementRef.VerifiedBy = &lecturer.ID

	if err := s.achievementRefRepo.Update(achievementRef); err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to verify achievement")
	}

	return utils.SuccessResponse(c, "Achievement verified successfully", fiber.Map{
		"id":          id,
		"status":      "verified",
		"verified_by": lecturer.ID,
		"verified_at": now,
		"comments":    req.Comments,
	})
}

func (s *verificationService) RejectAchievement(c *fiber.Ctx) error {
	id := c.Params("id")
	claims := middleware.GetUserFromContext(c)
	if claims == nil {
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "User not authenticated")
	}

	var req RejectRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body")
	}

	// Get achievement reference to check status and get student info
	achievementRef, err := s.achievementRefRepo.FindByMongoID(id)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusNotFound, "Achievement not found")
	}

	// Get student info to check advisor
	student, err := s.studentRepo.FindByID(achievementRef.StudentID)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusNotFound, "Student not found")
	}

	// SECURITY CHECK: Verify lecturer is the advisor of this student
	lecturer, err := s.lecturerRepo.FindByUserID(claims.UserID)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusForbidden, "Only lecturers can reject achievements")
	}

	if student.AdvisorID == nil || *student.AdvisorID != lecturer.ID {
		return utils.ErrorResponse(c, fiber.StatusForbidden, "You can only reject achievements from your advisees")
	}

	// Update achievement reference status
	achievementRef.Status = models.StatusRejected
	now := time.Now()
	achievementRef.VerifiedAt = &now
	achievementRef.VerifiedBy = &lecturer.ID
	achievementRef.RejectionNote = req.Reason

	if err := s.achievementRefRepo.Update(achievementRef); err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to reject achievement")
	}

	return utils.SuccessResponse(c, "Achievement rejected", fiber.Map{
		"id":          id,
		"status":      "rejected",
		"verified_by": lecturer.ID,
		"verified_at": now,
		"reason":      req.Reason,
	})
}

func (s *verificationService) GetAdviseeAchievements(c *fiber.Ctx) error {
	claims := middleware.GetUserFromContext(c)
	if claims == nil {
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "User not authenticated")
	}

	// Get all students under this advisor
	students, err := s.studentRepo.FindByAdvisorID(claims.UserID)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to fetch advisees")
	}

	return utils.SuccessResponse(c, "Advisee achievements retrieved", fiber.Map{
		"advisees": students,
		"note":     "Achievement details would be fetched from MongoDB",
	})
}

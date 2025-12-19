package service

import (
	"context"
	"fmt"
	"student-achievement-system/middleware"
	"student-achievement-system/models"
	"student-achievement-system/repository"
	"student-achievement-system/utils"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AchievementService interface {
	ListAchievements(c *fiber.Ctx) error
	GetAchievement(c *fiber.Ctx) error
	CreateAchievement(c *fiber.Ctx) error
	UpdateAchievement(c *fiber.Ctx) error
	DeleteAchievement(c *fiber.Ctx) error
	GetAchievementHistory(c *fiber.Ctx) error
	UploadAttachment(c *fiber.Ctx) error
}

type VerificationService interface {
	SubmitForVerification(c *fiber.Ctx) error
	VerifyAchievement(c *fiber.Ctx) error
	RejectAchievement(c *fiber.Ctx) error
	GetAdviseeAchievements(c *fiber.Ctx) error
}

type CreateAchievementRequest struct {
	AchievementType string                 `json:"achievement_type" validate:"required,oneof=academic competition organization publication certification other"`
	Title           string                 `json:"title" validate:"required"`
	Description     string                 `json:"description"`
	AchievedDate    string                 `json:"achieved_date" validate:"required"`
	Data            map[string]interface{} `json:"data" validate:"required"`
	Attachments     []AttachmentRequest    `json:"attachments,omitempty"`
	Tags            []string               `json:"tags,omitempty"`
}

type AttachmentRequest struct {
	FileName string `json:"file_name" validate:"required"`
	FileURL  string `json:"file_url" validate:"required"`
	FileType string `json:"file_type" validate:"required"`
}

type UpdateAchievementRequest struct {
	AchievementType string                 `json:"achievement_type,omitempty"`
	Title           string                 `json:"title,omitempty"`
	Description     string                 `json:"description,omitempty"`
	AchievedDate    string                 `json:"achieved_date,omitempty"`
	Data            map[string]interface{} `json:"data,omitempty"`
	Attachments     []AttachmentRequest    `json:"attachments,omitempty"`
	Tags            []string               `json:"tags,omitempty"`
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
	achievementRepo    repository.AchievementRepository
	achievementRefRepo repository.AchievementReferenceRepository
	studentRepo        repository.StudentRepository
	lecturerRepo       repository.LecturerRepository
	notificationRepo   repository.NotificationRepository
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
	achievementRepo repository.AchievementRepository,
	achievementRefRepo repository.AchievementReferenceRepository,
	studentRepo repository.StudentRepository,
	lecturerRepo repository.LecturerRepository,
	notificationRepo repository.NotificationRepository,
) VerificationService {
	return &verificationService{
		achievementRepo:    achievementRepo,
		achievementRefRepo: achievementRefRepo,
		studentRepo:        studentRepo,
		lecturerRepo:       lecturerRepo,
		notificationRepo:   notificationRepo,
	}
}

// ListAchievements godoc
// @Summary      List all achievements
// @Description  Get paginated list of achievements with optional status filter
// @Tags         Achievements
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        page    query    int     false  "Page number (default 1)"
// @Param        limit   query    int     false  "Items per page (default 10, max 100)"
// @Param        status  query    string  false  "Filter by status (draft/submitted/verified/rejected)"
// @Success      200 {object} map[string]interface{} "List of achievements with pagination"
// @Failure      401 {object} map[string]interface{} "Unauthorized"
// @Failure      500 {object} map[string]interface{} "Internal server error"
// @Router       /achievements [get]
func (s *achievementService) ListAchievements(c *fiber.Ctx) error {
	// Get pagination parameters
	pagination := utils.GetPaginationParams(c)

	// Get status filter if provided
	status := c.Query("status", "")

	// Get all achievement references with pagination
	achievementRefs, total, err := s.achievementRefRepo.FindAll(pagination.Offset, pagination.Limit, status)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to retrieve achievements")
	}

	// Combine PostgreSQL reference data with MongoDB achievement data
	var enrichedAchievements []fiber.Map
	for _, ref := range achievementRefs {
		// Get achievement details from MongoDB
		achievement, err := s.achievementRepo.FindByID(context.Background(), ref.MongoAchievementID)
		if err != nil {
			// If MongoDB record not found, skip this entry
			continue
		}

		// Prepare student info
		studentInfo := fiber.Map{
			"id":   ref.StudentID,
			"name": "Unknown Student",
		}
		if ref.Student != nil && ref.Student.User.FullName != "" {
			studentInfo = fiber.Map{
				"id":         ref.StudentID,
				"student_id": ref.Student.StudentID,
				"program":    ref.Student.ProgramStudy,
				"name":       ref.Student.User.FullName,
				"email":      ref.Student.User.Email,
			}
		}

		// Combine both data sources
		enrichedAchievement := fiber.Map{
			"id":                   ref.ID,
			"student_id":           ref.StudentID,
			"student":              studentInfo,
			"mongo_achievement_id": ref.MongoAchievementID,
			"status":               ref.Status,
			"submitted_at":         ref.SubmittedAt,
			"verified_at":          ref.VerifiedAt,
			"verified_by":          ref.VerifiedBy,
			"rejection_note":       ref.RejectionNote,
			"created_at":           ref.CreatedAt,
			"updated_at":           ref.UpdatedAt,
			// MongoDB fields
			"title":            achievement.Title,
			"description":      achievement.Description,
			"achievement_type": achievement.AchievementType,
			"achieved_date":    achievement.Details.EventDate,
			"details":          achievement.Details,
			"attachments":      achievement.Attachments,
			"tags":             achievement.Tags,
			"points":           achievement.Points,
		}
		enrichedAchievements = append(enrichedAchievements, enrichedAchievement)
	}

	return utils.PaginatedResponse(c, fiber.Map{
		"achievements": enrichedAchievements,
	}, total, pagination.Page, pagination.Limit)
}

// GetAchievement godoc
// @Summary      Get achievement by ID
// @Description  Get detailed information of a specific achievement
// @Tags         Achievements
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id  path     string  true  "Achievement ID (MongoDB ObjectID)"
// @Success      200 {object} map[string]interface{} "Achievement details"
// @Failure      401 {object} map[string]interface{} "Unauthorized"
// @Failure      404 {object} map[string]interface{} "Achievement not found"
// @Router       /achievements/{id} [get]
func (s *achievementService) GetAchievement(c *fiber.Ctx) error {
	id := c.Params("id")

	achievement, err := s.achievementRepo.FindByID(context.Background(), id)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusNotFound, "Achievement not found")
	}

	return utils.SuccessResponse(c, "Achievement retrieved successfully", achievement)
}

// CreateAchievement godoc
// @Summary      Create new achievement
// @Description  Create achievement (6 types: academic, competition, organization, publication, certification, other). See examples below.
// @Tags         Achievements
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        achievement  body     CreateAchievementRequest  true  "Achievement data - See examples in SWAGGER_EXAMPLES.md"
// @Success      201 {object} map[string]interface{} "Achievement created"
// @Failure      400 {object} map[string]interface{} "Invalid input"
// @Failure      401 {object} map[string]interface{} "Unauthorized"
// @Failure      500 {object} map[string]interface{} "Internal server error"
// @Router       /achievements [post]
//
// EXAMPLE 1 - COMPETITION:
//
//	{
//	  "achievement_type": "competition",
//	  "title": "Juara 1 Hackathon Nasional 2025",
//	  "description": "Memenangkan kompetisi hackathon tingkat nasional",
//	  "achieved_date": "2025-11-15",
//	  "data": {
//	    "competition_name": "Hackathon Indonesia 2025",
//	    "competition_level": "national",
//	    "rank": 1,
//	    "medal_type": "gold",
//	    "organizer": "Kementerian Pendidikan",
//	    "location": "Jakarta"
//	  },
//	  "attachments": [{
//	    "file_name": "certificate.pdf",
//	    "file_url": "/uploads/cert_123.pdf",
//	    "file_type": "application/pdf"
//	  }],
//	  "tags": ["hackathon", "programming", "AI"]
//	}
//
// EXAMPLE 2 - PUBLICATION:
//
//	{
//	  "achievement_type": "publication",
//	  "title": "Research on Machine Learning Applications",
//	  "description": "Published paper in international journal",
//	  "achieved_date": "2025-10-20",
//	  "data": {
//	    "publication_type": "journal",
//	    "publication_title": "ML Applications in Education",
//	    "authors": ["Ahmad Lazim", "Dr. Budi Santoso"],
//	    "publisher": "IEEE",
//	    "issn": "1234-5678",
//	    "journal_name": "IEEE Transactions on AI"
//	  },
//	  "tags": ["research", "machine-learning", "publication"]
//	}
//
// EXAMPLE 3 - ORGANIZATION:
//
//	{
//	  "achievement_type": "organization",
//	  "title": "Ketua HMTI 2024-2025",
//	  "description": "Memimpin organisasi mahasiswa teknik informatika",
//	  "achieved_date": "2024-08-01",
//	  "data": {
//	    "organization_name": "HMTI Universitas Airlangga",
//	    "position": "Ketua",
//	    "period_start": "2024-08-01",
//	    "period_end": "2025-07-31"
//	  },
//	  "tags": ["leadership", "organization"]
//	}
//
// EXAMPLE 4 - CERTIFICATION:
//
//	{
//	  "achievement_type": "certification",
//	  "title": "AWS Certified Solutions Architect",
//	  "description": "Professional certification from Amazon Web Services",
//	  "achieved_date": "2025-09-15",
//	  "data": {
//	    "certification_name": "AWS Certified Solutions Architect - Associate",
//	    "issued_by": "Amazon Web Services",
//	    "certification_number": "AWS-12345-ABCD",
//	    "valid_until": "2028-09-15"
//	  },
//	  "tags": ["cloud", "aws", "certification"]
//	}
//
// EXAMPLE 5 - ACADEMIC:
//
//	{
//	  "achievement_type": "academic",
//	  "title": "IPK Semester 4.00",
//	  "description": "Meraih IPK sempurna pada semester 5",
//	  "achieved_date": "2025-06-30",
//	  "data": {
//	    "score": 4.00,
//	    "semester": 5,
//	    "achievement_details": "Semua mata kuliah A"
//	  },
//	  "tags": ["academic", "gpa"]
//	}
//
// Valid values:
// - achievement_type: "academic", "competition", "organization", "publication", "certification", "other"
// - competition_level: "international", "national", "regional", "local"
// - medal_type: "gold", "silver", "bronze"
// - publication_type: "journal", "conference", "book"
func (s *achievementService) CreateAchievement(c *fiber.Ctx) error {
	var req CreateAchievementRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body")
	}

	// Validate input
	if err := utils.ValidateStruct(&req); err != nil {
		return utils.ValidationErrorResponse(c, err)
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

	// Parse achievement details from req.Data based on type
	achievementDetails := parseAchievementDetails(req.Data, req.AchievementType)

	// Parse attachments
	attachments := make([]models.Attachment, 0)
	for _, att := range req.Attachments {
		attachments = append(attachments, models.Attachment{
			FileName:   att.FileName,
			FileURL:    att.FileURL,
			FileType:   att.FileType,
			UploadedAt: time.Now(),
		})
	}

	// Calculate points based on achievement type and level
	points := calculatePoints(req.AchievementType, req.Data)

	achievement := &models.Achievement{
		StudentID:       student.ID.String(),
		AchievementType: models.AchievementType(req.AchievementType),
		Title:           req.Title,
		Description:     req.Description,
		Details:         achievementDetails,
		Attachments:     attachments,
		Tags:            req.Tags,
		Points:          points,
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
		// ROLLBACK: Delete MongoDB record if PostgreSQL insert fails
		_ = s.achievementRepo.Delete(context.Background(), id)
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to create achievement reference")
	}

	return utils.SuccessResponse(c, "Achievement created successfully", achievement)
}

// UpdateAchievement godoc
// @Summary      Update achievement
// @Description  Update achievement information by ID
// @Tags         Achievements
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id           path     string                    true  "Achievement ID (MongoDB ObjectID)"
// @Param        achievement  body     UpdateAchievementRequest  true  "Achievement update data"
// @Success      200 {object} map[string]interface{} "Achievement updated successfully"
// @Failure      400 {object} map[string]interface{} "Invalid achievement ID or input"
// @Failure      401 {object} map[string]interface{} "Unauthorized"
// @Failure      403 {object} map[string]interface{} "Forbidden"
// @Failure      404 {object} map[string]interface{} "Achievement not found"
// @Failure      500 {object} map[string]interface{} "Internal server error"
// @Router       /achievements/{id} [put]
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
	if req.AchievementType != "" {
		achievement.AchievementType = models.AchievementType(req.AchievementType)
	}

	// Update details from req.Data using helper function
	if req.Data != nil {
		achievementType := string(achievement.AchievementType)
		if req.AchievementType != "" {
			achievementType = req.AchievementType
		}
		achievement.Details = parseAchievementDetails(req.Data, achievementType)
		// Recalculate points
		achievement.Points = calculatePoints(achievementType, req.Data)
	}

	// Update attachments
	if req.Attachments != nil {
		attachments := make([]models.Attachment, 0)
		for _, att := range req.Attachments {
			attachments = append(attachments, models.Attachment{
				FileName:   att.FileName,
				FileURL:    att.FileURL,
				FileType:   att.FileType,
				UploadedAt: time.Now(),
			})
		}
		achievement.Attachments = attachments
	}

	// Update tags
	if req.Tags != nil {
		achievement.Tags = req.Tags
	}

	// Update achieved date if provided (stored in Details.EventDate)
	if req.AchievedDate != "" {
		eventDate, err := time.Parse("2006-01-02", req.AchievedDate)
		if err == nil {
			achievement.Details.EventDate = &eventDate
		}
	}

	achievement.UpdatedAt = time.Now()

	if err := s.achievementRepo.Update(context.Background(), id, achievement); err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to update achievement")
	}

	return utils.SuccessResponse(c, "Achievement updated successfully", achievement)
}

// DeleteAchievement godoc
// @Summary      Delete achievement
// @Description  Soft delete achievement by ID (changes status to 'deleted')
// @Tags         Achievements
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id  path     string  true  "Achievement ID (MongoDB ObjectID)"
// @Success      200 {object} map[string]interface{} "Achievement deleted successfully"
// @Failure      401 {object} map[string]interface{} "Unauthorized"
// @Failure      403 {object} map[string]interface{} "Forbidden"
// @Failure      404 {object} map[string]interface{} "Achievement not found"
// @Failure      500 {object} map[string]interface{} "Internal server error"
// @Router       /achievements/{id} [delete]
func (s *achievementService) DeleteAchievement(c *fiber.Ctx) error {
	id := c.Params("id")

	// Get achievement reference
	achievementRef, err := s.achievementRefRepo.FindByMongoID(id)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusNotFound, "Achievement not found")
	}

	// Soft delete: Update status to deleted
	achievementRef.Status = models.StatusDeleted
	achievementRef.UpdatedAt = time.Now()

	if err := s.achievementRefRepo.Update(achievementRef); err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to delete achievement")
	}

	return utils.SuccessResponse(c, "Achievement deleted successfully", fiber.Map{
		"id":     id,
		"status": "deleted",
	})
}

// Verification Service Methods
// SubmitForVerification godoc
// @Summary      Submit achievement for verification
// @Description  Submit an achievement to advisor for verification
// @Tags         Verification
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id  path     string  true  "Achievement ID (MongoDB ObjectID)"
// @Success      200 {object} map[string]interface{} "Achievement submitted for verification"
// @Failure      400 {object} map[string]interface{} "Achievement already verified/rejected"
// @Failure      401 {object} map[string]interface{} "Unauthorized"
// @Failure      404 {object} map[string]interface{} "Achievement not found"
// @Failure      500 {object} map[string]interface{} "Internal server error"
// @Router       /achievements/{id}/submit [post]
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

	// Get achievement details for notification
	achievement, _ := s.achievementRepo.FindByID(context.Background(), id)

	// Get student info
	student, _ := s.studentRepo.FindByID(achievementRef.StudentID)

	// Create notification for advisor if assigned
	if student.AdvisorID != nil {
		advisorLecturer, err := s.lecturerRepo.FindByID(*student.AdvisorID)
		if err == nil {
			CreateNotification(
				s.notificationRepo,
				advisorLecturer.UserID,
				models.NotificationTypeAchievementSubmitted,
				"New Achievement Submitted",
				fmt.Sprintf("%s submitted a new achievement: %s", student.User.FullName, achievement.Title),
				fiber.Map{
					"achievement_id": id,
					"student_id":     student.ID,
					"student_name":   student.User.FullName,
				},
			)
		}
	}

	return utils.SuccessResponse(c, "Achievement submitted for verification", fiber.Map{
		"id":           id,
		"status":       "submitted",
		"submitted_at": now,
	})
}

// VerifyAchievement godoc
// @Summary      Verify achievement
// @Description  Advisor verifies an achievement (only for own advisees)
// @Tags         Verification
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id       path     string         true  "Achievement ID (MongoDB ObjectID)"
// @Param        verify   body     VerifyRequest  false "Verification comments"
// @Success      200 {object} map[string]interface{} "Achievement verified successfully"
// @Failure      400 {object} map[string]interface{} "Achievement not submitted or already verified"
// @Failure      401 {object} map[string]interface{} "Unauthorized"
// @Failure      403 {object} map[string]interface{} "Forbidden - not student's advisor"
// @Failure      404 {object} map[string]interface{} "Achievement not found"
// @Failure      500 {object} map[string]interface{} "Internal server error"
// @Router       /achievements/{id}/verify [post]
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

	// SECURITY CHECK: Admin can verify any achievement, Lecturer can only verify their advisees
	var verifierID uuid.UUID

	// Check if user is Admin (has admin role in claims)
	if claims.RoleName == "Admin" {
		// Admin can verify any achievement
		verifierID = claims.UserID
	} else {
		// For lecturers, verify they are the advisor of this student
		lecturer, err := s.lecturerRepo.FindByUserID(claims.UserID)
		if err != nil {
			return utils.ErrorResponse(c, fiber.StatusForbidden, "Only admins or lecturers can verify achievements")
		}

		if student.AdvisorID == nil || *student.AdvisorID != lecturer.ID {
			return utils.ErrorResponse(c, fiber.StatusForbidden, "You can only verify achievements from your advisees")
		}
		verifierID = lecturer.ID
	}

	// Update achievement reference status
	achievementRef.Status = models.StatusVerified
	now := time.Now()
	achievementRef.VerifiedAt = &now
	achievementRef.VerifiedBy = &verifierID

	if err := s.achievementRefRepo.Update(achievementRef); err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to verify achievement")
	}

	// Get achievement details for notification
	achievement, _ := s.achievementRepo.FindByID(context.Background(), id)

	// Create notification for student
	CreateNotification(
		s.notificationRepo,
		student.UserID,
		models.NotificationTypeAchievementVerified,
		"Achievement Verified",
		fmt.Sprintf("Congratulations! Your achievement '%s' has been verified", achievement.Title),
		fiber.Map{
			"achievement_id": id,
			"verified_by":    verifierID,
			"comments":       req.Comments,
		},
	)

	return utils.SuccessResponse(c, "Achievement verified successfully", fiber.Map{
		"id":          id,
		"status":      "verified",
		"verified_by": verifierID,
		"verified_at": now,
		"comments":    req.Comments,
	})
}

// RejectAchievement godoc
// @Summary      Reject achievement
// @Description  Advisor rejects an achievement with reason (only for own advisees)
// @Tags         Verification
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id      path     string         true  "Achievement ID (MongoDB ObjectID)"
// @Param        reject  body     RejectRequest  true  "Rejection reason"
// @Success      200 {object} map[string]interface{} "Achievement rejected successfully"
// @Failure      400 {object} map[string]interface{} "Achievement not submitted or already verified"
// @Failure      401 {object} map[string]interface{} "Unauthorized"
// @Failure      403 {object} map[string]interface{} "Forbidden - not student's advisor"
// @Failure      404 {object} map[string]interface{} "Achievement not found"
// @Failure      500 {object} map[string]interface{} "Internal server error"
// @Router       /achievements/{id}/reject [post]
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

	// SECURITY CHECK: Admin can reject any achievement, Lecturer can only reject their advisees
	var verifierID uuid.UUID

	if claims.RoleName == "Admin" {
		// Admin can reject any achievement
		verifierID = claims.UserID
	} else {
		// For lecturers, verify they are the advisor of this student
		lecturer, err := s.lecturerRepo.FindByUserID(claims.UserID)
		if err != nil {
			return utils.ErrorResponse(c, fiber.StatusForbidden, "Only admins or lecturers can reject achievements")
		}

		if student.AdvisorID == nil || *student.AdvisorID != lecturer.ID {
			return utils.ErrorResponse(c, fiber.StatusForbidden, "You can only reject achievements from your advisees")
		}
		verifierID = lecturer.ID
	}

	// Update achievement reference status
	achievementRef.Status = models.StatusRejected
	now := time.Now()
	achievementRef.VerifiedAt = &now
	achievementRef.VerifiedBy = &verifierID
	achievementRef.RejectionNote = req.Reason

	if err := s.achievementRefRepo.Update(achievementRef); err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to reject achievement")
	}

	// Get achievement details for notification
	achievement, _ := s.achievementRepo.FindByID(context.Background(), id)

	// Create notification for student
	CreateNotification(
		s.notificationRepo,
		student.UserID,
		models.NotificationTypeAchievementRejected,
		"Achievement Rejected",
		fmt.Sprintf("Your achievement '%s' has been rejected. Reason: %s", achievement.Title, req.Reason),
		fiber.Map{
			"achievement_id": id,
			"rejection_note": req.Reason,
			"verified_by":    verifierID,
		},
	)

	return utils.SuccessResponse(c, "Achievement rejected", fiber.Map{
		"id":          id,
		"status":      "rejected",
		"verified_by": verifierID,
		"verified_at": now,
		"reason":      req.Reason,
	})
}

// Helper functions for achievement parsing and point calculation

func parseAchievementDetails(data map[string]interface{}, achievementType string) models.AchievementDetails {
	details := models.AchievementDetails{}

	if data == nil {
		return details
	}

	// Common fields
	if eventDate, ok := data["event_date"].(string); ok {
		if t, err := time.Parse("2006-01-02", eventDate); err == nil {
			details.EventDate = &t
		}
	}
	if location, ok := data["location"].(string); ok {
		details.Location = location
	}
	if organizer, ok := data["organizer"].(string); ok {
		details.Organizer = organizer
	}
	if score, ok := data["score"].(float64); ok {
		details.Score = &score
	}

	switch achievementType {
	case "competition":
		if competitionName, ok := data["competition_name"].(string); ok {
			details.CompetitionName = competitionName
		}
		if competitionLevel, ok := data["competition_level"].(string); ok {
			details.CompetitionLevel = competitionLevel
		}
		if medalType, ok := data["medal_type"].(string); ok {
			details.MedalType = medalType
		}
		if rank, ok := data["rank"].(float64); ok {
			rankInt := int(rank)
			details.Rank = &rankInt
		}

	case "publication":
		if publicationType, ok := data["publication_type"].(string); ok {
			details.PublicationType = publicationType
		}
		if publicationTitle, ok := data["publication_title"].(string); ok {
			details.PublicationTitle = publicationTitle
		}
		if publisher, ok := data["publisher"].(string); ok {
			details.Publisher = publisher
		}
		if issn, ok := data["issn"].(string); ok {
			details.ISSN = issn
		}
		if authors, ok := data["authors"].([]interface{}); ok {
			authorsList := make([]string, 0)
			for _, author := range authors {
				if authorStr, ok := author.(string); ok {
					authorsList = append(authorsList, authorStr)
				}
			}
			details.Authors = authorsList
		}

	case "organization":
		if orgName, ok := data["organization_name"].(string); ok {
			details.OrganizationName = orgName
		}
		if position, ok := data["position"].(string); ok {
			details.Position = position
		}
		if periodStart, ok := data["period_start"].(string); ok {
			if periodEnd, ok := data["period_end"].(string); ok {
				startTime, _ := time.Parse("2006-01-02", periodStart)
				endTime, _ := time.Parse("2006-01-02", periodEnd)
				details.Period = &models.Period{
					Start: startTime,
					End:   endTime,
				}
			}
		}

	case "certification":
		if certName, ok := data["certification_name"].(string); ok {
			details.CertificationName = certName
		}
		if issuedBy, ok := data["issued_by"].(string); ok {
			details.IssuedBy = issuedBy
		}
		if certNumber, ok := data["certification_number"].(string); ok {
			details.CertificationNumber = certNumber
		}
		if validUntil, ok := data["valid_until"].(string); ok {
			if t, err := time.Parse("2006-01-02", validUntil); err == nil {
				details.ValidUntil = &t
			}
		}

	case "academic":
		// Academic achievements use score and custom fields
		if customFields, ok := data["custom_fields"].(map[string]interface{}); ok {
			details.CustomFields = customFields
		} else {
			// Store all data as custom fields for academic
			details.CustomFields = data
		}

	case "other":
		// Store all data as custom fields for other types
		details.CustomFields = data
	}

	return details
}

func calculatePoints(achievementType string, data map[string]interface{}) int {
	basePoints := map[string]int{
		"competition":   100,
		"publication":   150,
		"organization":  50,
		"certification": 75,
		"academic":      25,
		"other":         10,
	}

	points := basePoints[achievementType]

	// Bonus points based on level for competitions
	if achievementType == "competition" {
		if level, ok := data["competition_level"].(string); ok {
			levelBonus := map[string]int{
				"international": 200,
				"national":      100,
				"regional":      50,
				"local":         25,
			}
			points += levelBonus[level]
		}

		// Bonus points for rank
		if rank, ok := data["rank"].(float64); ok {
			rankInt := int(rank)
			if rankInt == 1 {
				points += 100
			} else if rankInt == 2 {
				points += 75
			} else if rankInt == 3 {
				points += 50
			}
		}
	}

	// Bonus points for publications
	if achievementType == "publication" {
		if pubType, ok := data["publication_type"].(string); ok {
			pubBonus := map[string]int{
				"journal":    100,
				"conference": 75,
				"book":       150,
			}
			points += pubBonus[pubType]
		}
	}

	return points
}

// GetAdviseeAchievements godoc
// @Summary      Get advisee achievements
// @Description  Get all achievements from students under advisor's supervision
// @Tags         Verification
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      200 {object} map[string]interface{} "Advisee achievements retrieved"
// @Failure      401 {object} map[string]interface{} "Unauthorized"
// @Failure      403 {object} map[string]interface{} "Forbidden - not a lecturer"
// @Failure      500 {object} map[string]interface{} "Internal server error"
// @Router       /advisee-achievements [get]
func (s *verificationService) GetAdviseeAchievements(c *fiber.Ctx) error {
	claims := middleware.GetUserFromContext(c)
	if claims == nil {
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "User not authenticated")
	}

	// Get pagination parameters
	pagination := utils.GetPaginationParams(c)
	status := c.Query("status", "")

	// Get lecturer record first to get lecturer.ID
	lecturer, err := s.lecturerRepo.FindByUserID(claims.UserID)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusForbidden, "Only lecturers can access advisee achievements")
	}

	// Get all students under this advisor using lecturer.ID (not user_id)
	students, err := s.studentRepo.FindByAdvisorID(lecturer.ID)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to fetch advisees")
	}

	if len(students) == 0 {
		return utils.PaginatedResponse(c, fiber.Map{
			"achievements": []fiber.Map{},
		}, 0, pagination.Page, pagination.Limit)
	}

	// Get student IDs
	studentIDs := make([]uuid.UUID, len(students))
	for i, student := range students {
		studentIDs[i] = student.ID
	}

	// Get all achievement references for these students
	achievementRefs, total, err := s.achievementRefRepo.FindByStudentIDs(studentIDs, pagination.Offset, pagination.Limit, status)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to retrieve achievements")
	}

	// Combine PostgreSQL reference data with MongoDB achievement data
	var enrichedAchievements []fiber.Map
	for _, ref := range achievementRefs {
		// Get achievement details from MongoDB
		achievement, err := s.achievementRepo.FindByID(context.Background(), ref.MongoAchievementID)
		if err != nil {
			// If MongoDB record not found, skip this entry
			continue
		}

		// Find student info
		var studentInfo *models.Student
		for _, student := range students {
			if student.ID == ref.StudentID {
				studentInfo = &student
				break
			}
		}

		// Combine both data sources
		enrichedAchievement := fiber.Map{
			"id":                   ref.ID,
			"student_id":           ref.StudentID,
			"mongo_achievement_id": ref.MongoAchievementID,
			"status":               ref.Status,
			"submitted_at":         ref.SubmittedAt,
			"verified_at":          ref.VerifiedAt,
			"verified_by":          ref.VerifiedBy,
			"rejection_note":       ref.RejectionNote,
			"created_at":           ref.CreatedAt,
			"updated_at":           ref.UpdatedAt,
			// MongoDB fields
			"title":            achievement.Title,
			"description":      achievement.Description,
			"achievement_type": achievement.AchievementType,
			"achieved_date":    achievement.Details.EventDate,
			"details":          achievement.Details,
			"attachments":      achievement.Attachments,
			"tags":             achievement.Tags,
			"points":           achievement.Points,
		}

		// Add student info if found
		if studentInfo != nil {
			enrichedAchievement["student"] = fiber.Map{
				"id":         studentInfo.ID,
				"student_id": studentInfo.StudentID,
				"name":       studentInfo.User.FullName,
				"email":      studentInfo.User.Email,
				"program":    studentInfo.ProgramStudy,
			}
		}

		enrichedAchievements = append(enrichedAchievements, enrichedAchievement)
	}

	return utils.PaginatedResponse(c, fiber.Map{
		"achievements": enrichedAchievements,
	}, total, pagination.Page, pagination.Limit)
}

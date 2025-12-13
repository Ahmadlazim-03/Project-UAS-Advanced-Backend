package service

import (
	"context"
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
	achievementRepo    repository.AchievementRepository
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
	achievementRepo repository.AchievementRepository,
	achievementRefRepo repository.AchievementReferenceRepository,
	studentRepo repository.StudentRepository,
	lecturerRepo repository.LecturerRepository,
) VerificationService {
	return &verificationService{
		achievementRepo:    achievementRepo,
		achievementRefRepo: achievementRefRepo,
		studentRepo:        studentRepo,
		lecturerRepo:       lecturerRepo,
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
// @Param        status  query    string  false  "Filter by status (pending/approved/rejected)"
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
			"title":           achievement.Title,
			"description":     achievement.Description,
			"achievement_type": achievement.AchievementType,
			"achieved_date":   achievement.Details.EventDate,
			"details":         achievement.Details,
			"attachments":     achievement.Attachments,
			"tags":            achievement.Tags,
			"points":          achievement.Points,
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
// @Description  Create a new achievement record (dual-database: PostgreSQL + MongoDB)
// @Tags         Achievements
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        achievement  body     CreateAchievementRequest  true  "Achievement creation data"
// @Success      201 {object} map[string]interface{} "Achievement created successfully"
// @Failure      400 {object} map[string]interface{} "Invalid input or validation error"
// @Failure      401 {object} map[string]interface{} "Unauthorized"
// @Failure      403 {object} map[string]interface{} "Forbidden"
// @Failure      500 {object} map[string]interface{} "Internal server error"
// @Router       /achievements [post]
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

	// Update details from req.Data
	if req.Data != nil {
		// Competition fields
		if competitionName, ok := req.Data["competition_name"].(string); ok {
			achievement.Details.CompetitionName = competitionName
		}
		if competitionLevel, ok := req.Data["competition_level"].(string); ok {
			achievement.Details.CompetitionLevel = competitionLevel
		}
		if medalType, ok := req.Data["medal_type"].(string); ok {
			achievement.Details.MedalType = medalType
		}
		if rank, ok := req.Data["rank"].(float64); ok {
			rankInt := int(rank)
			achievement.Details.Rank = &rankInt
		}

		// Publication fields
		if publicationType, ok := req.Data["publication_type"].(string); ok {
			achievement.Details.PublicationType = publicationType
		}
		if publicationTitle, ok := req.Data["publication_title"].(string); ok {
			achievement.Details.PublicationTitle = publicationTitle
		}
		if publisher, ok := req.Data["publisher"].(string); ok {
			achievement.Details.Publisher = publisher
		}
		if issn, ok := req.Data["issn"].(string); ok {
			achievement.Details.ISSN = issn
		}

		// Common fields
		if organizer, ok := req.Data["organizer"].(string); ok {
			achievement.Details.Organizer = organizer
		}
		if location, ok := req.Data["location"].(string); ok {
			achievement.Details.Location = location
		}
	}

	// Update event date if provided
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
// @Description  Delete achievement by ID (soft delete)
// @Tags         Achievements
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id  path     string  true  "Achievement ID (MongoDB ObjectID)"
// @Success      200 {object} map[string]interface{} "Achievement deleted successfully"
// @Failure      401 {object} map[string]interface{} "Unauthorized"
// @Failure      403 {object} map[string]interface{} "Forbidden"
// @Failure      500 {object} map[string]interface{} "Internal server error"
// @Router       /achievements/{id} [delete]
func (s *achievementService) DeleteAchievement(c *fiber.Ctx) error {
	id := c.Params("id")

	if err := s.achievementRepo.Delete(context.Background(), id); err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to delete achievement")
	}

	return utils.SuccessResponse(c, "Achievement deleted successfully", nil)
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

	return utils.SuccessResponse(c, "Achievement submitted for verification", fiber.Map{
		"id":           id,
		"status":       "pending_verification",
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

	return utils.SuccessResponse(c, "Achievement rejected", fiber.Map{
		"id":          id,
		"status":      "rejected",
		"verified_by": verifierID,
		"verified_at": now,
		"reason":      req.Reason,
	})
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
			"title":           achievement.Title,
			"description":     achievement.Description,
			"achievement_type": achievement.AchievementType,
			"achieved_date":   achievement.Details.EventDate,
			"details":         achievement.Details,
			"attachments":     achievement.Attachments,
			"tags":            achievement.Tags,
			"points":          achievement.Points,
		}

		// Add student info if found
		if studentInfo != nil {
			enrichedAchievement["student"] = fiber.Map{
				"id":          studentInfo.ID,
				"student_id":  studentInfo.StudentID,
				"name":        studentInfo.User.FullName,
				"email":       studentInfo.User.Email,
				"program":     studentInfo.ProgramStudy,
			}
		}

		enrichedAchievements = append(enrichedAchievements, enrichedAchievement)
	}

	return utils.PaginatedResponse(c, fiber.Map{
		"achievements": enrichedAchievements,
	}, total, pagination.Page, pagination.Limit)
}

package service

import (
	"context"
	"student-achievement-system/models"
	"student-achievement-system/repository"
	"student-achievement-system/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type ReportService interface {
	GetStatistics(c *fiber.Ctx) error
	GetStudentReport(c *fiber.Ctx) error
}

type reportService struct {
	achievementRepo    repository.AchievementRepository
	achievementRefRepo repository.AchievementReferenceRepository
	studentRepo        repository.StudentRepository
	lecturerRepo       repository.LecturerRepository
}

func NewReportService(
	achievementRepo repository.AchievementRepository,
	achievementRefRepo repository.AchievementReferenceRepository,
	studentRepo repository.StudentRepository,
	lecturerRepo repository.LecturerRepository,
) ReportService {
	return &reportService{
		achievementRepo:    achievementRepo,
		achievementRefRepo: achievementRefRepo,
		studentRepo:        studentRepo,
		lecturerRepo:       lecturerRepo,
	}
}

// GetStatistics godoc
// @Summary      Get system statistics
// @Description  Get overall statistics (achievements by type, total students, total lecturers)
// @Tags         Reports
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      200 {object} map[string]interface{} "Statistics retrieved successfully"
// @Failure      401 {object} map[string]interface{} "Unauthorized"
// @Failure      500 {object} map[string]interface{} "Internal server error"
// @Router       /reports/statistics [get]
func (s *reportService) GetStatistics(c *fiber.Ctx) error {
	// Get statistics by type and status
	typeCounts, _ := s.achievementRepo.CountByType(context.Background())
	statusCounts, _ := s.achievementRefRepo.CountByStatus()

	// Count students and lecturers
	_, totalStudents, _ := s.studentRepo.FindAll(0, 0)
	_, totalLecturers, _ := s.lecturerRepo.FindAll(0, 0)

	return utils.SuccessResponse(c, "Statistics retrieved successfully", fiber.Map{
		"achievements":        statusCounts,
		"achievement_types":   typeCounts,
		"students":            totalStudents,
		"lecturers":           totalLecturers,
	})
}

// GetStudentReport godoc
// @Summary      Get student report
// @Description  Get achievement report for a specific student
// @Tags         Reports
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id  path     string  true  "Student ID (UUID)"
// @Success      200 {object} map[string]interface{} "Student report retrieved successfully"
// @Failure      400 {object} map[string]interface{} "Invalid student ID"
// @Failure      401 {object} map[string]interface{} "Unauthorized"
// @Failure      404 {object} map[string]interface{} "Student not found"
// @Failure      500 {object} map[string]interface{} "Internal server error"
// @Router       /reports/students/{id} [get]
func (s *reportService) GetStudentReport(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid student ID")
	}

	student, err := s.studentRepo.FindByUserID(id)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusNotFound, "Student not found")
	}

	// Get student's achievements count by status from PostgreSQL
	statusCounts, _ := s.achievementRefRepo.CountByStudentID(student.ID)
	
	// Get achievements count by type from MongoDB
	typeCounts, _ := s.achievementRepo.CountByStudentIDAndType(context.Background(), student.ID.String())

	// Calculate totals
	totalAchievements := int64(0)
	for _, count := range statusCounts {
		totalAchievements += count
	}

	return utils.SuccessResponse(c, "Student report retrieved successfully", fiber.Map{
		"student": fiber.Map{
			"id":            student.ID,
			"user_id":       student.UserID,
			"student_id":    student.StudentID,
			"program_study": student.ProgramStudy,
			"full_name":     student.User.FullName,
			"email":         student.User.Email,
		},
		"summary": fiber.Map{
			"total_achievements":    totalAchievements,
			"verified_achievements": statusCounts[string(models.StatusVerified)],
			"pending_achievements":  statusCounts[string(models.StatusSubmitted)],
			"rejected_achievements": statusCounts[string(models.StatusRejected)],
			"draft_achievements":    statusCounts[string(models.StatusDraft)],
		},
		"achievements_by_type":  typeCounts,
		"achievements_by_level": fiber.Map{}, // Can be expanded later if needed
	})
}

package service

import (
	"context"
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

func (s *reportService) GetStatistics(c *fiber.Ctx) error {
	// Get statistics by type
	typeCounts, _ := s.achievementRepo.CountByType(context.Background())

	// Count students and lecturers
	_, totalStudents, _ := s.studentRepo.FindAll(0, 0)
	_, totalLecturers, _ := s.lecturerRepo.FindAll(0, 0)

	return utils.SuccessResponse(c, "Statistics retrieved successfully", fiber.Map{
		"achievements": typeCounts,
		"students":     totalStudents,
		"lecturers":    totalLecturers,
	})
}

func (s *reportService) GetStudentReport(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid student ID")
	}

	student, err := s.studentRepo.FindByUserID(id)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusNotFound, "Student not found")
	}

	// Get student's achievements count by type
	typeCounts, _ := s.achievementRepo.CountByStudentIDAndType(context.Background(), id.String())

	return utils.SuccessResponse(c, "Student report retrieved successfully", fiber.Map{
		"student":      student,
		"achievements": typeCounts,
	})
}

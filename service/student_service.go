package service

import (
	"strconv"
	"student-achievement-system/repository"
	"student-achievement-system/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type StudentService interface {
	ListStudents(c *fiber.Ctx) error
	GetStudent(c *fiber.Ctx) error
	GetStudentAchievements(c *fiber.Ctx) error
	AssignAdvisor(c *fiber.Ctx) error
}

type LecturerService interface {
	ListLecturers(c *fiber.Ctx) error
	GetAdvisees(c *fiber.Ctx) error
}

type AssignAdvisorRequest struct {
	AdvisorID string `json:"advisor_id" validate:"required,uuid"`
}

type studentService struct {
	studentRepo        repository.StudentRepository
	lecturerRepo       repository.LecturerRepository
	achievementRefRepo repository.AchievementReferenceRepository
}

type lecturerService struct {
	lecturerRepo repository.LecturerRepository
	studentRepo  repository.StudentRepository
}

func NewStudentService(
	studentRepo repository.StudentRepository,
	lecturerRepo repository.LecturerRepository,
	achievementRefRepo repository.AchievementReferenceRepository,
) StudentService {
	return &studentService{
		studentRepo:        studentRepo,
		lecturerRepo:       lecturerRepo,
		achievementRefRepo: achievementRefRepo,
	}
}

func NewLecturerService(
	lecturerRepo repository.LecturerRepository,
	studentRepo repository.StudentRepository,
) LecturerService {
	return &lecturerService{
		lecturerRepo: lecturerRepo,
		studentRepo:  studentRepo,
	}
}

// Student Service Methods
func (s *studentService) ListStudents(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	offset := (page - 1) * limit
	students, total, err := s.studentRepo.FindAll(offset, limit)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to fetch students")
	}

	return utils.SuccessResponse(c, "Students retrieved successfully", fiber.Map{
		"students": students,
		"pagination": fiber.Map{
			"page":  page,
			"limit": limit,
			"total": total,
		},
	})
}

func (s *studentService) GetStudent(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid student ID")
	}

	student, err := s.studentRepo.FindByUserID(id)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusNotFound, "Student not found")
	}

	return utils.SuccessResponse(c, "Student retrieved successfully", student)
}

func (s *studentService) GetStudentAchievements(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid student ID")
	}

	student, err := s.studentRepo.FindByUserID(id)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusNotFound, "Student not found")
	}

	// Get achievements from MongoDB would be done via achievement repository
	// For now, return student info
	return utils.SuccessResponse(c, "Student achievements retrieved", fiber.Map{
		"student": student,
		"note":    "Achievement data will be fetched from MongoDB",
	})
}

func (s *studentService) AssignAdvisor(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid student ID")
	}

	var req AssignAdvisorRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body")
	}

	advisorID, err := uuid.Parse(req.AdvisorID)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid advisor ID")
	}

	// Verify lecturer exists
	_, err = s.lecturerRepo.FindByUserID(advisorID)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusNotFound, "Advisor not found")
	}

	student, err := s.studentRepo.FindByUserID(id)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusNotFound, "Student not found")
	}

	student.AdvisorID = &advisorID
	if err := s.studentRepo.Update(student); err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to assign advisor")
	}

	student, _ = s.studentRepo.FindByUserID(id)
	return utils.SuccessResponse(c, "Advisor assigned successfully", student)
}

// Lecturer Service Methods
func (s *lecturerService) ListLecturers(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	offset := (page - 1) * limit
	lecturers, total, err := s.lecturerRepo.FindAll(offset, limit)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to fetch lecturers")
	}

	return utils.SuccessResponse(c, "Lecturers retrieved successfully", fiber.Map{
		"lecturers": lecturers,
		"pagination": fiber.Map{
			"page":  page,
			"limit": limit,
			"total": total,
		},
	})
}

func (s *lecturerService) GetAdvisees(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid lecturer ID")
	}

	advisees, err := s.studentRepo.FindByAdvisorID(id)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to fetch advisees")
	}

	return utils.SuccessResponse(c, "Advisees retrieved successfully", advisees)
}

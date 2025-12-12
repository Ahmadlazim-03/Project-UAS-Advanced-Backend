package service

import (
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
// ListStudents godoc
// @Summary      List all students
// @Description  Get paginated list of students with their academic information
// @Tags         Students
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        page   query    int  false  "Page number (default 1)"
// @Param        limit  query    int  false  "Items per page (default 10, max 100)"
// @Success      200 {object} map[string]interface{} "List of students with pagination"
// @Failure      401 {object} map[string]interface{} "Unauthorized"
// @Failure      403 {object} map[string]interface{} "Forbidden"
// @Failure      500 {object} map[string]interface{} "Internal server error"
// @Router       /students [get]
func (s *studentService) ListStudents(c *fiber.Ctx) error {
	// Get pagination parameters
	pagination := utils.GetPaginationParams(c)

	students, total, err := s.studentRepo.FindAll(pagination.Offset, pagination.Limit)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to fetch students")
	}

	return utils.PaginatedResponse(c, fiber.Map{
		"students": students,
	}, total, pagination.Page, pagination.Limit)
}

// GetStudent godoc
// @Summary      Get student by ID
// @Description  Get detailed information of a specific student
// @Tags         Students
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id  path     string  true  "Student ID (UUID)"
// @Success      200 {object} map[string]interface{} "Student details"
// @Failure      400 {object} map[string]interface{} "Invalid student ID"
// @Failure      401 {object} map[string]interface{} "Unauthorized"
// @Failure      404 {object} map[string]interface{} "Student not found"
// @Router       /students/{id} [get]
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

// GetStudentAchievements godoc
// @Summary      Get student achievements
// @Description  Get all achievements of a specific student
// @Tags         Students
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id  path     string  true  "Student ID (UUID)"
// @Success      200 {object} map[string]interface{} "Student achievements"
// @Failure      400 {object} map[string]interface{} "Invalid student ID"
// @Failure      401 {object} map[string]interface{} "Unauthorized"
// @Failure      404 {object} map[string]interface{} "Student not found"
// @Router       /students/{id}/achievements [get]
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

// AssignAdvisor godoc
// @Summary      Assign advisor to student
// @Description  Assign a lecturer as advisor to a student
// @Tags         Students
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id       path     string                true  "Student ID (UUID)"
// @Param        advisor  body     AssignAdvisorRequest  true  "Advisor assignment data"
// @Success      200 {object} map[string]interface{} "Advisor assigned successfully"
// @Failure      400 {object} map[string]interface{} "Invalid student ID or advisor ID"
// @Failure      401 {object} map[string]interface{} "Unauthorized"
// @Failure      403 {object} map[string]interface{} "Forbidden"
// @Failure      404 {object} map[string]interface{} "Student or advisor not found"
// @Failure      500 {object} map[string]interface{} "Internal server error"
// @Router       /students/{id}/assign-advisor [post]
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

	// Verify lecturer exists using lecturer ID (not user_id)
	_, err = s.lecturerRepo.FindByID(advisorID)
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
// ListLecturers godoc
// @Summary      List all lecturers
// @Description  Get paginated list of lecturers with their department information
// @Tags         Lecturers
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        page   query    int  false  "Page number (default 1)"
// @Param        limit  query    int  false  "Items per page (default 10, max 100)"
// @Success      200 {object} map[string]interface{} "List of lecturers with pagination"
// @Failure      401 {object} map[string]interface{} "Unauthorized"
// @Failure      403 {object} map[string]interface{} "Forbidden"
// @Failure      500 {object} map[string]interface{} "Internal server error"
// @Router       /lecturers [get]
func (s *lecturerService) ListLecturers(c *fiber.Ctx) error {
	// Get pagination parameters
	pagination := utils.GetPaginationParams(c)

	lecturers, total, err := s.lecturerRepo.FindAll(pagination.Offset, pagination.Limit)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to fetch lecturers")
	}

	return utils.PaginatedResponse(c, fiber.Map{
		"lecturers": lecturers,
	}, total, pagination.Page, pagination.Limit)
}

// GetAdvisees godoc
// @Summary      Get lecturer's advisees
// @Description  Get all students advised by a specific lecturer
// @Tags         Lecturers
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id  path     string  true  "Lecturer ID (UUID)"
// @Success      200 {object} map[string]interface{} "List of advisees"
// @Failure      400 {object} map[string]interface{} "Invalid lecturer ID"
// @Failure      401 {object} map[string]interface{} "Unauthorized"
// @Failure      404 {object} map[string]interface{} "Lecturer not found"
// @Failure      500 {object} map[string]interface{} "Internal server error"
// @Router       /lecturers/{id}/advisees [get]
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

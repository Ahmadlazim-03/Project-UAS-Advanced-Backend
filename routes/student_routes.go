package routes

import (
	"github.com/Ahmadlazim-03/Project-UAS-Advanced-Backend/database"
	"github.com/Ahmadlazim-03/Project-UAS-Advanced-Backend/middleware"
	"github.com/Ahmadlazim-03/Project-UAS-Advanced-Backend/models"
	"github.com/Ahmadlazim-03/Project-UAS-Advanced-Backend/services"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// CreateStudent godoc
// @Summary Create a new student
// @Description Create a new student record
// @Tags Students
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body models.Student true "Student Data"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Router /students [post]
func CreateStudent() fiber.Handler {
	return func(c *fiber.Ctx) error {
		var student models.Student
		if err := c.BodyParser(&student); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "Invalid request body"})
		}

		// Generate UUID if not provided
		if student.ID == uuid.Nil {
			student.ID = uuid.New()
		}

		if err := database.DB.Create(&student).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": err.Error()})
		}

		// Reload with relations
		database.DB.Preload("User").Preload("Advisor").First(&student, student.ID)

		return c.Status(fiber.StatusCreated).JSON(fiber.Map{"status": "success", "data": student})
	}
}

// GetAllStudents godoc
// @Summary Get all students
// @Description Get a list of all students
// @Tags Students
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /students [get]
func GetAllStudents() fiber.Handler {
	return func(c *fiber.Ctx) error {
		var students []models.Student
		if err := database.DB.Preload("User").Preload("Advisor.User").Find(&students).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": err.Error()})
		}
		return c.JSON(fiber.Map{"status": "success", "data": students})
	}
}

// GetStudentByID godoc
// @Summary Get student by ID
// @Description Get details of a specific student
// @Tags Students
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Student ID"
// @Success 200 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /students/{id} [get]
func GetStudentByID() fiber.Handler {
	return func(c *fiber.Ctx) error {
		id := c.Params("id")
		var student models.Student
		if err := database.DB.Preload("User").Preload("Advisor").Where("id = ?", id).First(&student).Error; err != nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"status": "error", "message": "Student not found"})
		}
		return c.JSON(fiber.Map{"status": "success", "data": student})
	}
}

// UpdateStudent godoc
// @Summary Update student
// @Description Update student information
// @Tags Students
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Student ID"
// @Param request body models.Student true "Student Data"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Router /students/{id} [put]
func UpdateStudent() fiber.Handler {
	return func(c *fiber.Ctx) error {
		id := c.Params("id")
		
		var student models.Student
		if err := database.DB.Where("id = ?", id).First(&student).Error; err != nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"status": "error", "message": "Student not found"})
		}

		var updateData models.Student
		if err := c.BodyParser(&updateData); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "Invalid request body"})
		}

		// Update fields
		if updateData.StudentID != "" {
			student.StudentID = updateData.StudentID
		}
		if updateData.ProgramStudy != "" {
			student.ProgramStudy = updateData.ProgramStudy
		}
		if updateData.AcademicYear != "" {
			student.AcademicYear = updateData.AcademicYear
		}

		if err := database.DB.Save(&student).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": err.Error()})
		}

		// Reload with relations
		database.DB.Preload("User").Preload("Advisor").First(&student, student.ID)

		return c.JSON(fiber.Map{"status": "success", "data": student})
	}
}

// DeleteStudent godoc
// @Summary Delete student
// @Description Delete a student record
// @Tags Students
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Student ID"
// @Success 200 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /students/{id} [delete]
func DeleteStudent() fiber.Handler {
	return func(c *fiber.Ctx) error {
		id := c.Params("id")
		
		var student models.Student
		if err := database.DB.Where("id = ?", id).First(&student).Error; err != nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"status": "error", "message": "Student not found"})
		}

		if err := database.DB.Delete(&student).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": err.Error()})
		}

		return c.JSON(fiber.Map{"status": "success", "message": "Student deleted successfully"})
	}
}

// GetStudentAchievements godoc
// @Summary Get student achievements
// @Description Get all achievements of a specific student
// @Tags Students
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Student ID"
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /students/{id}/achievements [get]
func GetStudentAchievementsHandler(service services.AchievementService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		studentID := c.Params("id")
		
		// Get student user_id from student table
		var student models.Student
		if err := database.DB.Where("id = ?", studentID).First(&student).Error; err != nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"status": "error", "message": "Student not found"})
		}

		achievements, err := service.GetStudentAchievements(student.UserID.String())
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": err.Error()})
		}
		return c.JSON(fiber.Map{"status": "success", "data": achievements})
	}
}

// UpdateStudentAdvisor godoc
// @Summary Update student advisor
// @Description Assign or update advisor for a student
// @Tags Students
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Student ID"
// @Param request body map[string]string true "Advisor Data"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Router /students/{id}/advisor [put]
func UpdateStudentAdvisor() fiber.Handler {
	return func(c *fiber.Ctx) error {
		id := c.Params("id")
		var req struct {
			AdvisorID string `json:"advisorId"`
		}
		if err := c.BodyParser(&req); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "Invalid request body"})
		}

		if err := database.DB.Model(&models.Student{}).Where("id = ?", id).Update("advisor_id", req.AdvisorID).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": err.Error()})
		}

		return c.JSON(fiber.Map{"status": "success", "message": "Student advisor updated successfully"})
	}
}

func SetupStudentRoutes(app *fiber.App, achService services.AchievementService) {
	api := app.Group("/api/v1/students", middleware.Protected())

	api.Post("/", CreateStudent())
	api.Get("/", GetAllStudents())
	api.Get("/:id", GetStudentByID())
	api.Put("/:id", UpdateStudent())
	api.Delete("/:id", DeleteStudent())
	api.Get("/:id/achievements", GetStudentAchievementsHandler(achService))
	api.Put("/:id/advisor", UpdateStudentAdvisor())
}

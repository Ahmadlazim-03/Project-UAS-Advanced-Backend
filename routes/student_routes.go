package routes

import (
	"github.com/Ahmadlazim-03/Project-UAS-Advanced-Backend/database"
	"github.com/Ahmadlazim-03/Project-UAS-Advanced-Backend/middleware"
	"github.com/Ahmadlazim-03/Project-UAS-Advanced-Backend/models"
	"github.com/Ahmadlazim-03/Project-UAS-Advanced-Backend/services"
	"github.com/gofiber/fiber/v2"
)

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
		if err := database.DB.Preload("User").Preload("Advisor").Find(&students).Error; err != nil {
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

	api.Get("/", GetAllStudents())
	api.Get("/:id", GetStudentByID())
	api.Get("/:id/achievements", GetStudentAchievementsHandler(achService))
	api.Put("/:id/advisor", UpdateStudentAdvisor())
}

package routes

import (
	"github.com/Ahmadlazim-03/Project-UAS-Advanced-Backend/database"
	"github.com/Ahmadlazim-03/Project-UAS-Advanced-Backend/middleware"
	"github.com/Ahmadlazim-03/Project-UAS-Advanced-Backend/models"
	"github.com/gofiber/fiber/v2"
)

// GetAllLecturers godoc
// @Summary Get all lecturers
// @Description Get a list of all lecturers
// @Tags Lecturers
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /lecturers [get]
func GetAllLecturers() fiber.Handler {
	return func(c *fiber.Ctx) error {
		var lecturers []models.Lecturer
		if err := database.DB.Preload("User").Find(&lecturers).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": err.Error()})
		}
		return c.JSON(fiber.Map{"status": "success", "data": lecturers})
	}
}

// GetLecturerAdvisees godoc
// @Summary Get lecturer's advisees
// @Description Get all students under a specific lecturer's advisement
// @Tags Lecturers
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Lecturer ID"
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /lecturers/{id}/advisees [get]
func GetLecturerAdvisees() fiber.Handler {
	return func(c *fiber.Ctx) error {
		lecturerID := c.Params("id")
		
		var students []models.Student
		if err := database.DB.Preload("User").Where("advisor_id = ?", lecturerID).Find(&students).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": err.Error()})
		}
		
		return c.JSON(fiber.Map{"status": "success", "data": students})
	}
}

func SetupLecturerRoutes(app *fiber.App) {
	api := app.Group("/api/v1/lecturers", middleware.Protected())

	api.Get("/", GetAllLecturers())
	api.Get("/:id/advisees", GetLecturerAdvisees())
}

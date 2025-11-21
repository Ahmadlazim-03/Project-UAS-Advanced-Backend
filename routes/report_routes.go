package routes

import (
	"github.com/Ahmadlazim-03/Project-UAS-Advanced-Backend/middleware"
	"github.com/Ahmadlazim-03/Project-UAS-Advanced-Backend/repository"
	"github.com/Ahmadlazim-03/Project-UAS-Advanced-Backend/services"
	"github.com/gofiber/fiber/v2"
)

// GetStatistics godoc
// @Summary Get achievement statistics
// @Description Get statistics of achievements
// @Tags Reports
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /reports/statistics [get]
func GetStatistics(service services.ReportService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		stats, err := service.GetStatistics()
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": err.Error()})
		}
		return c.JSON(fiber.Map{"status": "success", "data": stats})
	}
}

// GetStudentReport godoc
// @Summary Get student report
// @Description Get achievement report for a specific student
// @Tags Reports
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Student User ID"
// @Success 200 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /reports/student/{id} [get]
func GetStudentReport(service services.ReportService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		studentID := c.Params("id")
		report, err := service.GetStudentReport(studentID)
		if err != nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"status": "error", "message": err.Error()})
		}
		return c.JSON(fiber.Map{"status": "success", "data": report})
	}
}

func SetupReportRoutes(app *fiber.App) {
	repo := repository.NewReportRepository()
	service := services.NewReportService(repo)

	api := app.Group("/api/v1/reports", middleware.Protected())

	api.Get("/statistics", GetStatistics(service))
	api.Get("/student/:id", GetStudentReport(service))
}

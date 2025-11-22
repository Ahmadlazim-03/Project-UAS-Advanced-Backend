package routes

import (
	"github.com/Ahmadlazim-03/Project-UAS-Advanced-Backend/database"
	"github.com/Ahmadlazim-03/Project-UAS-Advanced-Backend/middleware"
	"github.com/Ahmadlazim-03/Project-UAS-Advanced-Backend/models"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// CreateLecturer godoc
// @Summary Create a new lecturer
// @Description Create a new lecturer record
// @Tags Lecturers
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body models.Lecturer true "Lecturer Data"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Router /lecturers [post]
func CreateLecturer() fiber.Handler {
	return func(c *fiber.Ctx) error {
		       var lecturer models.Lecturer
		       if err := c.BodyParser(&lecturer); err != nil {
			       return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "Invalid request body"})
		       }

		       // Validasi lecturer_id (NIP) tidak boleh kosong
		       if lecturer.LecturerID == "" {
			       return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "Lecturer ID (NIP) is required"})
		       }

		       // Generate UUID if not provided
		       if lecturer.ID == uuid.Nil {
			       lecturer.ID = uuid.New()
		       }

		       if err := database.DB.Create(&lecturer).Error; err != nil {
			       return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": err.Error()})
		       }

		       // Reload with relations
		       database.DB.Preload("User").First(&lecturer, lecturer.ID)

		       return c.Status(fiber.StatusCreated).JSON(fiber.Map{"status": "success", "data": lecturer})
	}
}

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

		// Add advisees count for each lecturer
		for i := range lecturers {
			var count int64
			database.DB.Model(&models.Student{}).Where("advisor_id = ?", lecturers[i].ID).Count(&count)
			// Store count in a custom field (we'll add this to the response)
			lecturers[i].Department = lecturers[i].Department // Keep original
		}

		return c.JSON(fiber.Map{"status": "success", "data": lecturers})
	}
}

// GetLecturerByID godoc
// @Summary Get lecturer by ID
// @Description Get details of a specific lecturer
// @Tags Lecturers
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Lecturer ID"
// @Success 200 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /lecturers/{id} [get]
func GetLecturerByID() fiber.Handler {
	return func(c *fiber.Ctx) error {
		id := c.Params("id")
		var lecturer models.Lecturer
		if err := database.DB.Preload("User").Where("id = ?", id).First(&lecturer).Error; err != nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"status": "error", "message": "Lecturer not found"})
		}
		return c.JSON(fiber.Map{"status": "success", "data": lecturer})
	}
}

// UpdateLecturer godoc
// @Summary Update lecturer
// @Description Update lecturer information
// @Tags Lecturers
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Lecturer ID"
// @Param request body models.Lecturer true "Lecturer Data"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Router /lecturers/{id} [put]
func UpdateLecturer() fiber.Handler {
	return func(c *fiber.Ctx) error {
		id := c.Params("id")
		
		var lecturer models.Lecturer
		if err := database.DB.Where("id = ?", id).First(&lecturer).Error; err != nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"status": "error", "message": "Lecturer not found"})
		}

		var updateData models.Lecturer
		if err := c.BodyParser(&updateData); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "Invalid request body"})
		}

		// Update fields
		if updateData.LecturerID != "" {
			lecturer.LecturerID = updateData.LecturerID
		}
		if updateData.Department != "" {
			lecturer.Department = updateData.Department
		}

		if err := database.DB.Save(&lecturer).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": err.Error()})
		}

		// Reload with relations
		database.DB.Preload("User").First(&lecturer, lecturer.ID)

		return c.JSON(fiber.Map{"status": "success", "data": lecturer})
	}
}

// DeleteLecturer godoc
// @Summary Delete lecturer
// @Description Delete a lecturer record
// @Tags Lecturers
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Lecturer ID"
// @Success 200 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /lecturers/{id} [delete]
func DeleteLecturer() fiber.Handler {
	return func(c *fiber.Ctx) error {
		id := c.Params("id")
		
		var lecturer models.Lecturer
		if err := database.DB.Where("id = ?", id).First(&lecturer).Error; err != nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"status": "error", "message": "Lecturer not found"})
		}

		// Check if lecturer has advisees
		var count int64
		database.DB.Model(&models.Student{}).Where("advisor_id = ?", id).Count(&count)
		if count > 0 {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"status": "error", 
				"message": "Cannot delete lecturer with assigned advisees. Please reassign students first.",
			})
		}

		if err := database.DB.Delete(&lecturer).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": err.Error()})
		}

		return c.JSON(fiber.Map{"status": "success", "message": "Lecturer deleted successfully"})
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

	api.Post("/", CreateLecturer())
	api.Get("/", GetAllLecturers())
	api.Get("/:id", GetLecturerByID())
	api.Put("/:id", UpdateLecturer())
	api.Delete("/:id", DeleteLecturer())
	api.Get("/:id/advisees", GetLecturerAdvisees())
}

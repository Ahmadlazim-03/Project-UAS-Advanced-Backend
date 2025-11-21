package routes

import (
	"fmt"
	"strconv"
	"time"

	"github.com/Ahmadlazim-03/Project-UAS-Advanced-Backend/middleware"
	"github.com/Ahmadlazim-03/Project-UAS-Advanced-Backend/models"
	"github.com/Ahmadlazim-03/Project-UAS-Advanced-Backend/repository"
	"github.com/Ahmadlazim-03/Project-UAS-Advanced-Backend/services"
	"github.com/gofiber/fiber/v2"
)

// CreateAchievement godoc
// @Summary Create achievement
// @Description Create a new achievement (Student)
// @Tags Achievements
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body models.Achievement true "Achievement Data"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Router /achievements [post]
func CreateAchievement(service services.AchievementService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Parse as generic map to handle string dates from frontend
		var achievementData map[string]interface{}
		if err := c.BodyParser(&achievementData); err != nil {
			fmt.Printf("Body parsing error: %v\n", err)
			fmt.Printf("Raw body: %s\n", string(c.Body()))
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "Invalid request body", "error": err.Error()})
		}

		fmt.Printf("Received achievement data: %+v\n", achievementData)

		studentID := c.Locals("user_id").(string)
		
		// Convert map to Achievement struct
		achievement := models.Achievement{
			StudentID:       studentID,
			Title:           getString(achievementData, "title"),
			Description:     getString(achievementData, "description"),
			Category:        getString(achievementData, "category"),
			Organizer:       getString(achievementData, "organizer"),
			CertificateNo:   getString(achievementData, "certificate_number"),
			Points:          getInt(achievementData, "points"),
			AchievementType: getString(achievementData, "achievement_type"),
			Tags:            []string{},
			Attachments:     []models.Attachment{},
		}

		// Parse achievement_date
		if dateStr := getString(achievementData, "achievement_date"); dateStr != "" {
			parsedDate, err := time.Parse("2006-01-02", dateStr)
			if err == nil {
				achievement.AchievementDate = parsedDate
			}
		}

		if err := service.CreateAchievement(studentID, achievement); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": err.Error()})
		}

		return c.JSON(fiber.Map{"status": "success", "message": "Achievement created successfully"})
	}
}

// Helper functions for type conversion
func getString(m map[string]interface{}, key string) string {
	if val, ok := m[key]; ok {
		if str, ok := val.(string); ok {
			return str
		}
	}
	return ""
}

func getInt(m map[string]interface{}, key string) int {
	if val, ok := m[key]; ok {
		switch v := val.(type) {
		case int:
			return v
		case float64:
			return int(v)
		case string:
			if i, err := strconv.Atoi(v); err == nil {
				return i
			}
		}
	}
	return 0
}

// GetStudentAchievements godoc
// @Summary Get student achievements
// @Description Get all achievements for the logged-in student OR all achievements if admin/lecturer
// @Tags Achievements
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param status query string false "Filter by status (draft, submitted, verified, rejected, all)"
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /achievements [get]
func GetStudentAchievements(service services.AchievementService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		role := c.Locals("role").(string)
		userID := c.Locals("user_id").(string)
		status := c.Query("status", "")

		var achievements []map[string]interface{}
		var err error

		// Admin and Dosen Wali can see all achievements
		if role == "Admin" || role == "Dosen Wali" {
			achievements, err = service.GetAllAchievements(status)
		} else {
			// Mahasiswa can only see their own achievements
			achievements, err = service.GetStudentAchievements(userID)
		}

		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": err.Error()})
		}
		return c.JSON(fiber.Map{"status": "success", "data": achievements})
	}
}

// GetAchievement godoc
// @Summary Get achievement details
// @Description Get details of a specific achievement
// @Tags Achievements
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Achievement ID"
// @Success 200 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /achievements/{id} [get]
func GetAchievement(service services.AchievementService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id := c.Params("id")
		achievement, err := service.GetAchievement(id)
		if err != nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"status": "error", "message": "Achievement not found"})
		}
		return c.JSON(fiber.Map{"status": "success", "data": achievement})
	}
}

// UpdateAchievement godoc
// @Summary Update achievement
// @Description Update an existing achievement (Draft only)
// @Tags Achievements
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Achievement ID"
// @Param request body models.Achievement true "Achievement Data"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Router /achievements/{id} [put]
func UpdateAchievement(service services.AchievementService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id := c.Params("id")
		var achievement models.Achievement
		if err := c.BodyParser(&achievement); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "Invalid request body"})
		}

		if err := service.UpdateAchievement(id, achievement); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": err.Error()})
		}

		return c.JSON(fiber.Map{"status": "success", "message": "Achievement updated successfully"})
	}
}

// DeleteAchievement godoc
// @Summary Delete achievement
// @Description Delete an achievement (Draft only)
// @Tags Achievements
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Achievement ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Router /achievements/{id} [delete]
func DeleteAchievement(service services.AchievementService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id := c.Params("id")
		if err := service.DeleteAchievement(id); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": err.Error()})
		}
		return c.JSON(fiber.Map{"status": "success", "message": "Achievement deleted successfully"})
	}
}

// SubmitAchievement godoc
// @Summary Submit achievement
// @Description Submit an achievement for verification
// @Tags Achievements
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Achievement ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Router /achievements/{id}/submit [post]
func SubmitAchievement(service services.AchievementService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id := c.Params("id")
		if err := service.SubmitAchievement(id); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": err.Error()})
		}
		return c.JSON(fiber.Map{"status": "success", "message": "Achievement submitted successfully"})
	}
}

// VerifyAchievementHandler godoc
// @Summary Verify achievement
// @Description Verify a submitted achievement (Dosen Wali)
// @Tags Achievements
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Achievement ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Router /achievements/{id}/verify [post]
func VerifyAchievementHandler(service services.AchievementService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id := c.Params("id")
		verifierID := c.Locals("user_id").(string)

		if err := service.VerifyAchievement(id, verifierID); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": err.Error()})
		}

		return c.JSON(fiber.Map{"status": "success", "message": "Achievement verified successfully"})
	}
}

// RejectAchievementHandler godoc
// @Summary Reject achievement
// @Description Reject a submitted achievement (Dosen Wali)
// @Tags Achievements
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Achievement ID"
// @Param request body RejectRequest true "Rejection note"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Router /achievements/{id}/reject [post]
func RejectAchievementHandler(service services.AchievementService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id := c.Params("id")
		verifierID := c.Locals("user_id").(string)

		var req RejectRequest
		if err := c.BodyParser(&req); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "Invalid request body"})
		}

		if err := service.RejectAchievement(id, verifierID, req.Note); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": err.Error()})
		}

		return c.JSON(fiber.Map{"status": "success", "message": "Achievement rejected successfully"})
	}
}

// GetAchievementHistory godoc
// @Summary Get achievement history
// @Description Get status change history for an achievement
// @Tags Achievements
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Achievement ID"
// @Success 200 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /achievements/{id}/history [get]
func GetAchievementHistory(service services.AchievementService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id := c.Params("id")
		history, err := service.GetAchievementHistory(id)
		if err != nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"status": "error", "message": err.Error()})
		}
		return c.JSON(fiber.Map{"status": "success", "data": history})
	}
}

// UploadAttachment godoc
// @Summary Upload achievement attachment
// @Description Upload file attachment for an achievement
// @Tags Achievements
// @Accept multipart/form-data
// @Produce json
// @Security BearerAuth
// @Param id path string true "Achievement ID"
// @Param file formData file true "Attachment file"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Router /achievements/{id}/attachments [post]
func UploadAttachment(service services.AchievementService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id := c.Params("id")
		
		file, err := c.FormFile("file")
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "No file uploaded"})
		}

		fileURL, err := service.UploadAttachment(id, file)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": err.Error()})
		}

		return c.JSON(fiber.Map{"status": "success", "message": "File uploaded successfully", "data": fiber.Map{"url": fileURL}})
	}
}

func SetupAchievementRoutes(app *fiber.App) {
	repo := repository.NewAchievementRepository()
	service := services.NewAchievementService(repo)

	api := app.Group("/api/v1/achievements", middleware.Protected())

	api.Post("/", CreateAchievement(service))
	api.Get("/", GetStudentAchievements(service))
	api.Get("/:id", GetAchievement(service))
	api.Put("/:id", UpdateAchievement(service))
	api.Delete("/:id", DeleteAchievement(service))
	api.Post("/:id/submit", SubmitAchievement(service))
	api.Post("/:id/verify", VerifyAchievementHandler(service))
	api.Post("/:id/reject", RejectAchievementHandler(service))
	api.Get("/:id/history", GetAchievementHistory(service))
	api.Post("/:id/attachments", UploadAttachment(service))
}

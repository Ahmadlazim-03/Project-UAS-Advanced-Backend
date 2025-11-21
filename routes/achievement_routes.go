package routes

import (
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
		var achievement models.Achievement
		if err := c.BodyParser(&achievement); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "Invalid request body"})
		}

		studentID := c.Locals("user_id").(string)
		if err := service.CreateAchievement(studentID, achievement); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": err.Error()})
		}

		return c.JSON(fiber.Map{"status": "success", "message": "Achievement created successfully"})
	}
}

// GetStudentAchievements godoc
// @Summary Get student achievements
// @Description Get all achievements for the logged-in student
// @Tags Achievements
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /achievements [get]
func GetStudentAchievements(service services.AchievementService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		studentID := c.Locals("user_id").(string)
		achievements, err := service.GetStudentAchievements(studentID)
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
}

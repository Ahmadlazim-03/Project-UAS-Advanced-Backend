package routes

import (
	"github.com/Ahmadlazim-03/Project-UAS-Advanced-Backend/middleware"
	"github.com/Ahmadlazim-03/Project-UAS-Advanced-Backend/repository"
	"github.com/Ahmadlazim-03/Project-UAS-Advanced-Backend/services"
	"github.com/gofiber/fiber/v2"
)

type RejectRequest struct {
	Note string `json:"note"`
}

// GetPendingVerifications godoc
// @Summary Get pending verifications
// @Description Get list of achievements pending verification (Advisor)
// @Tags Verification
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /verification/pending [get]
func GetPendingVerifications(service services.VerificationService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		advisorID := c.Locals("user_id").(string)

		achievements, err := service.GetPendingVerifications(advisorID)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": err.Error()})
		}

		return c.JSON(fiber.Map{"status": "success", "data": achievements})
	}
}

// VerifyAchievement godoc
// @Summary Verify achievement
// @Description Verify a submitted achievement (Advisor)
// @Tags Verification
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Achievement ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Router /verification/{id}/verify [post]
func VerifyAchievement(service services.VerificationService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id := c.Params("id")
		verifierID := c.Locals("user_id").(string)

		if err := service.VerifyAchievement(id, verifierID); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": err.Error()})
		}

		return c.JSON(fiber.Map{"status": "success", "message": "Achievement verified successfully"})
	}
}

// RejectAchievement godoc
// @Summary Reject achievement
// @Description Reject a submitted achievement (Advisor)
// @Tags Verification
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Achievement ID"
// @Param request body RejectRequest true "Rejection Note"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Router /verification/{id}/reject [post]
func RejectAchievement(service services.VerificationService) fiber.Handler {
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

func SetupVerificationRoutes(app *fiber.App) {
	repo := repository.NewAchievementRepository()
	service := services.NewVerificationService(repo)

	api := app.Group("/api/v1/verification", middleware.Protected())

	api.Get("/pending", GetPendingVerifications(service))
	api.Post("/:id/verify", VerifyAchievement(service))
	api.Post("/:id/reject", RejectAchievement(service))
}

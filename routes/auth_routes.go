package routes

import (
	"github.com/Ahmadlazim-03/Project-UAS-Advanced-Backend/repository"
	"github.com/Ahmadlazim-03/Project-UAS-Advanced-Backend/services"
	"github.com/gofiber/fiber/v2"
)

func SetupAuthRoutes(app *fiber.App) {
	userRepo := repository.NewUserRepository()
	authService := services.NewAuthService(userRepo)

	api := app.Group("/api/v1/auth")

	api.Post("/login", func(c *fiber.Ctx) error {
		type LoginRequest struct {
			Username string `json:"username"`
			Password string `json:"password"`
		}

		var req LoginRequest
		if err := c.BodyParser(&req); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "Invalid request body"})
		}

		token, user, err := authService.Login(req.Username, req.Password)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"status": "error", "message": err.Error()})
		}

		return c.JSON(fiber.Map{
			"status": "success",
			"data": fiber.Map{
				"token": token,
				"user": fiber.Map{
					"id":       user.ID,
					"username": user.Username,
					"fullName": user.FullName,
					"role":     user.Role.Name,
				},
			},
		})
	})

	// Temporary register route for seeding/testing
	api.Post("/register", func(c *fiber.Ctx) error {
		type RegisterRequest struct {
			Username string `json:"username"`
			Email    string `json:"email"`
			Password string `json:"password"`
			FullName string `json:"fullName"`
			RoleName string `json:"roleName"`
		}

		var req RegisterRequest
		if err := c.BodyParser(&req); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "Invalid request body"})
		}

		err := authService.Register(req.Username, req.Email, req.Password, req.FullName, req.RoleName)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": err.Error()})
		}

		return c.JSON(fiber.Map{"status": "success", "message": "User registered successfully"})
	})
}

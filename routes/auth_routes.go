package routes

import (
	"github.com/Ahmadlazim-03/Project-UAS-Advanced-Backend/middleware"
	"github.com/Ahmadlazim-03/Project-UAS-Advanced-Backend/repository"
	"github.com/Ahmadlazim-03/Project-UAS-Advanced-Backend/services"
	"github.com/gofiber/fiber/v2"
)

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type RegisterRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
	FullName string `json:"fullName"`
	RoleName string `json:"roleName"`
}

// Login godoc
// @Summary Login user
// @Description Authenticate user and return JWT token
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body LoginRequest true "Login Credentials"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Router /auth/login [post]
func Login(authService services.AuthService) fiber.Handler {
	return func(c *fiber.Ctx) error {
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
	}
}

// Register godoc
// @Summary Register user
// @Description Register a new user (for testing purposes)
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body RegisterRequest true "Register Data"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Router /auth/register [post]
func Register(authService services.AuthService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var req RegisterRequest
		if err := c.BodyParser(&req); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "Invalid request body"})
		}

		err := authService.Register(req.Username, req.Email, req.Password, req.FullName, req.RoleName)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": err.Error()})
		}

		return c.JSON(fiber.Map{"status": "success", "message": "User registered successfully"})
	}
}

// GetProfile godoc
// @Summary Get user profile
// @Description Get current logged-in user profile
// @Tags Auth
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Router /auth/profile [get]
func GetProfile(userRepo repository.UserRepository) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userID := c.Locals("user_id").(string)
		
		user, err := userRepo.FindUserByID(userID)
		if err != nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"status": "error", "message": "User not found"})
		}

		return c.JSON(fiber.Map{
			"status": "success",
			"data": fiber.Map{
				"id":       user.ID,
				"username": user.Username,
				"email":    user.Email,
				"fullName": user.FullName,
				"isActive": user.IsActive,
				"role":     user.Role,
			},
		})
	}
}

// Logout godoc
// @Summary Logout user
// @Description Logout user (client should remove token)
// @Tags Auth
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]interface{}
// @Router /auth/logout [post]
func Logout() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// In JWT stateless auth, logout is handled client-side by removing the token
		// Optionally, you can implement token blacklisting here
		return c.JSON(fiber.Map{"status": "success", "message": "Logged out successfully"})
	}
}

func SetupAuthRoutes(app *fiber.App) {
	userRepo := repository.NewUserRepository()
	authService := services.NewAuthService(userRepo)

	api := app.Group("/api/v1/auth")

	api.Post("/login", Login(authService))
	api.Post("/register", Register(authService))
	api.Post("/logout", Logout())
	api.Get("/profile", middleware.Protected(), GetProfile(userRepo))
}

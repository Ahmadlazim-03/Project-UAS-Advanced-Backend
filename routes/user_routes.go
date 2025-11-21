package routes

import (
	"github.com/Ahmadlazim-03/Project-UAS-Advanced-Backend/middleware"
	"github.com/Ahmadlazim-03/Project-UAS-Advanced-Backend/repository"
	"github.com/Ahmadlazim-03/Project-UAS-Advanced-Backend/services"
	"github.com/gofiber/fiber/v2"
)

type CreateUserRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
	FullName string `json:"fullName"`
	RoleName string `json:"roleName"`
}

type UpdateUserRequest struct {
	FullName string `json:"fullName"`
	RoleName string `json:"roleName"`
	IsActive *bool  `json:"isActive"`
}

// GetAllUsers godoc
// @Summary Get all users
// @Description Get a list of all users (Admin only)
// @Tags Users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /users [get]
func GetAllUsers(userService services.UserService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		users, err := userService.GetAllUsers()
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": err.Error()})
		}
		return c.JSON(fiber.Map{"status": "success", "data": users})
	}
}

// GetUserByID godoc
// @Summary Get user by ID
// @Description Get details of a specific user
// @Tags Users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "User ID"
// @Success 200 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /users/{id} [get]
func GetUserByID(userService services.UserService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id := c.Params("id")
		user, err := userService.GetUserByID(id)
		if err != nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"status": "error", "message": "User not found"})
		}
		return c.JSON(fiber.Map{"status": "success", "data": user})
	}
}

// CreateUser godoc
// @Summary Create a new user
// @Description Create a new user (Admin only)
// @Tags Users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body CreateUserRequest true "User Data"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Router /users [post]
func CreateUser(userService services.UserService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var req CreateUserRequest
		if err := c.BodyParser(&req); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "Invalid request body"})
		}

		err := userService.CreateUser(req.Username, req.Email, req.Password, req.FullName, req.RoleName)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": err.Error()})
		}

		return c.JSON(fiber.Map{"status": "success", "message": "User created successfully"})
	}
}

// UpdateUser godoc
// @Summary Update user
// @Description Update user details
// @Tags Users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "User ID"
// @Param request body UpdateUserRequest true "Update Data"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Router /users/{id} [put]
func UpdateUser(userService services.UserService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id := c.Params("id")
		var req UpdateUserRequest
		if err := c.BodyParser(&req); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "Invalid request body"})
		}

		err := userService.UpdateUser(id, req.FullName, req.RoleName, req.IsActive)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": err.Error()})
		}

		return c.JSON(fiber.Map{"status": "success", "message": "User updated successfully"})
	}
}

// DeleteUser godoc
// @Summary Delete user
// @Description Delete a user
// @Tags Users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "User ID"
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /users/{id} [delete]
func DeleteUser(userService services.UserService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id := c.Params("id")
		err := userService.DeleteUser(id)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": err.Error()})
		}
		return c.JSON(fiber.Map{"status": "success", "message": "User deleted successfully"})
	}
}

// UpdateUserRole godoc
// @Summary Update user role
// @Description Update user's role (Admin only)
// @Tags Users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "User ID"
// @Param request body map[string]string true "Role Data"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Router /users/{id}/role [put]
func UpdateUserRole(userService services.UserService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id := c.Params("id")
		var req struct {
			RoleName string `json:"roleName"`
		}
		if err := c.BodyParser(&req); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "Invalid request body"})
		}

		err := userService.UpdateUser(id, "", req.RoleName, nil)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": err.Error()})
		}

		return c.JSON(fiber.Map{"status": "success", "message": "User role updated successfully"})
	}
}

// ToggleUserStatus godoc
// @Summary Toggle user active status
// @Description Toggle user's active/inactive status
// @Tags Users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "User ID"
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /users/{id}/toggle-status [patch]
func ToggleUserStatus(userService services.UserService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id := c.Params("id")
		err := userService.ToggleUserStatus(id)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": err.Error()})
		}
		return c.JSON(fiber.Map{"status": "success", "message": "User status toggled successfully"})
	}
}

func SetupUserRoutes(app *fiber.App) {
	userRepo := repository.NewUserRepository()
	userService := services.NewUserService(userRepo)

	api := app.Group("/api/v1/users", middleware.Protected())

	api.Get("/", GetAllUsers(userService))
	api.Get("/:id", GetUserByID(userService))
	api.Post("/", CreateUser(userService))
	api.Put("/:id", UpdateUser(userService))
	api.Delete("/:id", DeleteUser(userService))
	api.Put("/:id/role", UpdateUserRole(userService))
	api.Patch("/:id/toggle-status", ToggleUserStatus(userService))
}

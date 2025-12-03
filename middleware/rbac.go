package middleware

import (
	"student-achievement-system/utils"

	"github.com/gofiber/fiber/v2"
)

// RequirePermission checks if user has the required permission
func RequirePermission(permission string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		user := GetUserFromContext(c)
		if user == nil {
			return utils.ErrorResponse(c, fiber.StatusUnauthorized, "Unauthorized")
		}

		// Check if user has the required permission
		hasPermission := false
		for _, perm := range user.Permissions {
			if perm == permission {
				hasPermission = true
				break
			}
		}

		if !hasPermission {
			return utils.ErrorResponse(c, fiber.StatusForbidden, "Insufficient permissions")
		}

		return c.Next()
	}
}

// RequireAnyPermission checks if user has any of the required permissions
func RequireAnyPermission(permissions ...string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		user := GetUserFromContext(c)
		if user == nil {
			return utils.ErrorResponse(c, fiber.StatusUnauthorized, "Unauthorized")
		}

		// Check if user has any of the required permissions
		for _, requiredPerm := range permissions {
			for _, userPerm := range user.Permissions {
				if userPerm == requiredPerm {
					return c.Next()
				}
			}
		}

		return utils.ErrorResponse(c, fiber.StatusForbidden, "Insufficient permissions")
	}
}

// RequireRole checks if user has the required role
func RequireRole(roleName string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		user := GetUserFromContext(c)
		if user == nil {
			return utils.ErrorResponse(c, fiber.StatusUnauthorized, "Unauthorized")
		}

		if user.RoleName != roleName {
			return utils.ErrorResponse(c, fiber.StatusForbidden, "Insufficient permissions")
		}

		return c.Next()
	}
}

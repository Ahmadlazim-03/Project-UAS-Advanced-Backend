package middleware

import (
	"strings"
	"student-achievement-system/config"
	"student-achievement-system/utils"

	"github.com/gofiber/fiber/v2"
)

// AuthMiddleware validates JWT token
func AuthMiddleware(cfg *config.Config) fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return utils.ErrorResponse(c, fiber.StatusUnauthorized, "Missing authorization header")
		}

		// Extract token from "Bearer <token>"
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			return utils.ErrorResponse(c, fiber.StatusUnauthorized, "Invalid authorization header format")
		}

		token := parts[1]
		claims, err := utils.ValidateToken(token, cfg.JWTSecret)
		if err != nil {
			return utils.ErrorResponse(c, fiber.StatusUnauthorized, "Invalid or expired token")
		}

		// Store claims in context
		c.Locals("user", claims)
		return c.Next()
	}
}

// GetUserFromContext retrieves user claims from context
func GetUserFromContext(c *fiber.Ctx) *utils.JWTClaims {
	user := c.Locals("user")
	if user == nil {
		return nil
	}
	return user.(*utils.JWTClaims)
}

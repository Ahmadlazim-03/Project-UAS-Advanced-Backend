package middleware

import (
	"student-achievement-system/utils"
	"time"

	"github.com/gofiber/fiber/v2"
)

// ErrorHandlerMiddleware handles all errors globally with structured logging
func ErrorHandlerMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Continue to next middleware/handler
		err := c.Next()

		// If there's an error, log it
		if err != nil {
			code := fiber.StatusInternalServerError
			message := "Internal Server Error"

			// Check if it's a Fiber error
			if e, ok := err.(*fiber.Error); ok {
				code = e.Code
				message = e.Message
			}

			// Log the error
			utils.GlobalLogger.Error("Request error", err, map[string]interface{}{
				"method": c.Method(),
				"path":   c.Path(),
				"ip":     c.IP(),
				"status": code,
			})

			// Return error response
			return c.Status(code).JSON(fiber.Map{
				"status":  "error",
				"message": message,
			})
		}

		return nil
	}
}

// RequestLoggerMiddleware logs all incoming requests
func RequestLoggerMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()

		// Continue to next handler
		err := c.Next()

		// Calculate duration
		duration := time.Since(start)

		// Log the request
		utils.GlobalLogger.LogRequest(c, duration)

		return err
	}
}

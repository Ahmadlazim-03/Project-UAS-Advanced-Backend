package utils

import "github.com/gofiber/fiber/v2"

// Response represents a standard API response
type Response struct {
	Status  string      `json:"status"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

// SuccessResponse sends a success response (default 200 OK)
func SuccessResponse(c *fiber.Ctx, message string, data interface{}) error {
	return c.Status(fiber.StatusOK).JSON(Response{
		Status:  "success",
		Message: message,
		Data:    data,
	})
}

// ErrorResponse sends an error response
func ErrorResponse(c *fiber.Ctx, statusCode int, message string) error {
	return c.Status(statusCode).JSON(Response{
		Status: "error",
		Error:  message,
	})
}

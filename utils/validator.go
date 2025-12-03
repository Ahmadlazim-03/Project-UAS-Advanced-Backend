package utils

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

var validate *validator.Validate

func init() {
	validate = validator.New()
}

// ValidateStruct validates a struct and returns formatted error messages
func ValidateStruct(s interface{}) error {
	return validate.Struct(s)
}

// ValidationErrorResponse creates a user-friendly error response for validation errors
func ValidationErrorResponse(c *fiber.Ctx, err error) error {
	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		errors := make(map[string]string)
		
		for _, fieldError := range validationErrors {
			field := strings.ToLower(fieldError.Field())
			
			switch fieldError.Tag() {
			case "required":
				errors[field] = fmt.Sprintf("%s is required", fieldError.Field())
			case "email":
				errors[field] = "Invalid email format"
			case "min":
				errors[field] = fmt.Sprintf("%s must be at least %s characters", fieldError.Field(), fieldError.Param())
			case "max":
				errors[field] = fmt.Sprintf("%s must be at most %s characters", fieldError.Field(), fieldError.Param())
			case "uuid":
				errors[field] = fmt.Sprintf("%s must be a valid UUID", fieldError.Field())
			default:
				errors[field] = fmt.Sprintf("%s is invalid", fieldError.Field())
			}
		}
		
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status": "error",
			"error":  "Validation failed",
			"fields": errors,
		})
	}
	
	return ErrorResponse(c, fiber.StatusBadRequest, err.Error())
}

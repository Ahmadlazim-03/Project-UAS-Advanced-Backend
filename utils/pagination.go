package utils

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
)

// PaginationParams holds pagination parameters
type PaginationParams struct {
	Page   int
	Limit  int
	Offset int
}

// PaginationResponse represents paginated response structure
type PaginationResponse struct {
	Page       int         `json:"page"`
	Limit      int         `json:"limit"`
	Total      int64       `json:"total"`
	TotalPages int         `json:"total_pages"`
	Data       interface{} `json:"data"`
}

// GetPaginationParams extracts and validates pagination parameters from request
func GetPaginationParams(c *fiber.Ctx) PaginationParams {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))

	// Validation
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}
	if limit > 100 {
		limit = 100 // Max limit to prevent abuse
	}

	offset := (page - 1) * limit

	return PaginationParams{
		Page:   page,
		Limit:  limit,
		Offset: offset,
	}
}

// PaginatedResponse creates a standardized paginated response
func PaginatedResponse(c *fiber.Ctx, data interface{}, total int64, page, limit int) error {
	totalPages := int((total + int64(limit) - 1) / int64(limit))

	response := PaginationResponse{
		Page:       page,
		Limit:      limit,
		Total:      total,
		TotalPages: totalPages,
		Data:       data,
	}

	return c.JSON(fiber.Map{
		"status":     "success",
		"pagination": response,
	})
}

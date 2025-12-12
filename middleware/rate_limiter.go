package middleware

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
)

// LoginRateLimiter restricts login attempts to prevent brute force attacks
func LoginRateLimiter() fiber.Handler {
	return limiter.New(limiter.Config{
		Max:        1000,              // Increased for testing (original: 5)
		Expiration: 1 * time.Minute,   // Reduced for testing (original: 15 minutes)
		KeyGenerator: func(c *fiber.Ctx) string {
			// Rate limit by IP address
			return c.IP()
		},
		LimitReached: func(c *fiber.Ctx) error {
			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
				"status":  "error",
				"error":   "Too many login attempts",
				"message": "Please try again in 15 minutes",
			})
		},
		SkipFailedRequests:     false, // Count failed requests
		SkipSuccessfulRequests: false, // Count all requests
		Storage:                nil,   // Use in-memory storage (for production, use Redis)
	})
}

// APIRateLimiter restricts general API calls
func APIRateLimiter() fiber.Handler {
	return limiter.New(limiter.Config{
		Max:        100,              // Maximum 100 requests
		Expiration: 1 * time.Minute,  // Per minute
		KeyGenerator: func(c *fiber.Ctx) string {
			return c.IP()
		},
		LimitReached: func(c *fiber.Ctx) error {
			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
				"status":  "error",
				"error":   "Rate limit exceeded",
				"message": "Too many requests. Please slow down.",
			})
		},
	})
}

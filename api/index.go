package handler

import (
	"log"
	"net/http"
	"os"

	"github.com/Ahmadlazim-03/Project-UAS-Advanced-Backend/database"
	"github.com/Ahmadlazim-03/Project-UAS-Advanced-Backend/models"
	"github.com/Ahmadlazim-03/Project-UAS-Advanced-Backend/routes"
	"github.com/Ahmadlazim-03/Project-UAS-Advanced-Backend/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/adaptor"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/joho/godotenv"
)

var app *fiber.App

func init() {
	// Load environment variables (Vercel will provide these)
	godotenv.Load()

	// Initialize database connections
	database.ConnectPostgres()
	database.ConnectMongo()

	// Auto Migrate
	if err := models.MigrateUsers(database.DB); err != nil {
		log.Println("Failed to migrate users:", err)
	}
	if err := models.MigrateLecturers(database.DB); err != nil {
		log.Println("Failed to migrate lecturers:", err)
	}
	if err := models.MigrateStudents(database.DB); err != nil {
		log.Println("Failed to migrate students:", err)
	}
	if err := models.MigrateAchievements(database.DB); err != nil {
		log.Println("Failed to migrate achievements:", err)
	}

	// Seed Roles
	utils.SeedRoles(database.DB)

	// Initialize Fiber app for serverless
	app = fiber.New(fiber.Config{
		ServerHeader: "Vercel",
		// Disable startup message for serverless
		DisableStartupMessage: true,
	})

	// Middleware
	app.Use(logger.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: "GET,POST,PUT,DELETE,PATCH,OPTIONS",
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
		AllowCredentials: true,
	}))
	app.Use(compress.New(compress.Config{
		Level: compress.LevelBestSpeed,
	}))

	// Health check
	app.Get("/api/v1", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":  "success",
			"message": "Student Achievement System API - Running on Vercel",
			"version": "1.0",
		})
	})

	app.Get("/api/v1/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status": "success",
			"message": "API is healthy",
			"database": "connected",
		})
	})

	// Setup all routes
	routes.SetupAuthRoutes(app)
	routes.SetupUserRoutes(app)
	routes.SetupAchievementRoutes(app)
	routes.SetupVerificationRoutes(app)
	routes.SetupReportRoutes(app)
}

// Handler is the serverless function entry point for Vercel
func Handler(w http.ResponseWriter, r *http.Request) {
	// Ensure PORT is set for Fiber
	if os.Getenv("PORT") == "" {
		os.Setenv("PORT", "3000")
	}

	// Convert http.Request to Fiber and handle
	err := adaptor.FiberApp(app)(w, r)
	if err != nil {
		log.Printf("Error handling request: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

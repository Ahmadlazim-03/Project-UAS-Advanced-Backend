package main

import (
	"log"
	"os"

	"github.com/Ahmadlazim-03/Project-UAS-Advanced-Backend/database"
	"github.com/Ahmadlazim-03/Project-UAS-Advanced-Backend/models"
	"github.com/Ahmadlazim-03/Project-UAS-Advanced-Backend/repository"
	"github.com/Ahmadlazim-03/Project-UAS-Advanced-Backend/routes"
	"github.com/Ahmadlazim-03/Project-UAS-Advanced-Backend/services"
	"github.com/Ahmadlazim-03/Project-UAS-Advanced-Backend/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/joho/godotenv"

	_ "github.com/Ahmadlazim-03/Project-UAS-Advanced-Backend/docs"
	fiberSwagger "github.com/swaggo/fiber-swagger"
)

// @title Student Achievement System API
// @version 1.0
// @description This is a backend system for managing student achievements.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:3000
// @BasePath /api/v1
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	// Connect to Databases
	database.ConnectPostgres()
	database.ConnectMongo()

	// Auto Migrate
	if err := models.MigrateUsers(database.DB); err != nil {
		log.Fatal("Failed to migrate database:", err)
	}
	if err := models.MigrateLecturers(database.DB); err != nil {
		log.Fatal("Failed to migrate lecturers:", err)
	}
	if err := models.MigrateStudents(database.DB); err != nil {
		log.Fatal("Failed to migrate students:", err)
	}
	if err := models.MigrateAchievements(database.DB); err != nil {
		log.Fatal("Failed to migrate achievements:", err)
	}

	// Seed Roles
	utils.SeedRoles(database.DB)

	// Initialize Fiber app
	app := fiber.New()

	// Middleware
	app.Use(logger.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins:     "*",
		AllowMethods:     "GET,POST,PUT,DELETE,PATCH,OPTIONS",
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization",
		AllowCredentials: false,
		ExposeHeaders:    "Content-Length, Content-Type",
	}))
	app.Use(compress.New(compress.Config{
		Level: compress.LevelBestSpeed,
	}))

	// Swagger
	app.Get("/swagger/*", fiberSwagger.WrapHandler)

	// Health check endpoints
	app.Get("/api/v1", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":  "success",
			"message": "Student Achievement System API",
			"version": "1.0",
		})
	})

	app.Get("/api/v1/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":  "success",
			"message": "API is healthy",
		})
	})

	// Setup all API routes
	routes.SetupAuthRoutes(app)
	routes.SetupUserRoutes(app)
	
	// Achievement routes need repo instance
	achRepo := repository.NewAchievementRepository()
	achService := services.NewAchievementService(achRepo)
	
	routes.SetupAchievementRoutes(app)
	routes.SetupVerificationRoutes(app)
	routes.SetupReportRoutes(app)
	routes.SetupStudentRoutes(app, achService)
	routes.SetupLecturerRoutes(app)

	// Serve Static Files (Frontend) - must be after API routes
	app.Static("/", "./frontend/build", fiber.Static{
		Index: "index.html",
		Browse: false,
	})
	
	// Catch all route for SPA - must be last
	app.Get("/*", func(c *fiber.Ctx) error {
		// Skip if it's an API route
		if len(c.Path()) >= 4 && c.Path()[:4] == "/api" {
			return c.Status(404).JSON(fiber.Map{
				"status": "error",
				"message": "API endpoint not found",
			})
		}
		return c.SendFile("./frontend/build/index.html")
	})

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}
	log.Fatal(app.Listen(":" + port))
}

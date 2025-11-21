package main

import (
	"log"
	"os"

	"github.com/Ahmadlazim-03/Project-UAS-Advanced-Backend/database"
	"github.com/Ahmadlazim-03/Project-UAS-Advanced-Backend/models"
	"github.com/Ahmadlazim-03/Project-UAS-Advanced-Backend/routes"
	"github.com/Ahmadlazim-03/Project-UAS-Advanced-Backend/utils"
	"github.com/gofiber/fiber/v2"
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
	if err := models.MigrateAchievements(database.DB); err != nil {
		log.Fatal("Failed to migrate achievements:", err)
	}

	// Seed Roles
	utils.SeedRoles(database.DB)

	// Initialize Fiber app
	app := fiber.New()

	// Middleware
	app.Use(logger.New())
	app.Use(cors.New())

	// Swagger
	app.Get("/swagger/*", fiberSwagger.WrapHandler)

	// Serve Static Files (Frontend)
	app.Static("/", "./public")

	// Routes
	app.Get("/api/v1", func(c *fiber.Ctx) error {
		return c.SendString("Student Achievement System API")
	})

	routes.SetupAuthRoutes(app)
	routes.SetupUserRoutes(app)
	routes.SetupAchievementRoutes(app)
	routes.SetupVerificationRoutes(app)
	routes.SetupReportRoutes(app)

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}
	log.Fatal(app.Listen(":" + port))
}

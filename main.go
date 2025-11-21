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
)

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

	// Seed Roles
	utils.SeedRoles(database.DB)

	// Initialize Fiber app
	app := fiber.New()

	// Middleware
	app.Use(logger.New())
	app.Use(cors.New())

	// Routes
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Student Achievement System API")
	})

	routes.SetupAuthRoutes(app)

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}
	log.Fatal(app.Listen(":" + port))
}

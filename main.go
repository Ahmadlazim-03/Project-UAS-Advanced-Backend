package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"student-achievement-system/config"
	"student-achievement-system/database"
	_ "student-achievement-system/docs"
	"student-achievement-system/middleware"
	"student-achievement-system/repository"
	"student-achievement-system/routes"
	"student-achievement-system/service"
	"student-achievement-system/utils"

	"github.com/gofiber/fiber/v2"
	fiberCors "github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/helmet"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/swagger"
)

// @title Student Achievement Management System API
// @version 1.0
// @description API for managing student achievements with RBAC
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.email support@example.com

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:3000
// @BasePath /api/v1
// @schemes http https

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Enter token with "Bearer " prefix. Example: "Bearer eyJhbGc..."

func main() {
	// Parse command line flags
	cleanupFlag := flag.Bool("cleanup", false, "Delete all data except admin user")
	cleanupAllFlag := flag.Bool("cleanup-all", false, "Delete ALL data including admin (DANGER!)")
	flag.Parse()

	// Load configuration
	cfg := config.LoadConfig()

	// Initialize databases
	database.ConnectPostgres(cfg)
	database.ConnectMongoDB(cfg)
	defer database.ClosePostgres()
	defer database.CloseMongoDB()

	// Handle cleanup commands
	if *cleanupAllFlag {
		fmt.Println("⚠️  DANGER: You are about to delete ALL data including admin!")
		fmt.Print("Type 'YES' to confirm: ")
		var confirm string
		fmt.Scanln(&confirm)
		if confirm == "YES" {
			if err := database.CleanupAllData(); err != nil {
				log.Fatalf("Cleanup failed: %v", err)
			}
			log.Println("✅ Full cleanup completed. Exiting...")
			return
		} else {
			log.Println("Cleanup cancelled.")
			return
		}
	}

	if *cleanupFlag {
		fmt.Println("⚠️  WARNING: You are about to delete all data except admin user!")
		fmt.Print("Type 'yes' to confirm: ")
		var confirm string
		fmt.Scanln(&confirm)
		if confirm == "yes" {
			if err := database.CleanupAllDataExceptAdmin(); err != nil {
				log.Fatalf("Cleanup failed: %v", err)
			}
			log.Println("✅ Cleanup completed. Exiting...")
			return
		} else {
			log.Println("Cleanup cancelled.")
			return
		}
	}

	// Create uploads directory if not exists
	uploadsDir := "./uploads"
	if err := os.MkdirAll(uploadsDir, os.ModePerm); err != nil {
		log.Printf("Warning: Failed to create uploads directory: %v", err)
	} else {
		log.Printf("Uploads directory ready: %s", uploadsDir)
	}

	// Run migrations
	database.Migrate()

	// Seed initial data
	database.SeedData(database.PostgresDB)

	// Initialize repositories
	userRepo := repository.NewUserRepository(database.PostgresDB)
	studentRepo := repository.NewStudentRepository(database.PostgresDB)
	lecturerRepo := repository.NewLecturerRepository(database.PostgresDB)
	roleRepo := repository.NewRoleRepository(database.PostgresDB)
	achievementRefRepo := repository.NewAchievementReferenceRepository(database.PostgresDB)
	achievementRepo := repository.NewAchievementRepository(database.MongoDB)
	notificationRepo := repository.NewNotificationRepository(database.PostgresDB)

	// Initialize services
	authService := service.NewAuthService(userRepo, cfg)
	userService := service.NewUserService(userRepo, studentRepo, lecturerRepo, roleRepo)
	achievementService := service.NewAchievementService(achievementRepo, achievementRefRepo, studentRepo, lecturerRepo)
	verificationService := service.NewVerificationService(achievementRepo, achievementRefRepo, studentRepo, lecturerRepo, notificationRepo)
	studentService := service.NewStudentService(studentRepo, lecturerRepo, achievementRefRepo)
	lecturerService := service.NewLecturerService(lecturerRepo, studentRepo)
	reportService := service.NewReportService(achievementRepo, achievementRefRepo, studentRepo, lecturerRepo)
	fileService := service.NewFileService()
	notificationService := service.NewNotificationService(notificationRepo, userRepo)

	// Create services struct
	services := &routes.Services{
		AuthService:         authService,
		UserService:         userService,
		AchievementService:  achievementService,
		VerificationService: verificationService,
		StudentService:      studentService,
		LecturerService:     lecturerService,
		ReportService:       reportService,
		FileService:         fileService,
		NotificationService: notificationService,
	}

	// Create Fiber app
	app := fiber.New(fiber.Config{
		AppName:      "Student Achievement System",
		ErrorHandler: customErrorHandler,
	})

	// Middleware
	app.Use(recover.New())
	app.Use(middleware.RequestLoggerMiddleware()) // Use custom structured logger
	app.Use(helmet.New())
	app.Use(fiberCors.New(fiberCors.Config{
		AllowOrigins: cfg.CORSOrigin,
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
		AllowMethods: "GET, POST, PUT, DELETE, OPTIONS",
	}))

	// Swagger documentation
	app.Get("/swagger/*", swagger.HandlerDefault)
	app.Get("/swagger", func(c *fiber.Ctx) error {
		return c.Redirect("/swagger/index.html")
	})

	// Static file serving for uploads
	app.Static("/uploads", "./uploads")

	// Health check
	app.Get("/health", func(c *fiber.Ctx) error {
		utils.GlobalLogger.Info("Health check requested", map[string]interface{}{
			"ip": c.IP(),
		})
		return c.JSON(fiber.Map{
			"status":  "success",
			"message": "Server is running",
		})
	})

	// Setup routes
	api := app.Group("/api/" + cfg.APIVersion)
	routes.SetupRoutes(api, services, cfg)

	// Start server
	port := cfg.Port
	if port == "" {
		port = "3000"
	}

	utils.GlobalLogger.Info("Server starting", map[string]interface{}{
		"port":       port,
		"api_version": cfg.APIVersion,
	})
	log.Fatal(app.Listen(":" + port))
}

// customErrorHandler handles errors globally
func customErrorHandler(c *fiber.Ctx, err error) error {
	code := fiber.StatusInternalServerError
	message := "Internal Server Error"

	if e, ok := err.(*fiber.Error); ok {
		code = e.Code
		message = e.Message
	}
	
	// Log the error
	utils.GlobalLogger.Error("Request error", err, map[string]interface{}{
		"method": c.Method(),
		"path":   c.Path(),
		"status": code,
	})

	return c.Status(code).JSON(fiber.Map{
		"status":  "error",
		"message": message,
	})
}

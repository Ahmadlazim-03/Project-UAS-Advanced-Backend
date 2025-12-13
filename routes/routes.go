package routes

import (
	"student-achievement-system/config"
	"student-achievement-system/middleware"
	"student-achievement-system/service"

	"github.com/gofiber/fiber/v2"
)

type Services struct {
	AuthService         service.AuthService
	UserService         service.UserService
	AchievementService  service.AchievementService
	VerificationService service.VerificationService
	StudentService      service.StudentService
	LecturerService     service.LecturerService
	ReportService       service.ReportService
	FileService         service.FileService
	NotificationService service.NotificationService
}

func SetupRoutes(api fiber.Router, services *Services, cfg *config.Config) {
	// Public routes - Authentication with rate limiting
	auth := api.Group("/auth")
	{
		// Apply strict rate limiting to login endpoint (5 attempts per 15 minutes)
		auth.Post("/login", middleware.LoginRateLimiter(), services.AuthService.Login)
		auth.Post("/refresh", services.AuthService.RefreshToken)
		
		// Protected auth routes
		auth.Use(middleware.AuthMiddleware(cfg))
		auth.Post("/logout", services.AuthService.Logout)
		auth.Get("/profile", services.AuthService.GetProfile)
	}

	// Protected routes - require authentication
	api.Use(middleware.AuthMiddleware(cfg))
	
	// Apply general API rate limiting (100 requests per minute)
	api.Use(middleware.APIRateLimiter())

	// Roles routes (Public - needed for user creation)
	roles := api.Group("/roles")
	{
		roles.Get("/", services.UserService.ListRoles)
	}

	// User management routes (Admin only)
	users := api.Group("/users")
	{
		users.Get("/", middleware.RequirePermission("user:manage"), services.UserService.ListUsers)
		users.Get("/deleted", middleware.RequirePermission("user:manage"), services.UserService.ListDeletedUsers)
		users.Get("/:id", middleware.RequirePermission("user:read"), services.UserService.GetUser)
		users.Post("/", middleware.RequirePermission("user:create"), services.UserService.CreateUser)
		users.Put("/:id", middleware.RequirePermission("user:update"), services.UserService.UpdateUser)
		users.Delete("/:id", middleware.RequirePermission("user:delete"), services.UserService.DeleteUser)
		users.Post("/:id/restore", middleware.RequirePermission("user:manage"), services.UserService.RestoreUser)
		users.Delete("/:id/hard-delete", middleware.RequirePermission("user:manage"), services.UserService.HardDeleteUser)
		users.Put("/:id/role", middleware.RequirePermission("user:manage"), services.UserService.AssignRole)
	}

	// Achievement routes
	achievements := api.Group("/achievements")
	{
		achievements.Get("/", middleware.RequirePermission("achievement:read"), services.AchievementService.ListAchievements)
		achievements.Get("/:id", middleware.RequirePermission("achievement:read"), services.AchievementService.GetAchievement)
		achievements.Post("/", middleware.RequirePermission("achievement:create"), services.AchievementService.CreateAchievement)
		achievements.Put("/:id", middleware.RequirePermission("achievement:update"), services.AchievementService.UpdateAchievement)
		achievements.Delete("/:id", middleware.RequirePermission("achievement:delete"), services.AchievementService.DeleteAchievement)
		
		// Status history and attachments
		achievements.Get("/:id/history", middleware.RequirePermission("achievement:read"), services.AchievementService.GetAchievementHistory)
		achievements.Post("/:id/attachments", middleware.RequirePermission("achievement:create"), services.AchievementService.UploadAttachment)
		
		// Submission and verification
		achievements.Post("/:id/submit", middleware.RequirePermission("achievement:update"), services.VerificationService.SubmitForVerification)
		achievements.Post("/:id/verify", middleware.RequirePermission("achievement:verify"), services.VerificationService.VerifyAchievement)
		achievements.Post("/:id/reject", middleware.RequirePermission("achievement:verify"), services.VerificationService.RejectAchievement)
	}

	// Student routes
	students := api.Group("/students")
	{
		students.Get("/", middleware.RequireAnyPermission("user:read", "user:manage"), services.StudentService.ListStudents)
		students.Get("/:id", middleware.RequirePermission("user:read"), services.StudentService.GetStudent)
		students.Get("/:id/achievements", middleware.RequirePermission("achievement:read"), services.StudentService.GetStudentAchievements)
		students.Put("/:id/advisor", middleware.RequirePermission("user:manage"), services.StudentService.AssignAdvisor)
	}

	// Lecturer routes
	lecturers := api.Group("/lecturers")
	{
		lecturers.Get("/", middleware.RequireAnyPermission("user:read", "user:manage"), services.LecturerService.ListLecturers)
		lecturers.Get("/me/advisees", middleware.RequirePermission("achievement:verify"), services.LecturerService.GetMyAdvisees)
		lecturers.Get("/:id/advisees", middleware.RequirePermission("achievement:verify"), services.LecturerService.GetAdvisees)
		lecturers.Get("/advisees/achievements", middleware.RequirePermission("achievement:verify"), services.VerificationService.GetAdviseeAchievements)
	}

	// Report routes
	reports := api.Group("/reports")
	{
		reports.Get("/statistics", middleware.RequirePermission("report:read"), services.ReportService.GetStatistics)
		reports.Get("/students/:id", middleware.RequirePermission("report:read"), services.ReportService.GetStudentReport)
		reports.Get("/top-students", middleware.RequirePermission("report:read"), services.ReportService.GetTopStudents)
		reports.Get("/statistics/period", middleware.RequirePermission("report:read"), services.ReportService.GetStatisticsByPeriod)
		reports.Get("/statistics/competition-levels", middleware.RequirePermission("report:read"), services.ReportService.GetCompetitionLevelDistribution)
	}

	// Notification routes
	notifications := api.Group("/notifications")
	{
		notifications.Get("/", services.NotificationService.GetMyNotifications)
		notifications.Get("/unread", services.NotificationService.GetUnreadNotifications)
		notifications.Get("/unread/count", services.NotificationService.GetUnreadCount)
		notifications.Put("/:id/read", services.NotificationService.MarkAsRead)
		notifications.Put("/read-all", services.NotificationService.MarkAllAsRead)
	}

	// File upload routes
	files := api.Group("/files")
	{
		files.Post("/upload", middleware.RequirePermission("achievement:create"), services.FileService.UploadAchievementFile)
		files.Delete("/:filename", middleware.RequirePermission("achievement:delete"), services.FileService.DeleteAchievementFile)
	}
}

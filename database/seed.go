package database

import (
	"log"
	"student-achievement-system/models"
	"student-achievement-system/utils"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// SeedData creates initial data for testing
func SeedData(db *gorm.DB) {
	log.Println("Starting database seeding...")

	// Create Roles
	roles := []models.Role{
		{ID: uuid.New(), Name: "Admin", Description: "System Administrator"},
		{ID: uuid.New(), Name: "Mahasiswa", Description: "Student"},
		{ID: uuid.New(), Name: "Dosen Wali", Description: "Academic Advisor"},
	}

	for _, role := range roles {
		db.FirstOrCreate(&role, models.Role{Name: role.Name})
	}

	// Create Permissions
	permissions := []models.Permission{
		{ID: uuid.New(), Name: "user:manage", Description: "Manage all users"},
		{ID: uuid.New(), Name: "user:read", Description: "Read user data"},
		{ID: uuid.New(), Name: "user:create", Description: "Create users"},
		{ID: uuid.New(), Name: "user:update", Description: "Update users"},
		{ID: uuid.New(), Name: "user:delete", Description: "Delete users"},
		{ID: uuid.New(), Name: "achievement:create", Description: "Create achievements"},
		{ID: uuid.New(), Name: "achievement:read", Description: "Read achievements"},
		{ID: uuid.New(), Name: "achievement:update", Description: "Update achievements"},
		{ID: uuid.New(), Name: "achievement:delete", Description: "Delete achievements"},
		{ID: uuid.New(), Name: "achievement:verify", Description: "Verify achievements"},
		{ID: uuid.New(), Name: "report:read", Description: "Read reports"},
	}

	for _, perm := range permissions {
		db.FirstOrCreate(&perm, models.Permission{Name: perm.Name})
	}

	// Assign permissions to roles
	var adminRole, studentRole, lecturerRole models.Role
	db.Where("name = ?", "Admin").First(&adminRole)
	db.Where("name = ?", "Mahasiswa").First(&studentRole)
	db.Where("name = ?", "Dosen Wali").First(&lecturerRole)

	var allPerms []models.Permission
	db.Find(&allPerms)

	// Admin gets all permissions
	db.Model(&adminRole).Association("Permissions").Replace(allPerms)

	// Student permissions
	var studentPerms []models.Permission
	db.Where("name IN ?", []string{
		"achievement:create",
		"achievement:read",
		"achievement:update",
		"achievement:delete",
		"user:read",
	}).Find(&studentPerms)
	db.Model(&studentRole).Association("Permissions").Replace(studentPerms)

	// Lecturer permissions
	var lecturerPerms []models.Permission
	db.Where("name IN ?", []string{
		"achievement:read",
		"achievement:verify",
		"user:read",
		"report:read",
	}).Find(&lecturerPerms)
	db.Model(&lecturerRole).Association("Permissions").Replace(lecturerPerms)

	// Create test users
	adminPassword, _ := utils.HashPassword("admin123")
	studentPassword, _ := utils.HashPassword("student123")
	lecturerPassword, _ := utils.HashPassword("lecturer123")

	// Create or get existing users
	var adminUser models.User
	if err := db.Where("username = ?", "admin").First(&adminUser).Error; err != nil {
		adminUser = models.User{
			Username:     "admin",
			Email:        "admin@university.ac.id",
			PasswordHash: adminPassword,
			FullName:     "System Administrator",
			RoleID:       adminRole.ID,
			IsActive:     true,
		}
		db.Create(&adminUser)
	}

	var studentUser models.User
	if err := db.Where("username = ?", "student001").First(&studentUser).Error; err != nil {
		studentUser = models.User{
			Username:     "student001",
			Email:        "student001@university.ac.id",
			PasswordHash: studentPassword,
			FullName:     "John Doe",
			RoleID:       studentRole.ID,
			IsActive:     true,
		}
		db.Create(&studentUser)
	}

	var lecturerUser models.User
	if err := db.Where("username = ?", "lecturer001").First(&lecturerUser).Error; err != nil {
		lecturerUser = models.User{
			Username:     "lecturer001",
			Email:        "lecturer001@university.ac.id",
			PasswordHash: lecturerPassword,
			FullName:     "Dr. Jane Smith",
			RoleID:       lecturerRole.ID,
			IsActive:     true,
		}
		db.Create(&lecturerUser)
	}

	// Create Lecturer profile FIRST (before student references it)
	var lecturer models.Lecturer
	if err := db.Where("lecturer_id = ?", "LEC001").First(&lecturer).Error; err != nil {
		lecturer = models.Lecturer{
			UserID:     lecturerUser.ID,
			LecturerID: "LEC001",
			Department: "Computer Science",
		}
		db.Create(&lecturer)
	}

	// Create Student profile - ensure user exists and lecturer exists
	var student models.Student
	if err := db.Where("student_id = ?", "STU001").First(&student).Error; err != nil {
		student = models.Student{
			UserID:       studentUser.ID,
			StudentID:    "STU001",
			ProgramStudy: "Teknik Informatika",
			AcademicYear: "2024",
			AdvisorID:    &lecturer.ID,
		}
		db.Create(&student)
	}

	log.Println("Database seeding completed!")
	log.Println("Test credentials:")
	log.Println("  Admin    - username: admin, password: admin123")
	log.Println("  Student  - username: student001, password: student123 (STU001)")
	log.Println("  Lecturer - username: lecturer001, password: lecturer123 (LEC001)")
}

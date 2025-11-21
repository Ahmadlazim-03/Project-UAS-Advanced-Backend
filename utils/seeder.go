package utils

import (
	"log"

	"github.com/Ahmadlazim-03/Project-UAS-Advanced-Backend/models"
	"gorm.io/gorm"
)

func SeedRoles(db *gorm.DB) {
	// Seed Permissions
	seedPermissions(db)
	
	// Seed Roles
	roles := []string{"Admin", "Mahasiswa", "Dosen Wali"}

	for _, roleName := range roles {
		var role models.Role
		if err := db.Where("name = ?", roleName).First(&role).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				newRole := models.Role{Name: roleName}
				if err := db.Create(&newRole).Error; err != nil {
					log.Printf("Failed to create role %s: %v", roleName, err)
				} else {
					log.Printf("Role %s created", roleName)
				}
			}
		}
	}
	
	// Assign permissions to roles
	assignRolePermissions(db)
	
	// Seed sample lecturer and students for testing
	seedTestData(db)
}

func seedPermissions(db *gorm.DB) {
	permissions := []models.Permission{
		{Name: "achievement:create", Resource: "achievement", Action: "create", Description: "Create achievement"},
		{Name: "achievement:read", Resource: "achievement", Action: "read", Description: "Read achievement"},
		{Name: "achievement:update", Resource: "achievement", Action: "update", Description: "Update achievement"},
		{Name: "achievement:delete", Resource: "achievement", Action: "delete", Description: "Delete achievement"},
		{Name: "achievement:verify", Resource: "achievement", Action: "verify", Description: "Verify achievement"},
		{Name: "user:manage", Resource: "user", Action: "manage", Description: "Manage users"},
		{Name: "student:manage", Resource: "student", Action: "manage", Description: "Manage students"},
		{Name: "lecturer:manage", Resource: "lecturer", Action: "manage", Description: "Manage lecturers"},
		{Name: "report:view", Resource: "report", Action: "view", Description: "View reports"},
	}

	for _, perm := range permissions {
		var existing models.Permission
		if err := db.Where("name = ?", perm.Name).First(&existing).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				if err := db.Create(&perm).Error; err != nil {
					log.Printf("Failed to create permission %s: %v", perm.Name, err)
				} else {
					log.Printf("Permission %s created", perm.Name)
				}
			}
		}
	}
}

func assignRolePermissions(db *gorm.DB) {
	// Admin - all permissions
	var adminRole models.Role
	db.Preload("Permissions").Where("name = ?", "Admin").First(&adminRole)
	
	var allPermissions []models.Permission
	db.Find(&allPermissions)
	
	if len(adminRole.Permissions) == 0 {
		db.Model(&adminRole).Association("Permissions").Append(&allPermissions)
		log.Println("Assigned all permissions to Admin role")
	}
	
	// Mahasiswa - achievement CRUD
	var mahasiswaRole models.Role
	db.Preload("Permissions").Where("name = ?", "Mahasiswa").First(&mahasiswaRole)
	
	var mahasiswaPerms []models.Permission
	db.Where("name IN ?", []string{
		"achievement:create",
		"achievement:read",
		"achievement:update",
		"achievement:delete",
		"report:view",
	}).Find(&mahasiswaPerms)
	
	if len(mahasiswaRole.Permissions) == 0 {
		db.Model(&mahasiswaRole).Association("Permissions").Append(&mahasiswaPerms)
		log.Println("Assigned permissions to Mahasiswa role")
	}
	
	// Dosen Wali - read and verify
	var dosenRole models.Role
	db.Preload("Permissions").Where("name = ?", "Dosen Wali").First(&dosenRole)
	
	var dosenPerms []models.Permission
	db.Where("name IN ?", []string{
		"achievement:read",
		"achievement:verify",
		"student:manage",
		"report:view",
	}).Find(&dosenPerms)
	
	if len(dosenRole.Permissions) == 0 {
		db.Model(&dosenRole).Association("Permissions").Append(&dosenPerms)
		log.Println("Assigned permissions to Dosen Wali role")
	}
}

func seedTestData(db *gorm.DB) {
	// Check if admin user already exists
	var adminCount int64
	db.Model(&models.User{}).Where("username = ?", "admin").Count(&adminCount)
	if adminCount > 0 {
		return // Already seeded
	}

	hashedPassword, _ := HashPassword("password123")

	// Create admin user
	adminRole := models.Role{}
	db.Where("name = ?", "Admin").First(&adminRole)
	
	adminUser := models.User{
		Username:     "admin",
		Email:        "admin@example.com",
		PasswordHash: hashedPassword,
		FullName:     "Administrator",
		RoleID:       adminRole.ID,
		IsActive:     true,
	}
	if err := db.Create(&adminUser).Error; err == nil {
		log.Println("Admin user created: admin/password123")
	}

	// Create a test lecturer (Dosen Wali)
	dosenRole := models.Role{}
	db.Where("name = ?", "Dosen Wali").First(&dosenRole)

	dosenUser := models.User{
		Username:     "dosenwali",
		Email:        "dosen@example.com",
		PasswordHash: hashedPassword,
		FullName:     "Dr. Dosen Wali",
		RoleID:       dosenRole.ID,
		IsActive:     true,
	}
	if err := db.Create(&dosenUser).Error; err == nil {
		lecturer := models.Lecturer{
			UserID:     dosenUser.ID,
			LecturerID: "LEC001",
			Department: "Teknik Informatika",
		}
		db.Create(&lecturer)
		log.Println("Test lecturer created: dosenwali/password123")
	}
	
	// Create a test student (Mahasiswa)
	mahasiswaRole := models.Role{}
	db.Where("name = ?", "Mahasiswa").First(&mahasiswaRole)

	mahasiswaUser := models.User{
		Username:     "mahasiswa",
		Email:        "mahasiswa@example.com",
		PasswordHash: hashedPassword,
		FullName:     "Mahasiswa Test",
		RoleID:       mahasiswaRole.ID,
		IsActive:     true,
	}
	if err := db.Create(&mahasiswaUser).Error; err == nil {
		// Get lecturer for advisor
		var lecturer models.Lecturer
		db.First(&lecturer)
		
		student := models.Student{
			UserID:       mahasiswaUser.ID,
			StudentID:    "STD001",
			ProgramStudy: "Teknik Informatika",
			AcademicYear: "2024",
			AdvisorID:    &lecturer.ID,
		}
		db.Create(&student)
		log.Println("Test student created: mahasiswa/password123")
	}
}

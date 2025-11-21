package utils

import (
	"log"

	"github.com/Ahmadlazim-03/Project-UAS-Advanced-Backend/models"
	"gorm.io/gorm"
)

func SeedRoles(db *gorm.DB) {
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
	
	// Seed sample lecturer and students for testing
	seedTestData(db)
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

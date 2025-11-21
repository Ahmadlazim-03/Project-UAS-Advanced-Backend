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
	// Check if test data already exists
	var count int64
	db.Model(&models.Lecturer{}).Count(&count)
	if count > 0 {
		return // Already seeded
	}

	// Create a test lecturer (Dosen Wali)
	hashedPassword, _ := HashPassword("password123")
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
}

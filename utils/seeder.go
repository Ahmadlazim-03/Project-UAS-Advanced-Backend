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
}

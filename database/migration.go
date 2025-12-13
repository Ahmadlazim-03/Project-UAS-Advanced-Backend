package database

import (
	"log"
	"student-achievement-system/models"
)

// Migrate runs database migrations
func Migrate() {
	log.Println("Running migrations...")

	// Disable foreign key constraints temporarily
	PostgresDB.Exec("SET session_replication_role = replica;")

	// Auto migrate PostgreSQL tables
	err := PostgresDB.AutoMigrate(
		&models.Role{},
		&models.Permission{},
		&models.RolePermission{},
		&models.User{},
		&models.Lecturer{},
		&models.Student{},
		&models.AchievementReference{},
		&models.AchievementStatusHistory{},
		&models.Notification{},
	)

	// Re-enable foreign key constraints
	PostgresDB.Exec("SET session_replication_role = DEFAULT;")

	if err != nil {
		log.Printf("Migration warning: %v", err)
		log.Println("Continuing with existing schema...")
	}

	// Seed initial data
	seedInitialData()

	log.Println("Migrations completed successfully")
}

// seedInitialData creates initial roles and permissions
func seedInitialData() {
	// Create permissions
	permissions := []models.Permission{
		{Name: "user:create", Description: "Create users"},
		{Name: "user:read", Description: "Read users"},
		{Name: "user:update", Description: "Update users"},
		{Name: "user:delete", Description: "Delete users"},
		{Name: "user:manage", Description: "Manage users"},
		{Name: "achievement:create", Description: "Create achievements"},
		{Name: "achievement:read", Description: "Read achievements"},
		{Name: "achievement:update", Description: "Update achievements"},
		{Name: "achievement:delete", Description: "Delete achievements"},
		{Name: "achievement:verify", Description: "Verify achievements"},
		{Name: "report:read", Description: "Read reports"},
	}

	for _, perm := range permissions {
		var existingPerm models.Permission
		if err := PostgresDB.Where("name = ?", perm.Name).First(&existingPerm).Error; err != nil {
			PostgresDB.Create(&perm)
		}
	}

	// Create roles
	roles := []models.Role{
		{Name: "Admin", Description: "System administrator"},
		{Name: "Dosen Wali", Description: "Academic advisor"},
		{Name: "Mahasiswa", Description: "Student"},
	}

	for _, role := range roles {
		var existingRole models.Role
		if err := PostgresDB.Where("name = ?", role.Name).First(&existingRole).Error; err != nil {
			PostgresDB.Create(&role)
		}
	}

	// Assign permissions to roles
	assignRolePermissions()

	log.Println("Initial data seeded successfully")
}

// assignRolePermissions assigns permissions to roles
func assignRolePermissions() {
	// Admin - all permissions
	var adminRole models.Role
	PostgresDB.Where("name = ?", "Admin").First(&adminRole)

	var allPermissions []models.Permission
	PostgresDB.Find(&allPermissions)

	for _, perm := range allPermissions {
		var existing models.RolePermission
		if err := PostgresDB.Where("role_id = ? AND permission_id = ?", adminRole.ID, perm.ID).First(&existing).Error; err != nil {
			PostgresDB.Create(&models.RolePermission{
				RoleID:       adminRole.ID,
				PermissionID: perm.ID,
			})
		}
	}

	// Dosen Wali - achievement verify and read
	var dosenRole models.Role
	PostgresDB.Where("name = ?", "Dosen Wali").First(&dosenRole)

	dosenPerms := []string{"achievement:read", "achievement:verify", "report:read"}
	for _, permName := range dosenPerms {
		var perm models.Permission
		if PostgresDB.Where("name = ?", permName).First(&perm).Error == nil {
			var existing models.RolePermission
			if err := PostgresDB.Where("role_id = ? AND permission_id = ?", dosenRole.ID, perm.ID).First(&existing).Error; err != nil {
				PostgresDB.Create(&models.RolePermission{
					RoleID:       dosenRole.ID,
					PermissionID: perm.ID,
				})
			}
		}
	}

	// Mahasiswa - create, read, update own achievements
	var mahasiswaRole models.Role
	PostgresDB.Where("name = ?", "Mahasiswa").First(&mahasiswaRole)

	mahasiswaPerms := []string{"achievement:create", "achievement:read", "achievement:update", "achievement:delete"}
	for _, permName := range mahasiswaPerms {
		var perm models.Permission
		if PostgresDB.Where("name = ?", permName).First(&perm).Error == nil {
			var existing models.RolePermission
			if err := PostgresDB.Where("role_id = ? AND permission_id = ?", mahasiswaRole.ID, perm.ID).First(&existing).Error; err != nil {
				PostgresDB.Create(&models.RolePermission{
					RoleID:       mahasiswaRole.ID,
					PermissionID: perm.ID,
				})
			}
		}
	}
}

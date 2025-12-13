package database

import (
	"context"
	"log"
	"student-achievement-system/models"

	"go.mongodb.org/mongo-driver/bson"
)

// CleanupAllDataExceptAdmin menghapus semua data kecuali user admin
func CleanupAllDataExceptAdmin() error {
	log.Println("Starting database cleanup...")
	log.Println("⚠️  WARNING: This will delete all data except admin user!")

	// 1. Delete all MongoDB achievements
	log.Println("Deleting all achievements from MongoDB...")
	achievementCollection := MongoDB.Collection("achievements")
	deleteResult, err := achievementCollection.DeleteMany(context.Background(), bson.M{})
	if err != nil {
		log.Printf("Error deleting MongoDB achievements: %v", err)
		return err
	}
	log.Printf("✓ Deleted %d achievements from MongoDB", deleteResult.DeletedCount)

	// 2. Delete all achievement_references from PostgreSQL
	log.Println("Deleting all achievement_references from PostgreSQL...")
	result := PostgresDB.Exec("DELETE FROM achievement_references")
	if result.Error != nil {
		log.Printf("Error deleting achievement_references: %v", result.Error)
		return result.Error
	}
	log.Printf("✓ Deleted %d achievement_references from PostgreSQL", result.RowsAffected)

	// 3. Delete all students
	log.Println("Deleting all students...")
	result = PostgresDB.Exec("DELETE FROM students")
	if result.Error != nil {
		log.Printf("Error deleting students: %v", result.Error)
		return result.Error
	}
	log.Printf("✓ Deleted %d students from PostgreSQL", result.RowsAffected)

	// 4. Delete all lecturers
	log.Println("Deleting all lecturers...")
	result = PostgresDB.Exec("DELETE FROM lecturers")
	if result.Error != nil {
		log.Printf("Error deleting lecturers: %v", result.Error)
		return result.Error
	}
	log.Printf("✓ Deleted %d lecturers from PostgreSQL", result.RowsAffected)

	// 5. Delete all users except admin
	log.Println("Deleting all users except admin...")
	
	// First, get admin user ID
	var adminUser models.User
	if err := PostgresDB.Where("username = ?", "admin").First(&adminUser).Error; err != nil {
		log.Printf("⚠️  Admin user not found: %v", err)
	} else {
		// Delete all users except admin
		result = PostgresDB.Exec("DELETE FROM users WHERE id != ?", adminUser.ID)
		if result.Error != nil {
			log.Printf("Error deleting users: %v", result.Error)
			return result.Error
		}
		log.Printf("✓ Deleted %d users from PostgreSQL (kept admin)", result.RowsAffected)
	}

	// 6. Reset sequences (optional, untuk reset auto-increment)
	log.Println("Resetting sequences...")
	PostgresDB.Exec("ALTER SEQUENCE students_id_seq RESTART WITH 1")
	PostgresDB.Exec("ALTER SEQUENCE lecturers_id_seq RESTART WITH 1")
	PostgresDB.Exec("ALTER SEQUENCE achievement_references_id_seq RESTART WITH 1")
	log.Println("✓ Sequences reset")

	log.Println("✅ Database cleanup completed successfully!")
	log.Println("Remaining data:")
	log.Println("  - Admin user: admin / admin123")
	log.Println("  - Roles and permissions (unchanged)")
	
	return nil
}

// CleanupAllData menghapus SEMUA data termasuk admin (use with caution!)
func CleanupAllData() error {
	log.Println("Starting FULL database cleanup...")
	log.Println("⚠️  DANGER: This will delete ALL data including admin!")

	// MongoDB
	log.Println("Deleting all achievements from MongoDB...")
	achievementCollection := MongoDB.Collection("achievements")
	deleteResult, err := achievementCollection.DeleteMany(context.Background(), bson.M{})
	if err != nil {
		log.Printf("Error deleting MongoDB achievements: %v", err)
		return err
	}
	log.Printf("✓ Deleted %d achievements from MongoDB", deleteResult.DeletedCount)

	// PostgreSQL - order matters due to foreign keys
	tables := []string{
		"achievement_references",
		"students",
		"lecturers",
		"users",
	}

	for _, table := range tables {
		log.Printf("Deleting all from %s...", table)
		result := PostgresDB.Exec("DELETE FROM " + table)
		if result.Error != nil {
			log.Printf("Error deleting from %s: %v", table, result.Error)
			return result.Error
		}
		log.Printf("✓ Deleted %d rows from %s", result.RowsAffected, table)
	}

	// Reset sequences
	log.Println("Resetting sequences...")
	PostgresDB.Exec("ALTER SEQUENCE students_id_seq RESTART WITH 1")
	PostgresDB.Exec("ALTER SEQUENCE lecturers_id_seq RESTART WITH 1")
	PostgresDB.Exec("ALTER SEQUENCE achievement_references_id_seq RESTART WITH 1")
	PostgresDB.Exec("ALTER SEQUENCE users_id_seq RESTART WITH 1")
	log.Println("✓ Sequences reset")

	log.Println("✅ Full database cleanup completed!")
	log.Println("⚠️  All data has been deleted. Run migrations to recreate initial data.")
	
	return nil
}

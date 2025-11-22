package repository

import (
	"context"

	"github.com/Ahmadlazim-03/Project-UAS-Advanced-Backend/database"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"gorm.io/gorm"
)

type ReportRepository interface {
	GetAchievementStatistics(userID string, role string) (map[string]interface{}, error)
	GetStudentReport(studentID string) (map[string]interface{}, error)
}

type reportRepository struct {
	pgDB    *gorm.DB
	mongoDB *mongo.Database
}

func NewReportRepository() ReportRepository {
	return &reportRepository{
		pgDB:    database.DB,
		mongoDB: database.MongoDb,
	}
}

func (r *reportRepository) GetAchievementStatistics(userID string, role string) (map[string]interface{}, error) {
	// Postgres: Count by status
	var query string
	var args []interface{}

	if role == "Mahasiswa" {
		// Find student by user_id
		var studentID string
		err := r.pgDB.Raw("SELECT id FROM students WHERE user_id = ?", userID).Scan(&studentID).Error
		if err != nil {
			return nil, err
		}
		query = "SELECT status, count(*) as count FROM achievement_references WHERE student_id = ? GROUP BY status"
		args = []interface{}{studentID}
	} else {
		query = "SELECT status, count(*) as count FROM achievement_references GROUP BY status"
		args = []interface{}{}
	}

	rows, err := r.pgDB.Raw(query, args...).Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	statusMap := make(map[string]int)
	for rows.Next() {
		var status string
		var count int
		rows.Scan(&status, &count)
		statusMap[status] = count
	}

	// User statistics by role
	userRows, err := r.pgDB.Raw(`
		SELECT r.name as role_name, count(*) as count 
		FROM users u 
		JOIN roles r ON u.role_id = r.id 
		GROUP BY r.name
	`).Rows()
	if err != nil {
		return nil, err
	}
	defer userRows.Close()

	userMap := make(map[string]int)
	for userRows.Next() {
		var roleName string
		var count int
		userRows.Scan(&roleName, &count)
		userMap[roleName] = count
	}

	// Mongo: Count by Type
	collection := r.mongoDB.Collection("achievements")
	pipeline := mongo.Pipeline{
		{{Key: "$group", Value: bson.D{{Key: "_id", Value: "$achievementType"}, {Key: "count", Value: bson.D{{Key: "$sum", Value: 1}}}}}},
	}
	cursor, err := collection.Aggregate(context.Background(), pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	type TypeCount struct {
		ID    string `bson:"_id"`
		Count int    `bson:"count"`
	}
	var typeCounts []TypeCount
	if err = cursor.All(context.Background(), &typeCounts); err != nil {
		return nil, err
	}

	typeMap := make(map[string]int)
	for _, tc := range typeCounts {
		typeMap[tc.ID] = tc.Count
	}

	return map[string]interface{}{
		"byStatus": statusMap,
		"byType":   typeMap,
		"byRole":   userMap,
	}, nil
}

func (r *reportRepository) GetStudentReport(studentID string) (map[string]interface{}, error) {
	// Count achievements by status for this student
	rows, err := r.pgDB.Raw("SELECT status, count(*) as count FROM achievement_references WHERE student_id = ? GROUP BY status", studentID).Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	statusMap := make(map[string]int)
	totalAchievements := 0
	for rows.Next() {
		var status string
		var count int
		rows.Scan(&status, &count)
		statusMap[status] = count
		totalAchievements += count
	}

	// Get total points from verified achievements
	var totalPoints int
	r.pgDB.Raw("SELECT COALESCE(SUM(points), 0) FROM achievements WHERE student_id = ? AND status = 'verified'", studentID).Scan(&totalPoints)

	return map[string]interface{}{
		"studentID":         studentID,
		"totalAchievements": totalAchievements,
		"byStatus":          statusMap,
		"totalPoints":       totalPoints,
	}, nil
}

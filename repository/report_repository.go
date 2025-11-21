package repository

import (
	"context"

	"github.com/Ahmadlazim-03/Project-UAS-Advanced-Backend/database"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"gorm.io/gorm"
)

type ReportRepository interface {
	GetAchievementStatistics() (map[string]interface{}, error)
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

func (r *reportRepository) GetAchievementStatistics() (map[string]interface{}, error) {
	// Example statistics: Count by status (Postgres) and Count by Type (Mongo)

	// Postgres: Count by status
	// Raw SQL or GORM query
	// SELECT status, count(*) FROM achievement_references GROUP BY status
	// Since we don't have a struct for this result, we can use a map or a temporary struct
	// Let's use a raw query for simplicity
	rows, err := r.pgDB.Raw("SELECT status, count(*) as count FROM achievement_references GROUP BY status").Rows()
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

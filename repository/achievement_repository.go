package repository

import (
	"context"
	"time"

	"github.com/Ahmadlazim-03/Project-UAS-Advanced-Backend/database"
	"github.com/Ahmadlazim-03/Project-UAS-Advanced-Backend/models"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"gorm.io/gorm"
)

type AchievementRepository interface {
	CreateAchievement(achievement *models.Achievement, ref *models.AchievementReference) error
	GetAchievementsByStudentID(studentID string) ([]models.Achievement, []models.AchievementReference, error)
	GetAchievementByID(id string) (*models.Achievement, *models.AchievementReference, error)
	GetPendingVerifications(advisorID string) ([]models.Achievement, []models.AchievementReference, error)
	UpdateAchievement(id string, achievement *models.Achievement) error
	UpdateAchievementStatus(id string, status string) error
	UpdateAchievementReference(id string, ref *models.AchievementReference) error
	VerifyAchievement(id string, verifierID uuid.UUID) error
	RejectAchievement(id string, verifierID uuid.UUID, note string) error
	DeleteAchievement(id string) error
}

type achievementRepository struct {
	pgDB    *gorm.DB
	mongoDB *mongo.Database
}

func NewAchievementRepository() AchievementRepository {
	return &achievementRepository{
		pgDB:    database.DB,
		mongoDB: database.MongoDb,
	}
}

func (r *achievementRepository) CreateAchievement(achievement *models.Achievement, ref *models.AchievementReference) error {
	// Transaction for consistency
	tx := r.pgDB.Begin()

	// Insert into MongoDB
	collection := r.mongoDB.Collection("achievements")
	achievement.CreatedAt = time.Now()
	achievement.UpdatedAt = time.Now()
	res, err := collection.InsertOne(context.Background(), achievement)
	if err != nil {
		tx.Rollback()
		return err
	}

	// Set Mongo ID in Postgres reference
	ref.MongoAchievementID = res.InsertedID.(primitive.ObjectID).Hex()

	// Insert into Postgres
	if err := tx.Create(ref).Error; err != nil {
		// Rollback Mongo insert (manual compensation needed in real world, or just ignore orphan)
		// For simplicity, we just rollback Postgres transaction
		tx.Rollback()
		collection.DeleteOne(context.Background(), bson.M{"_id": res.InsertedID})
		return err
	}

	return tx.Commit().Error
}

func (r *achievementRepository) GetAchievementsByStudentID(studentID string) ([]models.Achievement, []models.AchievementReference, error) {
	var refs []models.AchievementReference
	if err := r.pgDB.Where("student_id = ?", studentID).Find(&refs).Error; err != nil {
		return nil, nil, err
	}

	var achievements []models.Achievement
	collection := r.mongoDB.Collection("achievements")

	for _, ref := range refs {
		objID, _ := primitive.ObjectIDFromHex(ref.MongoAchievementID)
		var achievement models.Achievement
		// Filter out soft deleted achievements
		filter := bson.M{
			"_id": objID,
			"$or": []bson.M{
				{"isDeleted": bson.M{"$exists": false}},
				{"isDeleted": false},
			},
		}
		if err := collection.FindOne(context.Background(), filter).Decode(&achievement); err == nil {
			achievements = append(achievements, achievement)
		}
	}

	return achievements, refs, nil
}

func (r *achievementRepository) GetAchievementByID(id string) (*models.Achievement, *models.AchievementReference, error) {
	var ref models.AchievementReference
	// Assuming ID passed is the Postgres ID (Reference ID)
	if err := r.pgDB.Where("id = ?", id).First(&ref).Error; err != nil {
		return nil, nil, err
	}

	objID, _ := primitive.ObjectIDFromHex(ref.MongoAchievementID)
	var achievement models.Achievement
	collection := r.mongoDB.Collection("achievements")
	if err := collection.FindOne(context.Background(), bson.M{"_id": objID}).Decode(&achievement); err != nil {
		return nil, &ref, err
	}

	return &achievement, &ref, nil
}

func (r *achievementRepository) GetPendingVerifications(advisorID string) ([]models.Achievement, []models.AchievementReference, error) {
	// Get lecturer record by user_id
	var lecturer models.Lecturer
	if err := r.pgDB.Where("user_id = ?", advisorID).First(&lecturer).Error; err != nil {
		return nil, nil, err
	}

	// Get students under this advisor
	var students []models.Student
	if err := r.pgDB.Where("advisor_id = ?", lecturer.ID).Find(&students).Error; err != nil {
		return nil, nil, err
	}

	// Get student IDs
	var studentIDs []uuid.UUID
	for _, s := range students {
		studentIDs = append(studentIDs, s.ID)
	}

	if len(studentIDs) == 0 {
		return []models.Achievement{}, []models.AchievementReference{}, nil
	}

	// Get achievement references for these students with status 'submitted'
	var refs []models.AchievementReference
	if err := r.pgDB.Where("student_id IN ? AND status = ?", studentIDs, "submitted").Find(&refs).Error; err != nil {
		return nil, nil, err
	}

	var achievements []models.Achievement
	collection := r.mongoDB.Collection("achievements")

	for _, ref := range refs {
		objID, _ := primitive.ObjectIDFromHex(ref.MongoAchievementID)
		var achievement models.Achievement
		// Filter out soft deleted achievements
		filter := bson.M{
			"_id": objID,
			"$or": []bson.M{
				{"isDeleted": bson.M{"$exists": false}},
				{"isDeleted": false},
			},
		}
		if err := collection.FindOne(context.Background(), filter).Decode(&achievement); err == nil {
			achievements = append(achievements, achievement)
		}
	}

	return achievements, refs, nil
}

func (r *achievementRepository) UpdateAchievement(id string, achievement *models.Achievement) error {
	// ID here is the Postgres ID
	var ref models.AchievementReference
	if err := r.pgDB.Where("id = ?", id).First(&ref).Error; err != nil {
		return err
	}

	objID, _ := primitive.ObjectIDFromHex(ref.MongoAchievementID)
	collection := r.mongoDB.Collection("achievements")
	
	update := bson.M{
		"$set": bson.M{
			"title":             achievement.Title,
			"description":       achievement.Description,
			"category":          achievement.Category,
			"organizer":         achievement.Organizer,
			"achievementDate":   achievement.AchievementDate,
			"certificateNumber": achievement.CertificateNo,
			"points":            achievement.Points,
			"details":           achievement.Details,
			"tags":              achievement.Tags,
			"updatedAt":         time.Now(),
		},
	}

	_, err := collection.UpdateOne(context.Background(), bson.M{"_id": objID}, update)
	return err
}

func (r *achievementRepository) UpdateAchievementStatus(id string, status string) error {
	updates := map[string]interface{}{
		"status": status,
	}
	
	// Set submitted_at timestamp when status changes to submitted
	if status == "submitted" {
		now := time.Now()
		updates["submitted_at"] = &now
	}
	
	return r.pgDB.Model(&models.AchievementReference{}).Where("id = ?", id).Updates(updates).Error
}

func (r *achievementRepository) UpdateAchievementReference(id string, ref *models.AchievementReference) error {
	return r.pgDB.Model(&models.AchievementReference{}).Where("id = ?", id).Updates(ref).Error
}

func (r *achievementRepository) VerifyAchievement(id string, verifierID uuid.UUID) error {
	now := time.Now()
	return r.pgDB.Model(&models.AchievementReference{}).Where("id = ?", id).Updates(map[string]interface{}{
		"status":      "verified",
		"verified_by": verifierID,
		"verified_at": &now,
	}).Error
}

func (r *achievementRepository) RejectAchievement(id string, verifierID uuid.UUID, note string) error {
	now := time.Now()
	return r.pgDB.Model(&models.AchievementReference{}).Where("id = ?", id).Updates(map[string]interface{}{
		"status":         "rejected",
		"verified_by":    verifierID,
		"verified_at":    &now,
		"rejection_note": note,
	}).Error
}

func (r *achievementRepository) DeleteAchievement(id string) error {
	var ref models.AchievementReference
	if err := r.pgDB.Where("id = ?", id).First(&ref).Error; err != nil {
		return err
	}

	// Soft delete from Postgres (GORM handles this automatically)
	if err := r.pgDB.Delete(&ref).Error; err != nil {
		return err
	}

	// Soft delete from MongoDB
	objID, _ := primitive.ObjectIDFromHex(ref.MongoAchievementID)
	collection := r.mongoDB.Collection("achievements")
	now := time.Now()
	update := bson.M{
		"$set": bson.M{
			"isDeleted": true,
			"deletedAt": now,
			"updatedAt": now,
		},
	}
	_, err := collection.UpdateOne(context.Background(), bson.M{"_id": objID}, update)
	return err
}

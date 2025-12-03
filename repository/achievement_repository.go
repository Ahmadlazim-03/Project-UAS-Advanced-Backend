package repository

import (
	"context"
	"student-achievement-system/models"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type AchievementRepository interface {
	Create(ctx context.Context, achievement *models.Achievement) (string, error)
	FindByID(ctx context.Context, id string) (*models.Achievement, error)
	Update(ctx context.Context, id string, achievement *models.Achievement) error
	Delete(ctx context.Context, id string) error
	CountByType(ctx context.Context) (map[string]int64, error)
	CountByStudentIDAndType(ctx context.Context, studentID string) (map[string]int64, error)
}

type achievementRepository struct {
	collection *mongo.Collection
}

func NewAchievementRepository(db *mongo.Database) AchievementRepository {
	return &achievementRepository{
		collection: db.Collection("achievements"),
	}
}

func (r *achievementRepository) Create(ctx context.Context, achievement *models.Achievement) (string, error) {
	achievement.CreatedAt = time.Now()
	achievement.UpdatedAt = time.Now()

	result, err := r.collection.InsertOne(ctx, achievement)
	if err != nil {
		return "", err
	}

	return result.InsertedID.(primitive.ObjectID).Hex(), nil
}

func (r *achievementRepository) FindByID(ctx context.Context, id string) (*models.Achievement, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var achievement models.Achievement
	err = r.collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&achievement)
	if err != nil {
		return nil, err
	}

	return &achievement, nil
}

func (r *achievementRepository) Update(ctx context.Context, id string, achievement *models.Achievement) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	update := bson.M{
		"$set": bson.M{
			"achievementType": achievement.AchievementType,
			"title":           achievement.Title,
			"description":     achievement.Description,
			"details":         achievement.Details,
			"tags":            achievement.Tags,
			"updatedAt":       time.Now(),
		},
	}

	_, err = r.collection.UpdateOne(ctx, bson.M{"_id": objectID}, update)
	return err
}

func (r *achievementRepository) Delete(ctx context.Context, id string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	_, err = r.collection.DeleteOne(ctx, bson.M{"_id": objectID})
	return err
}

func (r *achievementRepository) CountByType(ctx context.Context) (map[string]int64, error) {
	pipeline := []bson.M{
		{
			"$group": bson.M{
				"_id":   "$achievementType",
				"count": bson.M{"$sum": 1},
			},
		},
	}

	cursor, err := r.collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	typeCounts := make(map[string]int64)
	for cursor.Next(ctx) {
		var result struct {
			ID    string `bson:"_id"`
			Count int64  `bson:"count"`
		}
		if err := cursor.Decode(&result); err == nil {
			typeCounts[result.ID] = result.Count
		}
	}

	return typeCounts, nil
}

func (r *achievementRepository) CountByStudentIDAndType(ctx context.Context, studentID string) (map[string]int64, error) {
	pipeline := []bson.M{
		{
			"$match": bson.M{"studentId": studentID},
		},
		{
			"$group": bson.M{
				"_id":   "$achievementType",
				"count": bson.M{"$sum": 1},
			},
		},
	}

	cursor, err := r.collection.Aggregate(ctx, pipeline, options.Aggregate())
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	typeCounts := make(map[string]int64)
	for cursor.Next(ctx) {
		var result struct {
			ID    string `bson:"_id"`
			Count int64  `bson:"count"`
		}
		if err := cursor.Decode(&result); err == nil {
			typeCounts[result.ID] = result.Count
		}
	}

	return typeCounts, nil
}

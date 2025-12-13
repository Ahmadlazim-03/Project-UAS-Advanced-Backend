package repository

import (
	"student-achievement-system/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AchievementReferenceRepository interface {
	FindByID(id uuid.UUID) (*models.AchievementReference, error)
	FindByMongoID(mongoID string) (*models.AchievementReference, error)
	FindByStudentID(studentID uuid.UUID, offset, limit int, status string) ([]models.AchievementReference, int64, error)
	FindByStudentIDs(studentIDs []uuid.UUID, offset, limit int, status string) ([]models.AchievementReference, int64, error)
	FindAll(offset, limit int, status string) ([]models.AchievementReference, int64, error)
	Create(ref *models.AchievementReference) error
	Update(ref *models.AchievementReference) error
	Delete(id uuid.UUID) error
	CountByStatus() (map[string]int64, error)
	CountByStudentID(studentID uuid.UUID) (map[string]int64, error)
	GetTopStudents(limit int) ([]struct {
		StudentID uuid.UUID
		Count     int64
	}, error)
}

type achievementReferenceRepository struct {
	db *gorm.DB
}

func NewAchievementReferenceRepository(db *gorm.DB) AchievementReferenceRepository {
	return &achievementReferenceRepository{db: db}
}

func (r *achievementReferenceRepository) FindByID(id uuid.UUID) (*models.AchievementReference, error) {
	var ref models.AchievementReference
	err := r.db.Where("id = ?", id).
		Preload("Student.User").
		Preload("Student.Advisor.User").
		Preload("VerifiedByUser").
		First(&ref).Error
	return &ref, err
}

func (r *achievementReferenceRepository) FindByMongoID(mongoID string) (*models.AchievementReference, error) {
	var ref models.AchievementReference
	err := r.db.Where("mongo_achievement_id = ?", mongoID).First(&ref).Error
	if err != nil {
		return nil, err
	}
	return &ref, nil
}

func (r *achievementReferenceRepository) FindByStudentID(studentID uuid.UUID, offset, limit int, status string) ([]models.AchievementReference, int64, error) {
	var refs []models.AchievementReference
	var total int64

	query := r.db.Model(&models.AchievementReference{}).
		Where("student_id = ?", studentID).
		Where("status != ?", models.StatusDeleted)

	if status != "" {
		query = query.Where("status = ?", status)
	}

	query.Count(&total)
	err := query.Offset(offset).Limit(limit).
		Preload("Student.User").
		Order("created_at DESC").
		Find(&refs).Error

	return refs, total, err
}

func (r *achievementReferenceRepository) FindByStudentIDs(studentIDs []uuid.UUID, offset, limit int, status string) ([]models.AchievementReference, int64, error) {
	var refs []models.AchievementReference
	var total int64

	query := r.db.Model(&models.AchievementReference{}).
		Where("student_id IN ?", studentIDs).
		Where("status != ?", models.StatusDeleted)

	if status != "" {
		query = query.Where("status = ?", status)
	}

	query.Count(&total)
	err := query.Offset(offset).Limit(limit).
		Order("created_at DESC").
		Find(&refs).Error

	return refs, total, err
}

func (r *achievementReferenceRepository) FindAll(offset, limit int, status string) ([]models.AchievementReference, int64, error) {
	var refs []models.AchievementReference
	var total int64

	query := r.db.Model(&models.AchievementReference{}).
		Where("status != ?", models.StatusDeleted)

	if status != "" {
		query = query.Where("status = ?", status)
	}

	query.Count(&total)
	err := query.Offset(offset).Limit(limit).
		Preload("Student.User").
		Preload("Student.Advisor.User").
		Preload("VerifiedByUser").
		Order("created_at DESC").
		Find(&refs).Error

	return refs, total, err
}

func (r *achievementReferenceRepository) Create(ref *models.AchievementReference) error {
	return r.db.Create(ref).Error
}

func (r *achievementReferenceRepository) Update(ref *models.AchievementReference) error {
	return r.db.Save(ref).Error
}

func (r *achievementReferenceRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&models.AchievementReference{}, id).Error
}

func (r *achievementReferenceRepository) CountByStatus() (map[string]int64, error) {
	var results []struct {
		Status string
		Count  int64
	}

	err := r.db.Model(&models.AchievementReference{}).
		Select("status, count(*) as count").
		Where("status != ?", models.StatusDeleted).
		Group("status").
		Scan(&results).Error

	if err != nil {
		return nil, err
	}

	counts := make(map[string]int64)
	for _, r := range results {
		counts[r.Status] = r.Count
	}

	return counts, nil
}

func (r *achievementReferenceRepository) CountByStudentID(studentID uuid.UUID) (map[string]int64, error) {
	var results []struct {
		Status string
		Count  int64
	}

	err := r.db.Model(&models.AchievementReference{}).
		Select("status, count(*) as count").
		Where("student_id = ? AND status != ?", studentID, models.StatusDeleted).
		Group("status").
		Scan(&results).Error

	if err != nil {
		return nil, err
	}

	counts := make(map[string]int64)
	for _, r := range results {
		counts[r.Status] = r.Count
	}

	return counts, nil
}

func (r *achievementReferenceRepository) GetTopStudents(limit int) ([]struct {
	StudentID uuid.UUID
	Count     int64
}, error) {
	var results []struct {
		StudentID uuid.UUID
		Count     int64
	}

	err := r.db.Model(&models.AchievementReference{}).
		Select("student_id, count(*) as count").
		Where("status = ?", models.StatusVerified).
		Group("student_id").
		Order("count DESC").
		Limit(limit).
		Scan(&results).Error

	return results, err
}

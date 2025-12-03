package repository

import (
	"student-achievement-system/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type LecturerRepository interface {
	FindByID(id uuid.UUID) (*models.Lecturer, error)
	FindByUserID(userID uuid.UUID) (*models.Lecturer, error)
	FindByLecturerID(lecturerID string) (*models.Lecturer, error)
	FindAll(offset, limit int) ([]models.Lecturer, int64, error)
	Create(lecturer *models.Lecturer) error
	Update(lecturer *models.Lecturer) error
	Delete(id uuid.UUID) error
}

type lecturerRepository struct {
	db *gorm.DB
}

func NewLecturerRepository(db *gorm.DB) LecturerRepository {
	return &lecturerRepository{db: db}
}

func (r *lecturerRepository) FindByID(id uuid.UUID) (*models.Lecturer, error) {
	var lecturer models.Lecturer
	err := r.db.Where("id = ?", id).
		Preload("User").
		First(&lecturer).Error
	return &lecturer, err
}

func (r *lecturerRepository) FindByUserID(userID uuid.UUID) (*models.Lecturer, error) {
	var lecturer models.Lecturer
	err := r.db.Where("user_id = ?", userID).
		Preload("User").
		First(&lecturer).Error
	return &lecturer, err
}

func (r *lecturerRepository) FindByLecturerID(lecturerID string) (*models.Lecturer, error) {
	var lecturer models.Lecturer
	err := r.db.Where("lecturer_id = ?", lecturerID).
		Preload("User").
		First(&lecturer).Error
	return &lecturer, err
}

func (r *lecturerRepository) FindAll(offset, limit int) ([]models.Lecturer, int64, error) {
	var lecturers []models.Lecturer
	var total int64

	r.db.Model(&models.Lecturer{}).Count(&total)
	err := r.db.Offset(offset).Limit(limit).
		Preload("User").
		Find(&lecturers).Error

	return lecturers, total, err
}

func (r *lecturerRepository) Create(lecturer *models.Lecturer) error {
	return r.db.Create(lecturer).Error
}

func (r *lecturerRepository) Update(lecturer *models.Lecturer) error {
	return r.db.Save(lecturer).Error
}

func (r *lecturerRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&models.Lecturer{}, id).Error
}

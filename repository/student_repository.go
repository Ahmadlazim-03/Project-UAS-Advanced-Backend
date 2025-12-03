package repository

import (
	"student-achievement-system/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type StudentRepository interface {
	FindByID(id uuid.UUID) (*models.Student, error)
	FindByUserID(userID uuid.UUID) (*models.Student, error)
	FindByStudentID(studentID string) (*models.Student, error)
	FindAll(offset, limit int) ([]models.Student, int64, error)
	FindByAdvisorID(advisorID uuid.UUID) ([]models.Student, error)
	Create(student *models.Student) error
	Update(student *models.Student) error
	Delete(id uuid.UUID) error
}

type studentRepository struct {
	db *gorm.DB
}

func NewStudentRepository(db *gorm.DB) StudentRepository {
	return &studentRepository{db: db}
}

func (r *studentRepository) FindByID(id uuid.UUID) (*models.Student, error) {
	var student models.Student
	err := r.db.Where("id = ?", id).
		Preload("User").
		Preload("Advisor.User").
		First(&student).Error
	return &student, err
}

func (r *studentRepository) FindByUserID(userID uuid.UUID) (*models.Student, error) {
	var student models.Student
	err := r.db.Where("user_id = ?", userID).
		Preload("User").
		Preload("Advisor.User").
		First(&student).Error
	return &student, err
}

func (r *studentRepository) FindByStudentID(studentID string) (*models.Student, error) {
	var student models.Student
	err := r.db.Where("student_id = ?", studentID).
		Preload("User").
		Preload("Advisor.User").
		First(&student).Error
	return &student, err
}

func (r *studentRepository) FindAll(offset, limit int) ([]models.Student, int64, error) {
	var students []models.Student
	var total int64

	r.db.Model(&models.Student{}).Count(&total)
	err := r.db.Offset(offset).Limit(limit).
		Preload("User").
		Preload("Advisor.User").
		Find(&students).Error

	return students, total, err
}

func (r *studentRepository) FindByAdvisorID(advisorID uuid.UUID) ([]models.Student, error) {
	var students []models.Student
	err := r.db.Where("advisor_id = ?", advisorID).
		Preload("User").
		Find(&students).Error
	return students, err
}

func (r *studentRepository) Create(student *models.Student) error {
	return r.db.Create(student).Error
}

func (r *studentRepository) Update(student *models.Student) error {
	return r.db.Save(student).Error
}

func (r *studentRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&models.Student{}, id).Error
}

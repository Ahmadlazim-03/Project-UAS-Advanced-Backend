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
	DeleteByUserID(userID uuid.UUID) error
}

type lecturerRepository struct {
	db *gorm.DB
}

func NewLecturerRepository(db *gorm.DB) LecturerRepository {
	return &lecturerRepository{db: db}
}

func (r *lecturerRepository) FindByID(id uuid.UUID) (*models.Lecturer, error) {
	var lecturer models.Lecturer
	query := `SELECT * FROM lecturers WHERE id = ? LIMIT 1`
	err := r.db.Raw(query, id).Scan(&lecturer).Error
	if err != nil {
		return nil, err
	}
	
	// Load User
	if lecturer.UserID != uuid.Nil {
		r.db.Raw("SELECT * FROM users WHERE id = ?", lecturer.UserID).Scan(&lecturer.User)
	}
	
	return &lecturer, nil
}

func (r *lecturerRepository) FindByUserID(userID uuid.UUID) (*models.Lecturer, error) {
	var lecturer models.Lecturer
	query := `SELECT * FROM lecturers WHERE user_id = ? LIMIT 1`
	err := r.db.Raw(query, userID).Scan(&lecturer).Error
	if err != nil {
		return nil, err
	}
	
	// Load User
	if lecturer.UserID != uuid.Nil {
		r.db.Raw("SELECT * FROM users WHERE id = ?", lecturer.UserID).Scan(&lecturer.User)
	}
	
	return &lecturer, nil
}

func (r *lecturerRepository) FindByLecturerID(lecturerID string) (*models.Lecturer, error) {
	var lecturer models.Lecturer
	query := `SELECT * FROM lecturers WHERE lecturer_id = ? LIMIT 1`
	err := r.db.Raw(query, lecturerID).Scan(&lecturer).Error
	if err != nil {
		return nil, err
	}
	
	// Load User
	if lecturer.UserID != uuid.Nil {
		r.db.Raw("SELECT * FROM users WHERE id = ?", lecturer.UserID).Scan(&lecturer.User)
	}
	
	return &lecturer, nil
}

func (r *lecturerRepository) FindAll(offset, limit int) ([]models.Lecturer, int64, error) {
	var lecturers []models.Lecturer
	var total int64

	// Count total
	countQuery := `SELECT COUNT(*) FROM lecturers`
	r.db.Raw(countQuery).Scan(&total)

	// Get lecturers
	query := `SELECT * FROM lecturers ORDER BY created_at DESC LIMIT ? OFFSET ?`
	err := r.db.Raw(query, limit, offset).Scan(&lecturers).Error
	
	// Load User for each lecturer
	for i := range lecturers {
		if lecturers[i].UserID != uuid.Nil {
			r.db.Raw("SELECT * FROM users WHERE id = ?", lecturers[i].UserID).Scan(&lecturers[i].User)
		}
	}

	return lecturers, total, err
}

func (r *lecturerRepository) Create(lecturer *models.Lecturer) error {
	// Generate UUID if not set
	if lecturer.ID == uuid.Nil {
		lecturer.ID = uuid.New()
	}
	
	query := `
		INSERT INTO lecturers (id, user_id, lecturer_id, department, created_at)
		VALUES (?, ?, ?, ?, NOW())
	`
	return r.db.Exec(query,
		lecturer.ID, lecturer.UserID, lecturer.LecturerID, 
		lecturer.Department,
	).Error
}

func (r *lecturerRepository) Update(lecturer *models.Lecturer) error {
	query := `
		UPDATE lecturers 
		SET lecturer_id = ?, department = ?
		WHERE id = ?
	`
	return r.db.Exec(query,
		lecturer.LecturerID, lecturer.Department, lecturer.ID,
	).Error
}

func (r *lecturerRepository) Delete(id uuid.UUID) error {
	query := `DELETE FROM lecturers WHERE id = ?`
	return r.db.Exec(query, id).Error
}

func (r *lecturerRepository) DeleteByUserID(userID uuid.UUID) error {
	query := `DELETE FROM lecturers WHERE user_id = ?`
	return r.db.Exec(query, userID).Error
}

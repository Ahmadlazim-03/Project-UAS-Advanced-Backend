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
	DeleteByUserID(userID uuid.UUID) error
}

type studentRepository struct {
	db *gorm.DB
}

func NewStudentRepository(db *gorm.DB) StudentRepository {
	return &studentRepository{db: db}
}

func (r *studentRepository) FindByID(id uuid.UUID) (*models.Student, error) {
	var student models.Student
	query := `SELECT * FROM students WHERE id = ? LIMIT 1`
	err := r.db.Raw(query, id).Scan(&student).Error
	if err != nil {
		return nil, err
	}
	
	// Load User
	if student.UserID != uuid.Nil {
		r.db.Raw("SELECT * FROM users WHERE id = ?", student.UserID).Scan(&student.User)
		if student.User.RoleID != uuid.Nil {
			r.db.Raw("SELECT * FROM roles WHERE id = ?", student.User.RoleID).Scan(&student.User.Role)
		}
	}
	
	// Load Advisor
	if student.AdvisorID != nil && *student.AdvisorID != uuid.Nil {
		var lecturer models.Lecturer
		r.db.Raw("SELECT * FROM lecturers WHERE id = ?", student.AdvisorID).Scan(&lecturer)
		if lecturer.UserID != uuid.Nil {
			r.db.Raw("SELECT * FROM users WHERE id = ?", lecturer.UserID).Scan(&lecturer.User)
		}
		student.Advisor = &lecturer
	}
	
	return &student, nil
}

func (r *studentRepository) FindByUserID(userID uuid.UUID) (*models.Student, error) {
	var student models.Student
	query := `SELECT * FROM students WHERE user_id = ? LIMIT 1`
	err := r.db.Raw(query, userID).Scan(&student).Error
	if err != nil {
		return nil, err
	}
	
	// Load User
	if student.UserID != uuid.Nil {
		r.db.Raw("SELECT * FROM users WHERE id = ?", student.UserID).Scan(&student.User)
		if student.User.RoleID != uuid.Nil {
			r.db.Raw("SELECT * FROM roles WHERE id = ?", student.User.RoleID).Scan(&student.User.Role)
		}
	}
	
	// Load Advisor
	if student.AdvisorID != nil && *student.AdvisorID != uuid.Nil {
		var lecturer models.Lecturer
		r.db.Raw("SELECT * FROM lecturers WHERE id = ?", student.AdvisorID).Scan(&lecturer)
		if lecturer.UserID != uuid.Nil {
			r.db.Raw("SELECT * FROM users WHERE id = ?", lecturer.UserID).Scan(&lecturer.User)
		}
		student.Advisor = &lecturer
	}
	
	return &student, nil
}

func (r *studentRepository) FindByStudentID(studentID string) (*models.Student, error) {
	var student models.Student
	query := `SELECT * FROM students WHERE student_id = ? LIMIT 1`
	err := r.db.Raw(query, studentID).Scan(&student).Error
	if err != nil {
		return nil, err
	}
	
	// Load User
	if student.UserID != uuid.Nil {
		r.db.Raw("SELECT * FROM users WHERE id = ?", student.UserID).Scan(&student.User)
		if student.User.RoleID != uuid.Nil {
			r.db.Raw("SELECT * FROM roles WHERE id = ?", student.User.RoleID).Scan(&student.User.Role)
		}
	}
	
	// Load Advisor
	if student.AdvisorID != nil && *student.AdvisorID != uuid.Nil {
		var lecturer models.Lecturer
		r.db.Raw("SELECT * FROM lecturers WHERE id = ?", student.AdvisorID).Scan(&lecturer)
		if lecturer.UserID != uuid.Nil {
			r.db.Raw("SELECT * FROM users WHERE id = ?", lecturer.UserID).Scan(&lecturer.User)
		}
		student.Advisor = &lecturer
	}
	
	return &student, nil
}

func (r *studentRepository) FindAll(offset, limit int) ([]models.Student, int64, error) {
	var students []models.Student
	var total int64

	// Count total
	countQuery := `SELECT COUNT(*) FROM students`
	r.db.Raw(countQuery).Scan(&total)

	// Get students
	query := `SELECT * FROM students ORDER BY created_at DESC LIMIT ? OFFSET ?`
	err := r.db.Raw(query, limit, offset).Scan(&students).Error
	
	// Load related data for each student
	for i := range students {
		if students[i].UserID != uuid.Nil {
			r.db.Raw("SELECT * FROM users WHERE id = ?", students[i].UserID).Scan(&students[i].User)
			if students[i].User.RoleID != uuid.Nil {
				r.db.Raw("SELECT * FROM roles WHERE id = ?", students[i].User.RoleID).Scan(&students[i].User.Role)
			}
		}
		
		if students[i].AdvisorID != nil && *students[i].AdvisorID != uuid.Nil {
			var lecturer models.Lecturer
			r.db.Raw("SELECT * FROM lecturers WHERE id = ?", students[i].AdvisorID).Scan(&lecturer)
			if lecturer.UserID != uuid.Nil {
				r.db.Raw("SELECT * FROM users WHERE id = ?", lecturer.UserID).Scan(&lecturer.User)
			}
			students[i].Advisor = &lecturer
		}
	}

	return students, total, err
}

func (r *studentRepository) FindByAdvisorID(advisorID uuid.UUID) ([]models.Student, error) {
	var students []models.Student
	query := `SELECT * FROM students WHERE advisor_id = ? ORDER BY created_at DESC`
	err := r.db.Raw(query, advisorID).Scan(&students).Error
	
	// Load User for each student
	for i := range students {
		if students[i].UserID != uuid.Nil {
			r.db.Raw("SELECT * FROM users WHERE id = ?", students[i].UserID).Scan(&students[i].User)
		}
	}
	
	return students, err
}

func (r *studentRepository) Create(student *models.Student) error {
	// Generate UUID if not set
	if student.ID == uuid.Nil {
		student.ID = uuid.New()
	}
	
	query := `
		INSERT INTO students (id, user_id, student_id, program_study, academic_year, advisor_id, created_at)
		VALUES (?, ?, ?, ?, ?, ?, NOW())
	`
	return r.db.Exec(query,
		student.ID, student.UserID, student.StudentID, student.ProgramStudy,
		student.AcademicYear, student.AdvisorID,
	).Error
}

func (r *studentRepository) Update(student *models.Student) error {
	query := `
		UPDATE students 
		SET student_id = ?, program_study = ?, academic_year = ?, advisor_id = ?
		WHERE id = ?
	`
	return r.db.Exec(query,
		student.StudentID, student.ProgramStudy, student.AcademicYear,
		student.AdvisorID, student.ID,
	).Error
}

func (r *studentRepository) Delete(id uuid.UUID) error {
	query := `DELETE FROM students WHERE id = ?`
	return r.db.Exec(query, id).Error
}

func (r *studentRepository) DeleteByUserID(userID uuid.UUID) error {
	query := `DELETE FROM students WHERE user_id = ?`
	return r.db.Exec(query, userID).Error
}

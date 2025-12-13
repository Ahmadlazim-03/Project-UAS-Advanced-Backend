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
	query := `SELECT * FROM achievement_references WHERE id = ? LIMIT 1`
	err := r.db.Raw(query, id).Scan(&ref).Error
	if err != nil {
		return nil, err
	}
	
	// Load Student with User
	if ref.StudentID != uuid.Nil {
		var student models.Student
		r.db.Raw("SELECT * FROM students WHERE id = ?", ref.StudentID).Scan(&student)
		if student.UserID != uuid.Nil {
			r.db.Raw("SELECT * FROM users WHERE id = ?", student.UserID).Scan(&student.User)
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
		ref.Student = &student
	}
	
	// Load VerifiedByUser
	if ref.VerifiedBy != nil && *ref.VerifiedBy != uuid.Nil {
		var user models.User
		r.db.Raw("SELECT * FROM users WHERE id = ?", ref.VerifiedBy).Scan(&user)
		ref.VerifiedByUser = &user
	}
	
	return &ref, nil
}

func (r *achievementReferenceRepository) FindByMongoID(mongoID string) (*models.AchievementReference, error) {
	var ref models.AchievementReference
	query := `SELECT * FROM achievement_references WHERE mongo_achievement_id = ? LIMIT 1`
	err := r.db.Raw(query, mongoID).Scan(&ref).Error
	if err != nil {
		return nil, err
	}
	return &ref, nil
}

func (r *achievementReferenceRepository) FindByStudentID(studentID uuid.UUID, offset, limit int, status string) ([]models.AchievementReference, int64, error) {
	var refs []models.AchievementReference
	var total int64

	// Build count query
	countQuery := `SELECT COUNT(*) FROM achievement_references WHERE student_id = ? AND status != ?`
	countArgs := []interface{}{studentID, models.StatusDeleted}
	
	// Build main query
	mainQuery := `SELECT * FROM achievement_references WHERE student_id = ? AND status != ?`
	mainArgs := []interface{}{studentID, models.StatusDeleted}
	
	if status != "" {
		countQuery += ` AND status = ?`
		countArgs = append(countArgs, status)
		mainQuery += ` AND status = ?`
		mainArgs = append(mainArgs, status)
	}
	
	mainQuery += ` ORDER BY created_at DESC LIMIT ? OFFSET ?`
	mainArgs = append(mainArgs, limit, offset)

	// Execute count
	r.db.Raw(countQuery, countArgs...).Scan(&total)
	
	// Execute main query
	err := r.db.Raw(mainQuery, mainArgs...).Scan(&refs).Error
	
	// Load Student with User for each ref
	for i := range refs {
		if refs[i].StudentID != uuid.Nil {
			var student models.Student
			r.db.Raw("SELECT * FROM students WHERE id = ?", refs[i].StudentID).Scan(&student)
			if student.UserID != uuid.Nil {
				r.db.Raw("SELECT * FROM users WHERE id = ?", student.UserID).Scan(&student.User)
			}
			refs[i].Student = &student
		}
	}

	return refs, total, err
}

func (r *achievementReferenceRepository) FindByStudentIDs(studentIDs []uuid.UUID, offset, limit int, status string) ([]models.AchievementReference, int64, error) {
	var refs []models.AchievementReference
	var total int64

	if len(studentIDs) == 0 {
		return refs, 0, nil
	}

	// Build placeholders for IN clause
	placeholders := ""
	args := []interface{}{}
	for i, id := range studentIDs {
		if i > 0 {
			placeholders += ","
		}
		placeholders += "?"
		args = append(args, id)
	}
	
	// Count query
	countQuery := `SELECT COUNT(*) FROM achievement_references WHERE student_id IN (` + placeholders + `) AND status != ?`
	countArgs := append(args, models.StatusDeleted)
	
	// Main query
	mainQuery := `SELECT * FROM achievement_references WHERE student_id IN (` + placeholders + `) AND status != ?`
	mainArgs := append(args, models.StatusDeleted)
	
	if status != "" {
		countQuery += ` AND status = ?`
		countArgs = append(countArgs, status)
		mainQuery += ` AND status = ?`
		mainArgs = append(mainArgs, status)
	}
	
	mainQuery += ` ORDER BY created_at DESC LIMIT ? OFFSET ?`
	mainArgs = append(mainArgs, limit, offset)

	// Execute
	r.db.Raw(countQuery, countArgs...).Scan(&total)
	err := r.db.Raw(mainQuery, mainArgs...).Scan(&refs).Error

	return refs, total, err
}

func (r *achievementReferenceRepository) FindAll(offset, limit int, status string) ([]models.AchievementReference, int64, error) {
	var refs []models.AchievementReference
	var total int64

	// Build queries
	countQuery := `SELECT COUNT(*) FROM achievement_references WHERE status != ?`
	countArgs := []interface{}{models.StatusDeleted}
	
	mainQuery := `SELECT * FROM achievement_references WHERE status != ?`
	mainArgs := []interface{}{models.StatusDeleted}
	
	if status != "" {
		countQuery += ` AND status = ?`
		countArgs = append(countArgs, status)
		mainQuery += ` AND status = ?`
		mainArgs = append(mainArgs, status)
	}
	
	mainQuery += ` ORDER BY created_at DESC LIMIT ? OFFSET ?`
	mainArgs = append(mainArgs, limit, offset)

	// Execute count
	r.db.Raw(countQuery, countArgs...).Scan(&total)
	
	// Execute main query
	err := r.db.Raw(mainQuery, mainArgs...).Scan(&refs).Error
	
	// Load related data for each ref
	for i := range refs {
		if refs[i].StudentID != uuid.Nil {
			var student models.Student
			r.db.Raw("SELECT * FROM students WHERE id = ?", refs[i].StudentID).Scan(&student)
			if student.UserID != uuid.Nil {
				r.db.Raw("SELECT * FROM users WHERE id = ?", student.UserID).Scan(&student.User)
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
			refs[i].Student = &student
		}
		
		// Load VerifiedByUser
		if refs[i].VerifiedBy != nil && *refs[i].VerifiedBy != uuid.Nil {
			var user models.User
			r.db.Raw("SELECT * FROM users WHERE id = ?", refs[i].VerifiedBy).Scan(&user)
			refs[i].VerifiedByUser = &user
		}
	}

	return refs, total, err
}

func (r *achievementReferenceRepository) Create(ref *models.AchievementReference) error {
	query := `
		INSERT INTO achievement_references 
		(id, mongo_achievement_id, student_id, status, verified_by, verified_at, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`
	return r.db.Exec(query,
		ref.ID, ref.MongoAchievementID, ref.StudentID,
		ref.Status, ref.VerifiedBy, ref.VerifiedAt, ref.CreatedAt, ref.UpdatedAt,
	).Error
}

func (r *achievementReferenceRepository) Update(ref *models.AchievementReference) error {
	query := `
		UPDATE achievement_references 
		SET status = ?, verified_by = ?, verified_at = ?, updated_at = ?
		WHERE id = ?
	`
	return r.db.Exec(query,
		ref.Status, ref.VerifiedBy, ref.VerifiedAt,
		ref.UpdatedAt, ref.ID,
	).Error
}

func (r *achievementReferenceRepository) Delete(id uuid.UUID) error {
	query := `DELETE FROM achievement_references WHERE id = ?`
	return r.db.Exec(query, id).Error
}

func (r *achievementReferenceRepository) CountByStatus() (map[string]int64, error) {
	var results []struct {
		Status string
		Count  int64
	}

	query := `
		SELECT status, COUNT(*) as count 
		FROM achievement_references 
		WHERE status != ? 
		GROUP BY status
	`
	err := r.db.Raw(query, models.StatusDeleted).Scan(&results).Error
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

	query := `
		SELECT status, COUNT(*) as count 
		FROM achievement_references 
		WHERE student_id = ? AND status != ? 
		GROUP BY status
	`
	err := r.db.Raw(query, studentID, models.StatusDeleted).Scan(&results).Error
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

	query := `
		SELECT student_id, COUNT(*) as count 
		FROM achievement_references 
		WHERE status = ? 
		GROUP BY student_id 
		ORDER BY count DESC 
		LIMIT ?
	`
	err := r.db.Raw(query, models.StatusVerified, limit).Scan(&results).Error

	return results, err
}

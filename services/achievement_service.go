package services

import (
	"errors"
	"fmt"
	"mime/multipart"
	"path/filepath"
	"time"

	"github.com/Ahmadlazim-03/Project-UAS-Advanced-Backend/models"
	"github.com/Ahmadlazim-03/Project-UAS-Advanced-Backend/repository"
	"github.com/google/uuid"
)

type AchievementService interface {
	CreateAchievement(studentID string, achievement models.Achievement) error
	GetStudentAchievements(studentID string) ([]map[string]interface{}, error)
	GetAchievement(id string) (map[string]interface{}, error)
	UpdateAchievement(id string, achievement models.Achievement) error
	DeleteAchievement(id string) error
	SubmitAchievement(id string) error
	VerifyAchievement(id string, verifierID string) error
	RejectAchievement(id string, verifierID string, note string) error
	GetAchievementHistory(id string) ([]map[string]interface{}, error)
	UploadAttachment(id string, file *multipart.FileHeader) (string, error)
}

type achievementService struct {
	repo repository.AchievementRepository
}

func NewAchievementService(repo repository.AchievementRepository) AchievementService {
	return &achievementService{
		repo: repo,
	}
}

func (s *achievementService) CreateAchievement(studentID string, achievement models.Achievement) error {
	studentUUID, err := uuid.Parse(studentID)
	if err != nil {
		return errors.New("invalid student ID")
	}

	achievement.StudentID = studentID
	ref := &models.AchievementReference{
		StudentID: studentUUID,
		Status:    "draft",
	}

	return s.repo.CreateAchievement(&achievement, ref)
}

func (s *achievementService) GetStudentAchievements(studentID string) ([]map[string]interface{}, error) {
	achievements, refs, err := s.repo.GetAchievementsByStudentID(studentID)
	if err != nil {
		return nil, err
	}

	var result []map[string]interface{}
	for _, ref := range refs {
		// Find matching achievement (simple index matching might be risky if order differs, but repo implementation aligns them)
		// Better to map by ID, but for now assuming repo returns aligned lists or we iterate to find match.
		// Actually repo implementation iterates refs and fetches mongo docs in order, so they align if mongo fetch succeeds.
		// But if mongo fetch fails for one, the slice lengths differ.
		// Let's make it safer by finding the matching mongo doc.
		
		var mongoData models.Achievement
		found := false
		for _, ach := range achievements {
			if ach.ID.Hex() == ref.MongoAchievementID {
				mongoData = ach
				found = true
				break
			}
		}

		if found {
			result = append(result, map[string]interface{}{
				"id":          ref.ID,
				"status":      ref.Status,
				"submittedAt": ref.SubmittedAt,
				"verifiedAt":  ref.VerifiedAt,
				"data":        mongoData,
			})
		}
	}
	return result, nil
}

func (s *achievementService) GetAchievement(id string) (map[string]interface{}, error) {
	ach, ref, err := s.repo.GetAchievementByID(id)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"id":          ref.ID,
		"status":      ref.Status,
		"submittedAt": ref.SubmittedAt,
		"verifiedAt":  ref.VerifiedAt,
		"rejectionNote": ref.RejectionNote,
		"data":        ach,
	}, nil
}

func (s *achievementService) UpdateAchievement(id string, achievement models.Achievement) error {
	_, ref, err := s.repo.GetAchievementByID(id)
	if err != nil {
		return err
	}

	if ref.Status != "draft" {
		return errors.New("cannot update achievement that is not in draft status")
	}

	return s.repo.UpdateAchievement(id, &achievement)
}

func (s *achievementService) DeleteAchievement(id string) error {
	_, ref, err := s.repo.GetAchievementByID(id)
	if err != nil {
		return err
	}

	if ref.Status != "draft" {
		return errors.New("cannot delete achievement that is not in draft status")
	}

	return s.repo.DeleteAchievement(id)
}

func (s *achievementService) SubmitAchievement(id string) error {
	_, ref, err := s.repo.GetAchievementByID(id)
	if err != nil {
		return err
	}

	if ref.Status != "draft" {
		return errors.New("achievement already submitted")
	}

	return s.repo.UpdateAchievementStatus(id, "submitted")
}

func (s *achievementService) VerifyAchievement(id string, verifierID string) error {
	_, ref, err := s.repo.GetAchievementByID(id)
	if err != nil {
		return err
	}

	if ref.Status != "submitted" {
		return errors.New("only submitted achievements can be verified")
	}

	verifierUUID, err := uuid.Parse(verifierID)
	if err != nil {
		return errors.New("invalid verifier ID")
	}

	ref.Status = "verified"
	ref.VerifiedBy = &verifierUUID
	now := time.Now()
	ref.VerifiedAt = &now

	return s.repo.UpdateAchievementReference(id, ref)
}

func (s *achievementService) RejectAchievement(id string, verifierID string, note string) error {
	_, ref, err := s.repo.GetAchievementByID(id)
	if err != nil {
		return err
	}

	if ref.Status != "submitted" {
		return errors.New("only submitted achievements can be rejected")
	}

	verifierUUID, err := uuid.Parse(verifierID)
	if err != nil {
		return errors.New("invalid verifier ID")
	}

	ref.Status = "rejected"
	ref.VerifiedBy = &verifierUUID
	ref.RejectionNote = note
	now := time.Now()
	ref.VerifiedAt = &now

	return s.repo.UpdateAchievementReference(id, ref)
}

func (s *achievementService) GetAchievementHistory(id string) ([]map[string]interface{}, error) {
	_, ref, err := s.repo.GetAchievementByID(id)
	if err != nil {
		return nil, err
	}

	// Build history timeline
	history := []map[string]interface{}{
		{
			"status":    "draft",
			"timestamp": ref.CreatedAt,
			"note":      "Achievement created",
		},
	}

	if ref.SubmittedAt != nil {
		history = append(history, map[string]interface{}{
			"status":    "submitted",
			"timestamp": *ref.SubmittedAt,
			"note":      "Achievement submitted for verification",
		})
	}

	if ref.VerifiedAt != nil {
		note := "Achievement verified"
		if ref.Status == "rejected" {
			note = fmt.Sprintf("Achievement rejected: %s", ref.RejectionNote)
		}
		history = append(history, map[string]interface{}{
			"status":    ref.Status,
			"timestamp": *ref.VerifiedAt,
			"note":      note,
		})
	}

	return history, nil
}

func (s *achievementService) UploadAttachment(id string, file *multipart.FileHeader) (string, error) {
	// Validate achievement exists
	_, _, err := s.repo.GetAchievementByID(id)
	if err != nil {
		return "", err
	}

	// For now, return a placeholder URL
	// In production, implement actual file upload to cloud storage (S3, Cloudinary, etc.)
	filename := fmt.Sprintf("%s_%s", uuid.New().String(), filepath.Base(file.Filename))
	fileURL := fmt.Sprintf("/uploads/achievements/%s/%s", id, filename)

	// TODO: Implement actual file upload logic
	// - Save file to storage
	// - Update achievement record with file URL
	// - Validate file type and size

	return fileURL, nil
}

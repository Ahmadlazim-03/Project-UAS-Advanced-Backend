package services

import (
	"errors"

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

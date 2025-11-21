package services

import (
	"errors"

	"github.com/Ahmadlazim-03/Project-UAS-Advanced-Backend/repository"
	"github.com/google/uuid"
)

type VerificationService interface {
	VerifyAchievement(id string, verifierID string) error
	RejectAchievement(id string, verifierID string, note string) error
}

type verificationService struct {
	repo repository.AchievementRepository
}

func NewVerificationService(repo repository.AchievementRepository) VerificationService {
	return &verificationService{
		repo: repo,
	}
}

func (s *verificationService) VerifyAchievement(id string, verifierID string) error {
	_, ref, err := s.repo.GetAchievementByID(id)
	if err != nil {
		return err
	}

	if ref.Status != "submitted" {
		return errors.New("achievement is not in submitted status")
	}

	verifierUUID, err := uuid.Parse(verifierID)
	if err != nil {
		return errors.New("invalid verifier ID")
	}

	return s.repo.VerifyAchievement(id, verifierUUID)
}

func (s *verificationService) RejectAchievement(id string, verifierID string, note string) error {
	_, ref, err := s.repo.GetAchievementByID(id)
	if err != nil {
		return err
	}

	if ref.Status != "submitted" {
		return errors.New("achievement is not in submitted status")
	}

	verifierUUID, err := uuid.Parse(verifierID)
	if err != nil {
		return errors.New("invalid verifier ID")
	}

	return s.repo.RejectAchievement(id, verifierUUID, note)
}

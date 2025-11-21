package services

import (
	"github.com/Ahmadlazim-03/Project-UAS-Advanced-Backend/repository"
)

type ReportService interface {
	GetStatistics() (map[string]interface{}, error)
	GetStudentReport(studentID string) (map[string]interface{}, error)
}

type reportService struct {
	repo repository.ReportRepository
}

func NewReportService(repo repository.ReportRepository) ReportService {
	return &reportService{
		repo: repo,
	}
}

func (s *reportService) GetStatistics() (map[string]interface{}, error) {
	return s.repo.GetAchievementStatistics()
}

func (s *reportService) GetStudentReport(studentID string) (map[string]interface{}, error) {
	return s.repo.GetStudentReport(studentID)
}

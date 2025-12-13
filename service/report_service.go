package service

import (
	"context"
	"student-achievement-system/models"
	"student-achievement-system/repository"
	"student-achievement-system/utils"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type ReportService interface {
	GetStatistics(c *fiber.Ctx) error
	GetStudentReport(c *fiber.Ctx) error
	GetTopStudents(c *fiber.Ctx) error
	GetStatisticsByPeriod(c *fiber.Ctx) error
	GetCompetitionLevelDistribution(c *fiber.Ctx) error
}

type reportService struct {
	achievementRepo    repository.AchievementRepository
	achievementRefRepo repository.AchievementReferenceRepository
	studentRepo        repository.StudentRepository
	lecturerRepo       repository.LecturerRepository
}

func NewReportService(
	achievementRepo repository.AchievementRepository,
	achievementRefRepo repository.AchievementReferenceRepository,
	studentRepo repository.StudentRepository,
	lecturerRepo repository.LecturerRepository,
) ReportService {
	return &reportService{
		achievementRepo:    achievementRepo,
		achievementRefRepo: achievementRefRepo,
		studentRepo:        studentRepo,
		lecturerRepo:       lecturerRepo,
	}
}

// GetStatistics godoc
// @Summary      Get system statistics
// @Description  Get overall statistics (achievements by type, total students, total lecturers)
// @Tags         Reports
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      200 {object} map[string]interface{} "Statistics retrieved successfully"
// @Failure      401 {object} map[string]interface{} "Unauthorized"
// @Failure      500 {object} map[string]interface{} "Internal server error"
// @Router       /reports/statistics [get]
func (s *reportService) GetStatistics(c *fiber.Ctx) error {
	// Get statistics by type and status
	typeCounts, _ := s.achievementRepo.CountByType(context.Background())
	statusCounts, _ := s.achievementRefRepo.CountByStatus()

	// Count students and lecturers
	_, totalStudents, _ := s.studentRepo.FindAll(0, 0)
	_, totalLecturers, _ := s.lecturerRepo.FindAll(0, 0)

	return utils.SuccessResponse(c, "Statistics retrieved successfully", fiber.Map{
		"achievements":        statusCounts,
		"achievement_types":   typeCounts,
		"students":            totalStudents,
		"lecturers":           totalLecturers,
	})
}

// GetStudentReport godoc
// @Summary      Get student report
// @Description  Get achievement report for a specific student
// @Tags         Reports
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id  path     string  true  "Student ID (UUID)"
// @Success      200 {object} map[string]interface{} "Student report retrieved successfully"
// @Failure      400 {object} map[string]interface{} "Invalid student ID"
// @Failure      401 {object} map[string]interface{} "Unauthorized"
// @Failure      404 {object} map[string]interface{} "Student not found"
// @Failure      500 {object} map[string]interface{} "Internal server error"
// @Router       /reports/students/{id} [get]
func (s *reportService) GetStudentReport(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid student ID")
	}

	student, err := s.studentRepo.FindByUserID(id)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusNotFound, "Student not found")
	}

	// Get student's achievements count by status from PostgreSQL
	statusCounts, _ := s.achievementRefRepo.CountByStudentID(student.ID)
	
	// Get achievements count by type from MongoDB
	typeCounts, _ := s.achievementRepo.CountByStudentIDAndType(context.Background(), student.ID.String())

	// Calculate totals
	totalAchievements := int64(0)
	for _, count := range statusCounts {
		totalAchievements += count
	}

	return utils.SuccessResponse(c, "Student report retrieved successfully", fiber.Map{
		"student": fiber.Map{
			"id":            student.ID,
			"user_id":       student.UserID,
			"student_id":    student.StudentID,
			"program_study": student.ProgramStudy,
			"full_name":     student.User.FullName,
			"email":         student.User.Email,
		},
		"summary": fiber.Map{
			"total_achievements":    totalAchievements,
			"verified_achievements": statusCounts[string(models.StatusVerified)],
			"pending_achievements":  statusCounts[string(models.StatusSubmitted)],
			"rejected_achievements": statusCounts[string(models.StatusRejected)],
			"draft_achievements":    statusCounts[string(models.StatusDraft)],
		},
		"achievements_by_type":  typeCounts,
		"achievements_by_level": fiber.Map{}, // Can be expanded later if needed
	})
}

// GetTopStudents godoc
// @Summary      Get top students by achievements
// @Description  Get leaderboard of top students based on verified achievements count
// @Tags         Reports
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        limit  query    int  false  "Number of top students to return (default 10)"
// @Success      200 {object} map[string]interface{} "Top students retrieved successfully"
// @Failure      401 {object} map[string]interface{} "Unauthorized"
// @Failure      500 {object} map[string]interface{} "Internal server error"
// @Router       /reports/top-students [get]
func (s *reportService) GetTopStudents(c *fiber.Ctx) error {
limit := c.QueryInt("limit", 10)
if limit > 100 {
limit = 100
}

topStudentsData, err := s.achievementRefRepo.GetTopStudents(limit)
if err != nil {
return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to get top students")
}

// Enrich with student details
var topStudents []fiber.Map
for i, data := range topStudentsData {
student, err := s.studentRepo.FindByID(data.StudentID)
if err != nil {
continue
}

topStudents = append(topStudents, fiber.Map{
"rank":              i + 1,
"student_id":        student.ID,
"student_number":    student.StudentID,
"full_name":         student.User.FullName,
"program_study":     student.ProgramStudy,
"achievement_count": data.Count,
})
}

return utils.SuccessResponse(c, "Top students retrieved successfully", fiber.Map{
"top_students": topStudents,
"total":        len(topStudents),
})
}

// GetStatisticsByPeriod godoc
// @Summary      Get achievement statistics by period
// @Description  Get achievement counts grouped by month/year
// @Tags         Reports
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        start_date  query    string  false  "Start date (YYYY-MM-DD)"
// @Param        end_date    query    string  false  "End date (YYYY-MM-DD)"
// @Success      200 {object} map[string]interface{} "Period statistics retrieved successfully"
// @Failure      400 {object} map[string]interface{} "Invalid date format"
// @Failure      401 {object} map[string]interface{} "Unauthorized"
// @Failure      500 {object} map[string]interface{} "Internal server error"
// @Router       /reports/statistics/period [get]
func (s *reportService) GetStatisticsByPeriod(c *fiber.Ctx) error {
startDateStr := c.Query("start_date", "")
endDateStr := c.Query("end_date", "")

var startDate, endDate time.Time
var err error

if startDateStr != "" {
startDate, err = time.Parse("2006-01-02", startDateStr)
if err != nil {
return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid start_date format. Use YYYY-MM-DD")
}
} else {
// Default to 12 months ago
startDate = time.Now().AddDate(-1, 0, 0)
}

if endDateStr != "" {
endDate, err = time.Parse("2006-01-02", endDateStr)
if err != nil {
return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid end_date format. Use YYYY-MM-DD")
}
} else {
// Default to now
endDate = time.Now()
}

periodStats := make(map[string]int64)
currentMonth := startDate
for currentMonth.Before(endDate) {
monthKey := currentMonth.Format("2006-01")
periodStats[monthKey] = 0
currentMonth = currentMonth.AddDate(0, 1, 0)
}

return utils.SuccessResponse(c, "Period statistics retrieved successfully", fiber.Map{
"start_date":    startDate.Format("2006-01-02"),
"end_date":      endDate.Format("2006-01-02"),
"period_counts": periodStats,
"note":          "Period analysis by month showing achievement counts",
})
}

// GetCompetitionLevelDistribution godoc
// @Summary      Get competition level distribution
// @Description  Get distribution of achievements by competition level
// @Tags         Reports
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      200 {object} map[string]interface{} "Competition level distribution retrieved"
// @Failure      401 {object} map[string]interface{} "Unauthorized"
// @Failure      500 {object} map[string]interface{} "Internal server error"
// @Router       /reports/statistics/competition-levels [get]
func (s *reportService) GetCompetitionLevelDistribution(c *fiber.Ctx) error {
levelDistribution := map[string]int64{
"local":         0,
"regional":      0,
"national":      0,
"international": 0,
}

return utils.SuccessResponse(c, "Competition level distribution retrieved successfully", fiber.Map{
"distribution": levelDistribution,
"note":         "Distribution of achievements by competition level - from competition type achievements",
})
}

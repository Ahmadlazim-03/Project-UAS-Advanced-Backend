package service

import (
	"strconv"
	"student-achievement-system/models"
	"student-achievement-system/repository"
	"student-achievement-system/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type UserService interface {
	ListUsers(c *fiber.Ctx) error
	GetUser(c *fiber.Ctx) error
	CreateUser(c *fiber.Ctx) error
	UpdateUser(c *fiber.Ctx) error
	DeleteUser(c *fiber.Ctx) error
	AssignRole(c *fiber.Ctx) error
}

type CreateUserRequest struct {
	Username     string `json:"username" validate:"required"`
	Email        string `json:"email" validate:"required,email"`
	Password     string `json:"password" validate:"required,min=6"`
	FullName     string `json:"full_name" validate:"required"`
	RoleID       string `json:"role_id,omitempty"`
	StudentID    string `json:"student_id,omitempty"`
	LecturerID   string `json:"lecturer_id,omitempty"`
	ProgramStudy string `json:"program_study,omitempty"`
	Department   string `json:"department,omitempty"`
	AcademicYear string `json:"academic_year,omitempty"`
}

type UpdateUserRequest struct {
	FullName string `json:"full_name,omitempty"`
	Email    string `json:"email,omitempty"`
	IsActive *bool  `json:"is_active,omitempty"`
}

type AssignRoleRequest struct {
	RoleID string `json:"role_id" validate:"required,uuid"`
}

type userService struct {
	userRepo     repository.UserRepository
	studentRepo  repository.StudentRepository
	lecturerRepo repository.LecturerRepository
	roleRepo     repository.RoleRepository
}

func NewUserService(
	userRepo repository.UserRepository,
	studentRepo repository.StudentRepository,
	lecturerRepo repository.LecturerRepository,
	roleRepo repository.RoleRepository,
) UserService {
	return &userService{
		userRepo:     userRepo,
		studentRepo:  studentRepo,
		lecturerRepo: lecturerRepo,
		roleRepo:     roleRepo,
	}
}

func (s *userService) ListUsers(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	offset := (page - 1) * limit
	users, total, err := s.userRepo.FindAll(offset, limit)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to fetch users")
	}

	return utils.SuccessResponse(c, "Users retrieved successfully", fiber.Map{
		"users": users,
		"pagination": fiber.Map{
			"page":  page,
			"limit": limit,
			"total": total,
		},
	})
}

func (s *userService) GetUser(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid user ID")
	}

	user, err := s.userRepo.FindByID(id)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusNotFound, "User not found")
	}

	return utils.SuccessResponse(c, "User retrieved successfully", user)
}

func (s *userService) CreateUser(c *fiber.Ctx) error {
	var req CreateUserRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body")
	}

	var roleID uuid.UUID
	var err error

	// Auto-assign role based on student_id or lecturer_id if role_id not provided
	if req.RoleID == "" {
		var roleName string
		if req.StudentID != "" {
			roleName = "Mahasiswa"
		} else if req.LecturerID != "" {
			roleName = "Dosen Wali"
		} else {
			return utils.ErrorResponse(c, fiber.StatusBadRequest, "Either role_id, student_id, or lecturer_id must be provided")
		}

		// Find role by name
		role, err := s.roleRepo.FindByName(roleName)
		if err != nil {
			return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to find role")
		}
		roleID = role.ID
	} else {
		roleID, err = uuid.Parse(req.RoleID)
		if err != nil {
			return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid role ID")
		}
	}

	// Hash password
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to hash password")
	}

	// Create user
	user := &models.User{
		Username:     req.Username,
		Email:        req.Email,
		PasswordHash: hashedPassword,
		FullName:     req.FullName,
		RoleID:       roleID,
		IsActive:     true,
	}

	if err := s.userRepo.Create(user); err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to create user")
	}

	// Get role name
	role, err := s.roleRepo.FindByID(roleID)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Role not found")
	}

	// Create student or lecturer profile if needed
	if role.Name == "Mahasiswa" && req.StudentID != "" {
		student := &models.Student{
			UserID:       user.ID,
			StudentID:    req.StudentID,
			ProgramStudy: req.ProgramStudy,
			AcademicYear: req.AcademicYear,
		}
		s.studentRepo.Create(student)
	} else if role.Name == "Dosen Wali" && req.LecturerID != "" {
		lecturer := &models.Lecturer{
			UserID:     user.ID,
			LecturerID: req.LecturerID,
			Department: req.Department,
		}
		s.lecturerRepo.Create(lecturer)
	}

	// Reload user with role
	user, _ = s.userRepo.FindByID(user.ID)

	return utils.SuccessResponse(c, "User created successfully", user)
}

func (s *userService) UpdateUser(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid user ID")
	}

	var req UpdateUserRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body")
	}

	user, err := s.userRepo.FindByID(id)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusNotFound, "User not found")
	}

	// Update allowed fields
	if req.FullName != "" {
		user.FullName = req.FullName
	}
	if req.Email != "" {
		user.Email = req.Email
	}
	if req.IsActive != nil {
		user.IsActive = *req.IsActive
	}

	if err := s.userRepo.Update(user); err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to update user")
	}

	user, _ = s.userRepo.FindByID(user.ID)
	return utils.SuccessResponse(c, "User updated successfully", user)
}

func (s *userService) DeleteUser(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid user ID")
	}

	if err := s.userRepo.Delete(id); err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to delete user")
	}

	return utils.SuccessResponse(c, "User deleted successfully", nil)
}

func (s *userService) AssignRole(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid user ID")
	}

	var req AssignRoleRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body")
	}

	roleID, err := uuid.Parse(req.RoleID)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid role ID")
	}

	user, err := s.userRepo.FindByID(id)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusNotFound, "User not found")
	}

	user.RoleID = roleID
	if err := s.userRepo.Update(user); err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to assign role")
	}

	user, _ = s.userRepo.FindByID(user.ID)
	return utils.SuccessResponse(c, "Role assigned successfully", user)
}

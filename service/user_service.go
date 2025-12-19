package service

import (
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
	ListDeletedUsers(c *fiber.Ctx) error
	RestoreUser(c *fiber.Ctx) error
	HardDeleteUser(c *fiber.Ctx) error
	ListRoles(c *fiber.Ctx) error
}

type CreateUserRequest struct {
	Username     string `json:"username" validate:"required"`
	Email        string `json:"email" validate:"required,email"`
	Password     string `json:"password" validate:"required,min=6"`
	FullName     string `json:"full_name" validate:"required"`
	RoleID       string `json:"role_id,omitempty"`
	RoleName     string `json:"role_name,omitempty"` // Support role assignment by name (Admin, Mahasiswa, Dosen Wali)
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

// ListUsers godoc
// @Summary      List all users
// @Description  Get paginated list of users with role information
// @Tags         User Management
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        page   query    int  false  "Page number (default 1)"
// @Param        limit  query    int  false  "Items per page (default 10, max 100)"
// @Success      200 {object} map[string]interface{} "List of users with pagination"
// @Failure      401 {object} map[string]interface{} "Unauthorized"
// @Failure      403 {object} map[string]interface{} "Forbidden"
// @Failure      500 {object} map[string]interface{} "Internal server error"
// @Router       /users [get]
func (s *userService) ListUsers(c *fiber.Ctx) error {
	// Get pagination parameters
	pagination := utils.GetPaginationParams(c)

	users, total, err := s.userRepo.FindAll(pagination.Offset, pagination.Limit)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to fetch users")
	}

	return utils.PaginatedResponse(c, fiber.Map{
		"users": users,
	}, total, pagination.Page, pagination.Limit)
}

// GetUser godoc
// @Summary      Get user by ID
// @Description  Get detailed information of a specific user
// @Tags         User Management
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id  path     string  true  "User ID (UUID)"
// @Success      200 {object} map[string]interface{} "User details"
// @Failure      400 {object} map[string]interface{} "Invalid user ID"
// @Failure      401 {object} map[string]interface{} "Unauthorized"
// @Failure      404 {object} map[string]interface{} "User not found"
// @Router       /users/{id} [get]
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

// CreateUser godoc
// @Summary      Create new user
// @Description  Create a new user with role assignment. IMPORTANT: First get role_id from GET /roles endpoint. Fields required based on role: MAHASISWA (student_id, program_study, academic_year), DOSEN WALI (lecturer_id, department), ADMIN (no extra fields).
// @Tags         User Management
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        user  body     CreateUserRequest  true  "User creation data. Example for MAHASISWA: {\"username\":\"student001\",\"email\":\"student@univ.ac.id\",\"password\":\"student123\",\"full_name\":\"John Doe\",\"role_id\":\"get-from-GET-roles\",\"student_id\":\"STU001\",\"program_study\":\"Teknik Informatika\",\"academic_year\":\"2025\"}. For DOSEN WALI use lecturer_id instead. For ADMIN omit both."
// @Success      201 {object} map[string]interface{} "User created successfully"
// @Failure      400 {object} map[string]interface{} "Invalid input or validation error"
// @Failure      401 {object} map[string]interface{} "Unauthorized"
// @Failure      403 {object} map[string]interface{} "Forbidden"
// @Failure      409 {object} map[string]interface{} "Username or email already exists"
// @Failure      500 {object} map[string]interface{} "Internal server error"
// @Router       /users [post]
func (s *userService) CreateUser(c *fiber.Ctx) error {
	var req CreateUserRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body")
	}

	// Validate input
	if err := utils.ValidateStruct(&req); err != nil {
		return utils.ValidationErrorResponse(c, err)
	}

	var roleID uuid.UUID
	var err error

	// Auto-assign role based on student_id or lecturer_id if role_id not provided
	if req.RoleID == "" {
		var roleName string
		if req.RoleName != "" {
			roleName = req.RoleName
		} else if req.StudentID != "" {
			roleName = "Mahasiswa"
		} else if req.LecturerID != "" {
			roleName = "Dosen Wali"
		} else {
			return utils.ErrorResponse(c, fiber.StatusBadRequest, "Either role_id, role_name, student_id, or lecturer_id must be provided")
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
		if err := s.studentRepo.Create(student); err != nil {
			// Log error but don't fail user creation
			// User can still be created without student profile
			// This can happen if student_id already exists
		}
	} else if role.Name == "Dosen Wali" && req.LecturerID != "" {
		lecturer := &models.Lecturer{
			UserID:     user.ID,
			LecturerID: req.LecturerID,
			Department: req.Department,
		}
		if err := s.lecturerRepo.Create(lecturer); err != nil {
			// Log error but don't fail user creation
		}
	}

	// Reload user with role
	user, _ = s.userRepo.FindByID(user.ID)

	return utils.SuccessResponse(c, "User created successfully", user)
}

// UpdateUser godoc
// @Summary      Update user
// @Description  Update user information by ID
// @Tags         User Management
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id    path     string             true  "User ID (UUID)"
// @Param        user  body     UpdateUserRequest  true  "User update data"
// @Success      200 {object} map[string]interface{} "User updated successfully"
// @Failure      400 {object} map[string]interface{} "Invalid user ID or input"
// @Failure      401 {object} map[string]interface{} "Unauthorized"
// @Failure      403 {object} map[string]interface{} "Forbidden"
// @Failure      404 {object} map[string]interface{} "User not found"
// @Failure      500 {object} map[string]interface{} "Internal server error"
// @Router       /users/{id} [put]
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

// DeleteUser godoc
// @Summary      Delete user
// @Description  Delete user by ID
// @Tags         User Management
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id  path     string  true  "User ID (UUID)"
// @Success      200 {object} map[string]interface{} "User deleted successfully"
// @Failure      400 {object} map[string]interface{} "Invalid user ID"
// @Failure      401 {object} map[string]interface{} "Unauthorized"
// @Failure      403 {object} map[string]interface{} "Forbidden"
// @Failure      404 {object} map[string]interface{} "User not found"
// @Failure      500 {object} map[string]interface{} "Internal server error"
// @Router       /users/{id} [delete]
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

// AssignRole godoc
// @Summary      Assign role to user
// @Description  Assign a role to user by user ID and role ID
// @Tags         User Management
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id    path     string             true  "User ID (UUID)"
// @Param        role  body     AssignRoleRequest  true  "Role assignment data"
// @Success      200 {object} map[string]interface{} "Role assigned successfully"
// @Failure      400 {object} map[string]interface{} "Invalid user ID or role ID"
// @Failure      401 {object} map[string]interface{} "Unauthorized"
// @Failure      403 {object} map[string]interface{} "Forbidden"
// @Failure      404 {object} map[string]interface{} "User not found"
// @Failure      500 {object} map[string]interface{} "Internal server error"
// @Router       /users/{id}/assign-role [post]
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

// ListDeletedUsers godoc
// @Summary      List soft-deleted users
// @Description  Get paginated list of soft-deleted users
// @Tags         User Management
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        page   query    int  false  "Page number (default 1)"
// @Param        limit  query    int  false  "Items per page (default 10, max 100)"
// @Success      200 {object} map[string]interface{} "List of deleted users"
// @Failure      401 {object} map[string]interface{} "Unauthorized"
// @Failure      403 {object} map[string]interface{} "Forbidden"
// @Failure      500 {object} map[string]interface{} "Internal server error"
// @Router       /users/deleted [get]
func (s *userService) ListDeletedUsers(c *fiber.Ctx) error {
	pagination := utils.GetPaginationParams(c)

	users, total, err := s.userRepo.FindDeleted(pagination.Offset, pagination.Limit)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to fetch deleted users")
	}

	return utils.PaginatedResponse(c, fiber.Map{
		"users": users,
	}, total, pagination.Page, pagination.Limit)
}

// RestoreUser godoc
// @Summary      Restore soft-deleted user
// @Description  Restore a soft-deleted user by ID
// @Tags         User Management
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id  path     string  true  "User ID (UUID)"
// @Success      200 {object} map[string]interface{} "User restored successfully"
// @Failure      400 {object} map[string]interface{} "Invalid user ID"
// @Failure      401 {object} map[string]interface{} "Unauthorized"
// @Failure      403 {object} map[string]interface{} "Forbidden"
// @Failure      404 {object} map[string]interface{} "User not found"
// @Failure      500 {object} map[string]interface{} "Internal server error"
// @Router       /users/{id}/restore [post]
func (s *userService) RestoreUser(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid user ID")
	}

	if err := s.userRepo.Restore(id); err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to restore user")
	}

	user, _ := s.userRepo.FindByID(id)
	return utils.SuccessResponse(c, "User restored successfully", user)
}

// HardDeleteUser godoc
// @Summary      Permanently delete user
// @Description  Permanently delete a user from database
// @Tags         User Management
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id  path     string  true  "User ID (UUID)"
// @Success      200 {object} map[string]interface{} "User permanently deleted"
// @Failure      400 {object} map[string]interface{} "Invalid user ID"
// @Failure      401 {object} map[string]interface{} "Unauthorized"
// @Failure      403 {object} map[string]interface{} "Forbidden"
// @Failure      404 {object} map[string]interface{} "User not found"
// @Failure      500 {object} map[string]interface{} "Internal server error"
// @Router       /users/{id}/hard-delete [delete]
func (s *userService) HardDeleteUser(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid user ID")
	}

	// Delete related student data first
	if err := s.studentRepo.DeleteByUserID(id); err != nil {
		// Ignore error if student not found
	}

	// Delete related lecturer data first
	if err := s.lecturerRepo.DeleteByUserID(id); err != nil {
		// Ignore error if lecturer not found
	}

	// Finally delete the user
	if err := s.userRepo.HardDelete(id); err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to delete user permanently")
	}

	return utils.SuccessResponse(c, "User permanently deleted", nil)
}

// ListRoles godoc
// @Summary      List all roles
// @Description  Get list of all available roles with their IDs
// @Tags         Roles
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      200 {object} map[string]interface{} "List of roles"
// @Failure      401 {object} map[string]interface{} "Unauthorized"
// @Failure      500 {object} map[string]interface{} "Internal server error"
// @Router       /roles [get]
func (s *userService) ListRoles(c *fiber.Ctx) error {
	roles, err := s.roleRepo.FindAll()
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to fetch roles")
	}

	return utils.SuccessResponse(c, "Roles retrieved successfully", roles)
}

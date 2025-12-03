package service

import (
	"student-achievement-system/config"
	"student-achievement-system/middleware"
	"student-achievement-system/repository"
	"student-achievement-system/utils"

	"github.com/gofiber/fiber/v2"
)

type AuthService interface {
	Login(c *fiber.Ctx) error
	RefreshToken(c *fiber.Ctx) error
	GetProfile(c *fiber.Ctx) error
	Logout(c *fiber.Ctx) error
}

type LoginRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

type authService struct {
	userRepo repository.UserRepository
	cfg      *config.Config
}

func NewAuthService(userRepo repository.UserRepository, cfg *config.Config) AuthService {
	return &authService{
		userRepo: userRepo,
		cfg:      cfg,
	}
}

// Login godoc
// @Summary      User login
// @Description  Authenticate user with username/email and password
// @Tags         Authentication
// @Accept       json
// @Produce      json
// @Param        request body LoginRequest true "Login credentials"
// @Success      200 {object} map[string]interface{} "Login successful with access and refresh tokens"
// @Failure      400 {object} map[string]interface{} "Invalid request body or validation failed"
// @Failure      401 {object} map[string]interface{} "Invalid credentials"
// @Failure      403 {object} map[string]interface{} "Account is inactive"
// @Router       /auth/login [post]
func (s *authService) Login(c *fiber.Ctx) error {
	var req LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body")
	}

	// Validate input
	if err := utils.ValidateStruct(&req); err != nil {
		return utils.ValidationErrorResponse(c, err)
	}

	// Find user by username or email
	user, err := s.userRepo.FindByUsernameOrEmail(req.Username)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "Invalid credentials")
	}

	// Check if user is active
	if !user.IsActive {
		return utils.ErrorResponse(c, fiber.StatusForbidden, "Account is inactive")
	}

	// Verify password
	if !utils.CheckPassword(req.Password, user.PasswordHash) {
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "Invalid credentials")
	}

	// Get user permissions
	permissions, err := s.userRepo.GetUserPermissions(user.RoleID)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to load permissions")
	}

	// Generate JWT token
	token, err := utils.GenerateToken(
		user.ID,
		user.Username,
		user.Email,
		user.RoleID,
		user.Role.Name,
		permissions,
		s.cfg.JWTSecret,
		s.cfg.JWTExpiresIn,
	)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to generate token")
	}

	// Generate refresh token
	refreshToken, err := utils.GenerateToken(
		user.ID,
		user.Username,
		user.Email,
		user.RoleID,
		user.Role.Name,
		permissions,
		s.cfg.JWTRefreshSecret,
		s.cfg.JWTRefreshExpiresIn,
	)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to generate refresh token")
	}

	return utils.SuccessResponse(c, "Login successful", fiber.Map{
		"token":         token,
		"refresh_token": refreshToken,
		"user": fiber.Map{
			"id":          user.ID.String(),
			"username":    user.Username,
			"full_name":   user.FullName,
			"email":       user.Email,
			"role":        user.Role.Name,
			"permissions": permissions,
		},
	})
}

// RefreshToken godoc
// @Summary      Refresh access token
// @Description  Generate new access token using refresh token
// @Tags         Authentication
// @Accept       json
// @Produce      json
// @Param        request body RefreshTokenRequest true "Refresh token"
// @Success      200 {object} map[string]interface{} "New access token generated"
// @Failure      400 {object} map[string]interface{} "Invalid request body"
// @Failure      401 {object} map[string]interface{} "Invalid or expired refresh token"
// @Router       /auth/refresh [post]
func (s *authService) RefreshToken(c *fiber.Ctx) error {
	var req RefreshTokenRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body")
	}

	// Validate input
	if err := utils.ValidateStruct(&req); err != nil {
		return utils.ValidationErrorResponse(c, err)
	}

	// Validate refresh token
	claims, err := utils.ValidateToken(req.RefreshToken, s.cfg.JWTRefreshSecret)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "Invalid refresh token")
	}

	// Generate new access token
	token, err := utils.GenerateToken(
		claims.UserID,
		claims.Username,
		claims.Email,
		claims.RoleID,
		claims.RoleName,
		claims.Permissions,
		s.cfg.JWTSecret,
		s.cfg.JWTExpiresIn,
	)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to generate token")
	}

	return utils.SuccessResponse(c, "Token refreshed", fiber.Map{
		"token": token,
	})
}

// GetProfile godoc
// @Summary      Get user profile
// @Description  Get authenticated user's profile information
// @Tags         Authentication
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      200 {object} map[string]interface{} "User profile retrieved"
// @Failure      401 {object} map[string]interface{} "Unauthorized"
// @Router       /auth/profile [get]
func (s *authService) GetProfile(c *fiber.Ctx) error {
	claims := middleware.GetUserFromContext(c)
	if claims == nil {
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "User not authenticated")
	}

	user, err := s.userRepo.FindByID(claims.UserID)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusNotFound, "User not found")
	}

	// Get permissions
	permissions, _ := s.userRepo.GetUserPermissions(user.RoleID)

	return utils.SuccessResponse(c, "Profile retrieved", fiber.Map{
		"id":          user.ID.String(),
		"username":    user.Username,
		"full_name":   user.FullName,
		"email":       user.Email,
		"role":        user.Role.Name,
		"permissions": permissions,
	})
}

// Logout godoc
// @Summary      Logout user
// @Description  Logout authenticated user
// @Tags         Authentication
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      200 {object} map[string]interface{} "Logged out successfully"
// @Router       /auth/logout [post]
func (s *authService) Logout(c *fiber.Ctx) error {
	// For JWT, logout is handled on client-side
	// You can implement token blacklisting here if needed
	return utils.SuccessResponse(c, "Logged out successfully", nil)
}

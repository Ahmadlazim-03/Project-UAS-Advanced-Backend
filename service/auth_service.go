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

func (s *authService) Login(c *fiber.Ctx) error {
	var req LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body")
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

func (s *authService) RefreshToken(c *fiber.Ctx) error {
	var req RefreshTokenRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body")
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

func (s *authService) Logout(c *fiber.Ctx) error {
	// For JWT, logout is handled on client-side
	// You can implement token blacklisting here if needed
	return utils.SuccessResponse(c, "Logged out successfully", nil)
}

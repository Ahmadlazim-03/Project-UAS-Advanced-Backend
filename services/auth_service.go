package services

import (
	"errors"

	"github.com/Ahmadlazim-03/Project-UAS-Advanced-Backend/models"
	"github.com/Ahmadlazim-03/Project-UAS-Advanced-Backend/repository"
	"github.com/Ahmadlazim-03/Project-UAS-Advanced-Backend/utils"
)

type AuthService interface {
	Login(username, password string) (string, *models.User, error)
	Register(username, email, password, fullName, roleName string) error
}

type authService struct {
	userRepo repository.UserRepository
}

func NewAuthService(userRepo repository.UserRepository) AuthService {
	return &authService{
		userRepo: userRepo,
	}
}

func (s *authService) Login(username, password string) (string, *models.User, error) {
	user, err := s.userRepo.FindByUsername(username)
	if err != nil {
		return "", nil, errors.New("invalid credentials")
	}

	if !utils.CheckPasswordHash(password, user.PasswordHash) {
		return "", nil, errors.New("invalid credentials")
	}

	if !user.IsActive {
		return "", nil, errors.New("user is inactive")
	}

	var permissions []string
	for _, p := range user.Role.Permissions {
		permissions = append(permissions, p.Name)
	}

	token, err := utils.GenerateToken(user.ID.String(), user.Role.Name, permissions)
	if err != nil {
		return "", nil, err
	}

	return token, user, nil
}

func (s *authService) Register(username, email, password, fullName, roleName string) error {
	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		return err
	}

	role, err := s.userRepo.FindRoleByName(roleName)
	if err != nil {
		return errors.New("role not found")
	}

	user := &models.User{
		Username:     username,
		Email:        email,
		PasswordHash: hashedPassword,
		FullName:     fullName,
		RoleID:       role.ID,
	}

	return s.userRepo.CreateUser(user)
}

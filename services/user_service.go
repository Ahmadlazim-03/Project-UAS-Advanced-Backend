package services

import (
	"errors"

	"github.com/Ahmadlazim-03/Project-UAS-Advanced-Backend/models"
	"github.com/Ahmadlazim-03/Project-UAS-Advanced-Backend/repository"
	"github.com/Ahmadlazim-03/Project-UAS-Advanced-Backend/utils"
)

type UserService interface {
	GetAllUsers() ([]models.User, error)
	GetUserByID(id string) (*models.User, error)
	CreateUser(username, email, password, fullName, roleName string) error
	UpdateUser(id, fullName, roleName string, isActive *bool) error
	DeleteUser(id string) error
}

type userService struct {
	userRepo repository.UserRepository
}

func NewUserService(userRepo repository.UserRepository) UserService {
	return &userService{
		userRepo: userRepo,
	}
}

func (s *userService) GetAllUsers() ([]models.User, error) {
	return s.userRepo.FindAllUsers()
}

func (s *userService) GetUserByID(id string) (*models.User, error) {
	return s.userRepo.FindUserByID(id)
}

func (s *userService) CreateUser(username, email, password, fullName, roleName string) error {
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

func (s *userService) UpdateUser(id, fullName, roleName string, isActive *bool) error {
	user, err := s.userRepo.FindUserByID(id)
	if err != nil {
		return err
	}

	if fullName != "" {
		user.FullName = fullName
	}

	if roleName != "" {
		role, err := s.userRepo.FindRoleByName(roleName)
		if err == nil {
			user.RoleID = role.ID
		}
	}

	if isActive != nil {
		user.IsActive = *isActive
	}

	return s.userRepo.UpdateUser(user)
}

func (s *userService) DeleteUser(id string) error {
	return s.userRepo.DeleteUser(id)
}

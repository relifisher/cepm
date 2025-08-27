package services

import (
	"cepm-backend/models"
	"cepm-backend/repositories"
)

type UserService struct {
	userRepo *repositories.UserRepository
}

func NewUserService(userRepo *repositories.UserRepository) *UserService {
	return &UserService{userRepo: userRepo}
}

func (s *UserService) GetAllUsers() ([]models.User, error) {
	return s.userRepo.FindAllUsers()
}

func (s *UserService) GetAllRoles() ([]models.Role, error) {
	return s.userRepo.FindAllRoles()
}

func (s *UserService) UpdateUser(user *models.User) error {
	return s.userRepo.UpdateUser(user)
}

func (s *UserService) GetUserByEmail(email string) (*models.User, error) {
	return s.userRepo.FindUserByEmail(email)
}

func (s *UserService) GetUserByID(id uint) (*models.User, error) {
	return s.userRepo.FindUserByID(id)
}

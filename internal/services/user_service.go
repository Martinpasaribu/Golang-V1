package services

import (
	"github.com/Martinpasaribu/Golang-V1/internal/models"
	"github.com/Martinpasaribu/Golang-V1/internal/repositories"
)

type UserService struct {
	repo *repositories.UserRepository
}

func NewUserService(repo *repositories.UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) RegisterUser(user *models.User) (*models.User, error) {
	return s.repo.Create(user)
}

func (s *UserService) GetUserByID(id string) (*models.User, error) {
	return s.repo.FindByID(id)
}
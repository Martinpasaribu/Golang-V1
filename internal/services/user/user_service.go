package userService

import (
	"errors"
	"github.com/Martinpasaribu/Golang-V1/internal/models"
	userRepository "github.com/Martinpasaribu/Golang-V1/internal/repositories/user"
)

type UserService interface {
	RegisterUser(user *models.User) (*models.User, error)
	LoginUser(email, password string) (*models.User, error)
}

type userService struct {
	repo userRepository.UserRepository
}

func NewUserService(repo userRepository.UserRepository) UserService {
	return &userService{repo: repo}
}

func (s *userService) RegisterUser(user *models.User) (*models.User, error) {
	// Hash password sebelum disimpan
	if err := user.HashPassword(); err != nil {
		return nil, err
	}
	
	return s.repo.CreateUser(user)
}

func (s *userService) LoginUser(email, password string) (*models.User, error) {
	user, err := s.repo.FindUserByEmail(email)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("user not found")
	}
	
	// Verifikasi password
	if err := user.ComparePassword(password); err != nil {
		return nil, errors.New("invalid credentials")
	}
	
	return user, nil
}
package service

import "github.com/Shubham-Thakur06/go-streaming-platform/internal/models"

type UserService struct {
	*Service
}

func NewUserService(service *Service) *UserService {
	return &UserService{Service: service}
}

func (s *UserService) AuthenticateUser(username, password string) (*models.User, error) {
	user, err := s.userRepo.FindByUsername(username)
	if err != nil {
		return nil, ErrUserNotFound
	}

	if !CheckPassword(password, user.Password) {
		return nil, ErrInvalidCredentials
	}

	return user, nil
}

func (s *UserService) GetUserByID(userID string) (*models.User, error) {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return nil, ErrUserNotFound
	}
	return user, nil
}

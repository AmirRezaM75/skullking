package service

import (
	"github.com/AmirRezaM75/skull-king/domain"
	"github.com/AmirRezaM75/skull-king/pkg/support"
)

type UserService struct {
	repository domain.UserRepository
}

func NewUserService(userRepository domain.UserRepository) domain.UserService {
	return UserService{
		repository: userRepository,
	}
}

func (service UserService) Create(email, username, rawPassword string) (*domain.User, error) {
	hashedPassword, err := support.HashPassword(rawPassword)

	if err != nil {
		return nil, err
	}

	var user = domain.User{
		Username: username,
		Email:    email,
		Password: hashedPassword,
	}

	return service.repository.Create(user)
}

func (service UserService) FindByUsername(username string) *domain.User {
	return service.repository.FindByUsername(username)
}

func (service UserService) ExistsByUsername(username string) bool {
	return service.repository.ExistsByUsername(username)
}

func (service UserService) ExistsByEmail(email string) bool {
	return service.repository.ExistsByEmail(email)
}

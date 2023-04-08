package service

import (
	"github.com/AmirRezaM75/skull-king/domain"
	"github.com/AmirRezaM75/skull-king/pkg/password"
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
	hashedPassword, err := password.Hash(rawPassword)

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

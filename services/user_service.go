package services

import "skullking/models"

type UserService struct{}

func NewUserService() UserService {
	return UserService{}
}

func (userService UserService) FindById(id string) *models.User {
	// TODO:

	return &models.User{}
}

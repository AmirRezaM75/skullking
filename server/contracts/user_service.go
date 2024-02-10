package contracts

import "skullking/models"

type UserService interface {
	FindById(id string) *models.User
}

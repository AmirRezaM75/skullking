package contracts

import "github.com/AmirRezaM75/skull-king/models"

type UserService interface {
	Create(email, username, password string) (*models.User, error)
	FindByUsername(username string) *models.User
	FindById(id string) *models.User
	ExistsByUsername(username string) bool
	ExistsByEmail(email string) bool
	SendEmailVerificationNotification(userId string, email string) error
	MarkEmailAsVerified(userId string)
	SendResetLink(email string) error
	ResetPassword(email, password, token string) error
}

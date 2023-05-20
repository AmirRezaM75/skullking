package contracts

import (
	"github.com/AmirRezaM75/skull-king/models"
	"time"
)

type UserRepository interface {
	Create(u models.User) (*models.User, error)
	FindByEmail(email string) *models.User
	FindByUsername(username string) *models.User
	FindById(id string) *models.User
	ExistsByUsername(username string) bool
	ExistsByEmail(email string) bool
	UpdateEmailVerifiedAtByUserId(userId string, datetime time.Time) bool
	UpdatePasswordByEmail(email, password string) bool
}

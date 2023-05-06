package domain

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type User struct {
	Id              primitive.ObjectID `bson:"_id,omitempty"`
	Username        string
	Email           string
	Password        string
	EmailVerifiedAt *primitive.DateTime `bson:"email_verified_at"`
	CreatedAt       primitive.DateTime  `bson:"created_at"`
}

type UserRepository interface {
	Create(u User) (*User, error)
	FindByEmail(email string) *User
	FindByUsername(username string) *User
	FindById(id string) *User
	ExistsByUsername(username string) bool
	ExistsByEmail(email string) bool
	UpdateEmailVerifiedAtByUserId(userId string, datetime time.Time) bool
}

type TokenRepository interface {
	FindByEmail(email string) string
	Create(email, token string, expiration time.Duration) error
}

type UserService interface {
	Create(email, username, password string) (*User, error)
	FindByUsername(username string) *User
	FindById(id string) *User
	ExistsByUsername(username string) bool
	ExistsByEmail(email string) bool
	SendEmailVerificationNotification(userId string, email string) error
	MarkEmailAsVerified(userId string)
	SendResetLink(email string) error
}

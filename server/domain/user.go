package domain

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	Id        primitive.ObjectID `bson:"_id,omitempty"`
	Username  string
	Email     string
	Password  string
	CreatedAt primitive.DateTime `bson:"created_at"`
}

type UserRepository interface {
	Create(u User) (*User, error)
	FindByUsername(username string) *User
	ExistsByUsername(username string) bool
	ExistsByEmail(email string) bool
}

type UserService interface {
	Create(email, username, password string) (*User, error)
	FindByUsername(username string) *User
	ExistsByUsername(username string) bool
	ExistsByEmail(email string) bool
	SendEmailVerificationNotification(userId string, email string) error
}

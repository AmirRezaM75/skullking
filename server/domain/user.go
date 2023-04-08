package domain

type User struct {
	Id        string `bson:"omitempty"`
	Username  string
	Email     string
	Password  string
	CreatedAt int64 `bson:"created_at"`
}

type UserRepository interface {
	Create(u User) (*User, error)
}

type UserService interface {
	Create(email, username, password string) (*User, error)
}

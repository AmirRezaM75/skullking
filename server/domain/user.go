package domain

import "context"

type User struct {
	username string
	email    string
	password string
}

type UserRepository interface {
	Create(ctx context.Context, u User)
}

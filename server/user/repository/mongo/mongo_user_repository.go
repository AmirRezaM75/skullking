package mongo

import (
	"context"
	"github.com/AmirRezaM75/skull-king/domain"
)

type mongoUserRepository struct {
}

func newMongoUserRepository() domain.UserRepository {
	return mongoUserRepository{}
}

func (ur mongoUserRepository) Create(ctx context.Context, user domain.User) {

}

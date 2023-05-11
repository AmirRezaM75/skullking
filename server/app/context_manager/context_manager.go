package context_manager

import (
	"context"
	"github.com/AmirRezaM75/skull-king/domain"
)

const userCtxKey = "user"

func WithUser(ctx context.Context, user *domain.User) context.Context {
	return context.WithValue(ctx, userCtxKey, user)
}

func GetUser(ctx context.Context) *domain.User {
	user, ok := ctx.Value(userCtxKey).(*domain.User)

	if !ok {
		return nil
	}

	return user
}

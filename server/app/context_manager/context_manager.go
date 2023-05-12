package context_manager

import (
	"context"
	"github.com/AmirRezaM75/skull-king/domain"
)

const userCtxKey = "user"
const requestParamsCtxKey = "params"

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

func WithRequestParams(ctx context.Context, params map[string]string) context.Context {
	return context.WithValue(ctx, requestParamsCtxKey, params)
}

func GetRequestParams(ctx context.Context) map[string]string {
	return ctx.Value(requestParamsCtxKey).(map[string]string)
}

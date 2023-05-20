package services

import (
	"context"
	"github.com/AmirRezaM75/skull-king/models"
)

const userCtxKey = "user"
const requestParamsCtxKey = "params"

type ContextService struct{}

func (cs ContextService) WithUser(ctx context.Context, user *models.User) context.Context {
	return context.WithValue(ctx, userCtxKey, user)
}

func (cs ContextService) GetUser(ctx context.Context) *models.User {
	user, ok := ctx.Value(userCtxKey).(*models.User)

	if !ok {
		return nil
	}

	return user
}

func (cs ContextService) WithRequestParams(ctx context.Context, params map[string]string) context.Context {
	return context.WithValue(ctx, requestParamsCtxKey, params)
}

func (cs ContextService) GetRequestParams(ctx context.Context) map[string]string {
	return ctx.Value(requestParamsCtxKey).(map[string]string)
}

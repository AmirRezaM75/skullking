package repositories

import (
	"context"
	"github.com/AmirRezaM75/skull-king/contracts"
	"github.com/redis/go-redis/v9"
	"time"
)

const TokenKey = "token_"

type tokenRepository struct {
	db *redis.Client
}

func NewTokenRepository(db *redis.Client) contracts.TokenRepository {
	return tokenRepository{
		db: db,
	}
}

func (tr tokenRepository) FindByEmail(email string) string {
	key := TokenKey + email
	value, _ := tr.db.Get(context.Background(), key).Result()

	return value
}

func (tr tokenRepository) Create(email, token string, expiration time.Duration) error {
	key := TokenKey + email

	return tr.db.Set(context.Background(), key, token, expiration).Err()
}

func (tr tokenRepository) DeleteByEmail(email string) error {
	key := TokenKey + email

	return tr.db.Del(context.Background(), key).Err()
}

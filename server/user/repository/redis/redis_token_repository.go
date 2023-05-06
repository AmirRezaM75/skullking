package redis

import (
	"context"
	"github.com/AmirRezaM75/skull-king/domain"
	"github.com/redis/go-redis/v9"
	"time"
)

const TokenKey = "token_"

type redisTokenRepository struct {
	db *redis.Client
}

func NewRedisTokenRepository(db *redis.Client) domain.TokenRepository {
	return redisTokenRepository{
		db: db,
	}
}

func (tr redisTokenRepository) FindByEmail(email string) string {
	key := TokenKey + email
	value, _ := tr.db.Get(context.Background(), key).Result()

	return value
}

func (tr redisTokenRepository) Create(email, token string, expiration time.Duration) error {
	key := TokenKey + email

	return tr.db.Set(context.Background(), key, token, expiration).Err()
}

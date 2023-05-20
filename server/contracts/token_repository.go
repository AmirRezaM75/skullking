package contracts

import "time"

type TokenRepository interface {
	FindByEmail(email string) string
	Create(email, token string, expiration time.Duration) error
	DeleteByEmail(email string) error
}

package support

import (
	"github.com/golang-jwt/jwt/v5"
	"io/ioutil"
	"os"
	"time"
)

func GenerateJWT(id string) (string, error) {
	claims := jwt.RegisteredClaims{
		ID:        id,
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
		Issuer:    os.Getenv("APP_URL"),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	key, _ := ioutil.ReadFile("private.key")
	return token.SignedString(key)
}

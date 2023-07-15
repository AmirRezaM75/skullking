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
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 3)),
		Issuer:    os.Getenv("APP_URL"),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	key, _ := ioutil.ReadFile("private.key")
	return token.SignedString(key)
}

func ParseJWT(token string) (*jwt.RegisteredClaims, error) {
	key, _ := ioutil.ReadFile("private.key")
	claims := &jwt.RegisteredClaims{}

	_, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		return key, nil
	})

	if err != nil {
		return nil, err
	}

	return claims, nil
}

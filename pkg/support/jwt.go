package support

import (
	"github.com/golang-jwt/jwt/v5"
	"io/ioutil"
)

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

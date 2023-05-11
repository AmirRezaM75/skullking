package middlewares

import (
	"errors"
	"fmt"
	"github.com/AmirRezaM75/skull-king/pkg/support"
	"net/http"
	"strings"
)

type Authenticate struct {
}

func (vs Authenticate) Execute(w http.ResponseWriter, r *http.Request) error {
	auth := r.Header.Get("Authorization")

	if auth == "" {
		return errors.New("unauthenticated")
	}

	token := strings.Replace(auth, "Bearer ", "", 1)

	claims, err := support.ParseJWT(token)

	if err != nil {
		return errors.New("unauthenticated")
	}

	fmt.Println(claims)
	return nil
}

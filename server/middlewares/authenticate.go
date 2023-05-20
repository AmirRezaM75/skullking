package middlewares

import (
	"github.com/AmirRezaM75/skull-king/contracts"
	"github.com/AmirRezaM75/skull-king/pkg/support"
	"github.com/AmirRezaM75/skull-king/services"
	"net/http"
	"strings"
)

type Authenticate struct {
	UserService contracts.UserService
}

func (a Authenticate) Handle(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authorizationHeader := r.Header.Get("Authorization")

		if authorizationHeader == "" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		token := strings.Replace(authorizationHeader, "Bearer ", "", 1)

		claims, err := support.ParseJWT(token)

		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		user := a.UserService.FindById(claims.ID)

		if user == nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		ctx := services.ContextService{}.WithUser(r.Context(), user)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

package middlewares

import (
	"github.com/AmirRezaM75/skull-king/app/context_manager"
	"github.com/AmirRezaM75/skull-king/domain"
	"github.com/AmirRezaM75/skull-king/pkg/support"
	"net/http"
	"strings"
)

type Authenticate struct {
	UserService domain.UserService
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

		ctx := context_manager.WithUser(r.Context(), user)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

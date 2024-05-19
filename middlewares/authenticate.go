package middlewares

import (
	"net/http"
	"skullking/contracts"
	"skullking/pkg/support"
	"skullking/services"
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

		if claims.Subject == "" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		user := a.UserService.FindById(claims.Subject)

		if user == nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		ctx := services.ContextService{}.WithUser(r.Context(), user)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

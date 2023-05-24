package middlewares

import (
	"github.com/AmirRezaM75/skull-king/pkg/support"
	"net/http"
	"os"
)

type ValidateSignature struct {
}

func (vs ValidateSignature) Handle(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		urlGenerator := support.NewUrlGenerator(
			os.Getenv("APP_URL"),
			os.Getenv("APP_KEY"),
		)

		valid := urlGenerator.HasValidSignature(r.URL)

		if valid {
			next.ServeHTTP(w, r)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
	})
}
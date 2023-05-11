package middlewares

import (
	"github.com/AmirRezaM75/skull-king/pkg/url_generator"
	"net/http"
)

type ValidateSignature struct {
}

func (vs ValidateSignature) Handle(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		urlGenerator := url_generator.NewUrlGenerator()
		valid := urlGenerator.HasValidSignature(r.URL)

		if valid {
			next.ServeHTTP(w, r)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
	})
}

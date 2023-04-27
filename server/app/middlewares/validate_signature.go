package middlewares

import (
	"errors"
	"github.com/AmirRezaM75/skull-king/pkg/url_generator"
	"net/http"
)

type ValidateSignature struct {
}

func (vs ValidateSignature) Execute(w http.ResponseWriter, r *http.Request) error {
	urlGenerator := url_generator.NewUrlGenerator()
	valid := urlGenerator.HasValidSignature(r.URL)

	if valid {
		return nil
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusUnauthorized)

	return errors.New("invalid signature")
}

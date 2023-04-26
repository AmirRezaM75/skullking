package middlewares

import (
	"errors"
	"github.com/AmirRezaM75/skull-king/pkg/router"
	"net/http"
	"os"
)

type CorsPolicy struct {
	next router.Middleware
}

func (cp CorsPolicy) Execute(w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Access-Control-Allow-Origin", os.Getenv("FRONTEND_URL"))
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusNoContent)
		return errors.New("OPTIONS HTTP method received")
	}

	return nil
}

func (cp CorsPolicy) Next(m router.Middleware) {
	cp.next = m
}

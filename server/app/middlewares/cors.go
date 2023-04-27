package middlewares

import (
	"errors"
	"net/http"
	"os"
)

type CorsPolicy struct {
}

func (cp CorsPolicy) Execute(w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Access-Control-Allow-Origin", os.Getenv("FRONTEND_URL"))
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusNoContent)
		// TODO: It's not actually an error
		return errors.New("OPTIONS HTTP method received")
	}

	return nil
}

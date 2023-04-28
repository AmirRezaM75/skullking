package middlewares

import (
	"net/http"
	"os"
)

type CorsPolicy struct {
}

func (cp CorsPolicy) Execute(w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Access-Control-Allow-Origin", os.Getenv("FRONTEND_URL"))
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Max-Age", "3600")

	return nil
}

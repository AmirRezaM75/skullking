package http

import (
	"net/http"
)

func (userHandler UserHandler) emailVerificationNotification(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

}

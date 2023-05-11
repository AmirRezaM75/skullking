package http

import (
	"github.com/AmirRezaM75/skull-king/app/context_manager"
	"net/http"
)

func (userHandler UserHandler) emailVerificationNotification(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	user := context_manager.GetUser(r.Context())

	err := userHandler.service.SendEmailVerificationNotification(user.Id.Hex(), user.Email)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}

	w.WriteHeader(http.StatusAccepted)
}

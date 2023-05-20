package handlers

import (
	"github.com/AmirRezaM75/skull-king/services"
	"net/http"
)

func (userHandler UserHandler) verifyEmail(w http.ResponseWriter, r *http.Request) {
	params := services.ContextService{}.GetRequestParams(r.Context())
	userId := params["id"]
	userHandler.service.FindById(userId)
	userHandler.service.MarkEmailAsVerified(userId)
}

func (userHandler UserHandler) emailVerificationNotification(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	user := services.ContextService{}.GetUser(r.Context())

	err := userHandler.service.SendEmailVerificationNotification(user.Id.Hex(), user.Email)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}

	w.WriteHeader(http.StatusAccepted)
}

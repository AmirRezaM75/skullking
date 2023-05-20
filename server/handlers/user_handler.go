package handlers

import (
	"encoding/json"
	"github.com/AmirRezaM75/skull-king/contracts"
	"github.com/AmirRezaM75/skull-king/middlewares"
	"github.com/AmirRezaM75/skull-king/pkg/router"
	"github.com/AmirRezaM75/skull-king/pkg/validator"
	"net/http"
)

type ErrorResponse struct {
	Message string `json:"message"`
}

type UserHandler struct {
	service   contracts.UserService
	validator validator.Validator
}

func NewUserHandler(userService contracts.UserService, validator validator.Validator, r *router.Router) {
	var handler = UserHandler{
		service:   userService,
		validator: validator,
	}

	r.Get("/verify-email/:id/:hash", handler.verifyEmail).
		Middleware(middlewares.ValidateSignature{})
	r.Post("/register", handler.register)
	r.Post("/login", handler.login)
	r.Post("/email/verification-notification", handler.emailVerificationNotification).
		Middleware(middlewares.Authenticate{UserService: userService})
	r.Post("/forgot-password", handler.forgotPassword)
	r.Post("/reset-password", handler.resetPassword)
}

func decoder(payload any, w http.ResponseWriter, r *http.Request) error {
	d := json.NewDecoder(r.Body)

	d.DisallowUnknownFields()

	err := d.Decode(payload)

	if err != nil {
		var response = ErrorResponse{
			Message: "Can't decode payload.",
		}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return err
	}

	return nil
}

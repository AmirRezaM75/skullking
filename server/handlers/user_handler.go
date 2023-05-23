package handlers

import (
	"encoding/json"
	"github.com/AmirRezaM75/skull-king/contracts"
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

func NewUserHandler(userService contracts.UserService, validator validator.Validator) UserHandler {
	return UserHandler{
		service:   userService,
		validator: validator,
	}
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

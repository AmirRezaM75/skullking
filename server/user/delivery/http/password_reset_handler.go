package http

import (
	"encoding/json"
	"net/http"
)

type ForgotPasswordRequest struct {
	Email string `json:"email" validate:"required,email"`
}

func (userHandler UserHandler) forgotPassword(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var payload ForgotPasswordRequest

	if err := decoder(&payload, w, r); err != nil {
		return
	}

	validationError := userHandler.validator.ValidateStruct(payload)

	if validationError != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		json.NewEncoder(w).Encode(validationError)
		return
	}

	err := userHandler.service.SendResetLink(payload.Email)

	if err != nil {
		var response struct {
			Message string `json:"message"`
		}
		w.WriteHeader(http.StatusInternalServerError)
		response.Message = err.Error()
		json.NewEncoder(w).Encode(response)
		return
	}

	w.WriteHeader(http.StatusAccepted)
}

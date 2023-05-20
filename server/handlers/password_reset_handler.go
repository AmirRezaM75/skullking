package handlers

import (
	"encoding/json"
	"net/http"
)

func (userHandler UserHandler) forgotPassword(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	type ForgotPasswordRequest struct {
		Email string `json:"email" validate:"required,email"`
	}

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

func (userHandler UserHandler) resetPassword(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	type ResetPasswordRequest struct {
		Email    string `json:"email" validate:"required,email"`
		Token    string `json:"token" validate:"required"`
		Password string `json:"password" validate:"required,min=6,max=255"`
	}

	var payload ResetPasswordRequest

	if err := decoder(&payload, w, r); err != nil {
		return
	}

	validationError := userHandler.validator.ValidateStruct(payload)

	if validationError != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		_ = json.NewEncoder(w).Encode(validationError)
		return
	}

	err := userHandler.service.ResetPassword(payload.Email, payload.Password, payload.Token)

	if err != nil {
		var response struct {
			Message string `json:"message"`
		}
		w.WriteHeader(http.StatusInternalServerError)
		response.Message = err.Error()
		_ = json.NewEncoder(w).Encode(response)
		return
	}

	w.WriteHeader(http.StatusOK)
}

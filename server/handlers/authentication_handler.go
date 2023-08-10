package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/AmirRezaM75/skull-king/pkg/support"
	"github.com/AmirRezaM75/skull-king/responses"
	"net/http"
)

type CreateUserRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Username string `json:"username" validate:"required,min=3,max=32,alphanum"`
	Password string `json:"password" validate:"required,min=6,max=255"`
}

func (userHandler UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var payload CreateUserRequest

	if err := decoder(&payload, w, r); err != nil {
		return
	}

	validationError := userHandler.validator.ValidateStruct(payload)

	if validationError != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		json.NewEncoder(w).Encode(validationError)
		return
	}

	usernameExists := userHandler.service.ExistsByUsername(payload.Username)
	emailExists := userHandler.service.ExistsByEmail(payload.Email)

	if usernameExists || emailExists {
		w.WriteHeader(http.StatusUnprocessableEntity)

		var r = struct {
			Message string            `json:"message"`
			Errors  map[string]string `json:"errors"`
		}{
			Message: "The given data is invalid.",
			Errors:  map[string]string{},
		}

		if emailExists {
			r.Errors["email"] = "The email has already been taken."
		}

		if usernameExists {
			r.Errors["username"] = "The username has already been taken."
		}

		json.NewEncoder(w).Encode(r)
		return
	}

	user, err := userHandler.service.Create(payload.Email, payload.Username, payload.Password)

	if err != nil {
		var response = ErrorResponse{
			Message: "Can not persist user into database.",
		}

		fmt.Println("Can not persist user into database: ", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	_ = userHandler.service.SendEmailVerificationNotification(user.Id.Hex(), user.Email)

	token, err := support.GenerateJWT(user.Id.Hex())

	if err != nil {
		http.Error(w, "Can't generate JWT.", http.StatusInternalServerError)
		return
	}

	var response responses.Authentication

	response.User.Id = user.Id.Hex()
	response.User.Email = user.Email
	response.User.Username = user.Username
	response.User.Verified = false
	response.Token = token

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

func (userHandler UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	payload := struct {
		Identifier string `json:"identifier"`
		Password   string `json:"password"`
	}{}

	if err := decoder(&payload, w, r); err != nil {
		return
	}

	user := userHandler.service.FindByUsernameOrEmail(payload.Identifier)

	if user == nil {
		var response = ErrorResponse{
			Message: "These credentials doesn't match our records.",
		}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	err := support.VerifyPassword(user.Password, payload.Password)

	if err != nil {
		var response = ErrorResponse{
			Message: "These credentials doesn't match our records.",
		}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	var response responses.Authentication

	token, err := support.GenerateJWT(user.Id.Hex())

	if err != nil {
		http.Error(w, "Can't generate JWT.", http.StatusInternalServerError)
		return
	}

	response.User.Id = user.Id.Hex()
	response.User.Email = user.Email
	response.User.Username = user.Username
	response.User.Verified = user.EmailVerifiedAt != nil
	response.Token = token
	json.NewEncoder(w).Encode(response)
}

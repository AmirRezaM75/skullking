package http

import (
	"encoding/json"
	"github.com/AmirRezaM75/skull-king/domain"
	"github.com/AmirRezaM75/skull-king/pkg/support"
	"github.com/AmirRezaM75/skull-king/pkg/validator"
	"net/http"
)

type ErrorResponse struct {
	Message string `json:"message"`
}

type UserHandler struct {
	Service   domain.UserService
	validator validator.Validator
}

func NewUserHandler(userService domain.UserService, validator validator.Validator) {
	var handler = UserHandler{
		Service:   userService,
		validator: validator,
	}

	http.HandleFunc("/register", handler.register)
	http.HandleFunc("/login", handler.login)

}

type CreateUserRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Username string `json:"username" validate:"required,min=3,max=255"`
	Password string `json:"password" validate:"required,min=6,max=255"`
}

func (userHandler UserHandler) register(w http.ResponseWriter, r *http.Request) {
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

	user, err := userHandler.Service.Create(payload.Email, payload.Username, payload.Password)

	var response struct {
		User struct {
			Email    string `json:"email"`
			Username string `json:"username"`
		} `json:"user"`
		Token string `json:"token"`
	}

	token, err := support.GenerateJWT(user.Id.Hex())

	if err != nil {
		http.Error(w, "Can't generate JWT.", http.StatusInternalServerError)
		return
	}

	response.User.Email = user.Email
	response.User.Username = user.Username
	response.Token = token
	json.NewEncoder(w).Encode(response)
}

func (userHandler UserHandler) login(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	payload := struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}{}

	if err := decoder(&payload, w, r); err != nil {
		return
	}

	user := userHandler.Service.FindByUsername(payload.Username)

	if user == nil {
		var response = ErrorResponse{
			Message: "You have entered an invalid username or password",
		}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	err := support.CompareHashAndPassword(user.Password, payload.Password)

	if err != nil {
		var response = ErrorResponse{
			Message: "You have entered an invalid username or password",
		}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	var response struct {
		User struct {
			Email    string `json:"email"`
			Username string `json:"username"`
		} `json:"user"`
		Token string `json:"token"`
	}

	token, err := support.GenerateJWT(user.Id.Hex())

	if err != nil {
		http.Error(w, "Can't generate JWT.", http.StatusInternalServerError)
		return
	}

	response.User.Email = user.Email
	response.User.Username = user.Username
	response.Token = token
	json.NewEncoder(w).Encode(response)
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

package http

import (
	"encoding/json"
	"github.com/AmirRezaM75/skull-king/app/middlewares"
	"github.com/AmirRezaM75/skull-king/domain"
	"github.com/AmirRezaM75/skull-king/pkg/router"
	"github.com/AmirRezaM75/skull-king/pkg/support"
	"github.com/AmirRezaM75/skull-king/pkg/validator"
	"net/http"
)

type ErrorResponse struct {
	Message string `json:"message"`
}

type UserHandler struct {
	service   domain.UserService
	validator validator.Validator
}

func NewUserHandler(userService domain.UserService, validator validator.Validator, r *router.Router) {
	var handler = UserHandler{
		service:   userService,
		validator: validator,
	}

	r.Get("/verify-email/:id/:hash", handler.verifyEmail).Middleware(middlewares.ValidateSignature{})
	r.Post("/register", handler.register)
	r.Post("/login", handler.login)
	r.Post("/email/verification-notification", handler.emailVerificationNotification).Middleware(middlewares.Authenticate{})
	r.Post("/forgot-password", handler.forgotPassword) // TODO: Guest middleware
	r.Post("/reset-password", handler.resetPassword)   // TODO: Guest middleware
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

	_ = userHandler.service.SendEmailVerificationNotification(user.Id.Hex(), user.Email)

	var response struct {
		User struct {
			Email    string `json:"email"`
			Username string `json:"username"`
			Verified bool   `json:"verified"`
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
	response.User.Verified = false
	response.Token = token

	w.WriteHeader(http.StatusCreated)
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

	user := userHandler.service.FindByUsername(payload.Username)

	if user == nil {
		var response = ErrorResponse{
			Message: "You have entered an invalid username or password",
		}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	err := support.VerifyPassword(user.Password, payload.Password)

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
			Verified bool   `json:"verified"`
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
	response.User.Verified = user.EmailVerifiedAt != nil
	response.Token = token
	json.NewEncoder(w).Encode(response)
}

func (userHandler UserHandler) verifyEmail(w http.ResponseWriter, r *http.Request) {
	params := r.Context().Value("params").(map[string]string)
	userId := params["id"]
	userHandler.service.FindById(userId)
	userHandler.service.MarkEmailAsVerified(userId)
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

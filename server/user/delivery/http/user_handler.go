package http

import (
	"encoding/json"
	"github.com/AmirRezaM75/skull-king/domain"
	"net/http"
)

type UserHandler struct {
	Service domain.UserService
}

func NewUserHandler(userService domain.UserService) {
	var handler = UserHandler{
		Service: userService,
	}

	http.HandleFunc("/register", handler.register)

}

func (userHandler UserHandler) register(w http.ResponseWriter, r *http.Request) {
	d := json.NewDecoder(r.Body)
	d.DisallowUnknownFields()
	payload := struct {
		Email    string `json:"email"`
		Username string `json:"username"`
		Password string `json:"password"`
	}{}

	err := d.Decode(&payload)

	if err != nil {
		http.Error(w, "Can't decode payload.", http.StatusBadRequest)
		return
	}

	userHandler.Service.Create(payload.Email, payload.Username, payload.Password)
}

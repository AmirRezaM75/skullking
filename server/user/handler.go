package user

import (
	"fmt"
	"net/http"
)

type Handler struct{}

func NewHandler() Handler {
	return Handler{}
}

func (h Handler) Register(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Register")
}

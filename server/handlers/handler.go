package handlers

import (
	"encoding/json"
	"net/http"
)

func decoder(payload any, w http.ResponseWriter, r *http.Request) error {
	d := json.NewDecoder(r.Body)

	d.DisallowUnknownFields()

	err := d.Decode(payload)

	if err != nil {
		errorResponse(w, "Can not decode payload.", http.StatusUnprocessableEntity)
		return err
	}

	return nil
}

func errorResponse(w http.ResponseWriter, message string, status int) {
	var response = struct {
		Message string `json:"message"`
	}{
		Message: message,
	}
	json.NewEncoder(w).Encode(response)
	w.WriteHeader(status)
}

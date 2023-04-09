package api

import (
	"encoding/json"
	"net/http"
)

func render(w http.ResponseWriter, statusCode int, v any) {
	w.Header().Add("content-type", "application/json")
	w.WriteHeader(statusCode)

	err := json.NewEncoder(w).Encode(v)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

type ErrResponse struct {
	StatusCode int    `json:"-"`
	Message    string `json:"message"`
	Err        error  `json:"-"`
}

func ErrInternalServer(err error) ErrResponse {
	return ErrResponse{
		StatusCode: http.StatusInternalServerError,
		Message:    err.Error(),
		Err:        err,
	}
}

func renderError(w http.ResponseWriter, err ErrResponse) {
	render(w, err.StatusCode, err)
}

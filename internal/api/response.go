package api

import (
	"encoding/json"
	"net/http"
)

type ErrorResponse struct {
	Error 	string 	`json:"error"`
	Code 	string 	`json:"code,omitempty"`
	Details string 	`json:"details,omitempty"`
}

type SuccessResponse struct {
	Status 	string 		`json:"status"`
	Message string 		`json:"message"`
	Data 	interface{} `json:"data"`
}

func WriteJSON(w http.ResponseWriter, status int, v interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(v)
}

func WriteError(w http.ResponseWriter, status int, err error) {
	WriteJSON(w, status, ErrorResponse{
		Error: err.Error(),
	})
}

func WriteSuccess(w http.ResponseWriter, message string, data interface{}) {
	WriteJSON(w, http.StatusOK, SuccessResponse{
		Status: "success",
		Message: message,
		Data: data,
	})
}

func WriteCreated(w http.ResponseWriter, message string, data interface{}) {
	WriteJSON(w, http.StatusCreated, SuccessResponse{
		Status:  "created",
		Message: message,
		Data:    data,
	})
}

func WriteNoContent(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNoContent)
}
package main

import (
	"encoding/json"
	"net/http"
)

type ErrorResponse struct {
	Status      int                    `json:"status"`
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	Meta        map[string]interface{} `json:"meta"`
}

func writeErrorResponse(w http.ResponseWriter, errR ErrorResponse) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(errR.Status)
	_ = json.NewEncoder(w).Encode(errR)
}

func newErrorResponse(status int, name, description string, meta map[string]any) ErrorResponse {
	return ErrorResponse{
		Status:      status,
		Name:        name,
		Description: description,
		Meta:        meta,
	}
}

func InternalServerError(w http.ResponseWriter, description string, meta map[string]any) ErrorResponse {
	errResponse := newErrorResponse(http.StatusInternalServerError, "INTERNAL_SERVER_ERROR", description, meta)
	writeErrorResponse(w, errResponse)
	return errResponse
}

func ValidationError(w http.ResponseWriter, description string, meta map[string]any) ErrorResponse {
	errResponse := newErrorResponse(http.StatusBadRequest, "VALIDATION_ERROR", description, meta)
	writeErrorResponse(w, errResponse)
	return errResponse
}

func NotFoundError(w http.ResponseWriter, description string, meta map[string]any) ErrorResponse {
	errResponse := newErrorResponse(http.StatusNotFound, "NOT_FOUND_ERROR", description, meta)
	writeErrorResponse(w, errResponse)
	return errResponse
}

func WriteStructToJson(w http.ResponseWriter, data interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(data)
}

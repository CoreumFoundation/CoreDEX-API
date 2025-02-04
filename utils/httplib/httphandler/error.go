package handler

import (
	"encoding/json"
	"net/http"

	"github.com/CoreumFoundation/CoreDEX-API/utils/logger"
)

// TODO: Add multiple errors handling.
// TODO: Make this type interface so each package can define own errors.
type APIError struct {
	Status int    `json:"-"`
	Name   string `json:"name"`
}

// Satisfy the error interface.
func (apiErr APIError) Error() string {
	return apiErr.Name
}

// Follow internal errors format convention.
func (apiErr APIError) MarshalJSON() ([]byte, error) {
	type Errors struct {
		Name string `json:"name"`
	}
	res := struct {
		Errors            []Errors `json:"errors"`
		Alerts            []Errors `json:"alerts"` // Compatibility.
		DeprecationNotice string   `json:"_"`
	}{
		Errors:            []Errors{{Name: apiErr.Name}},
		Alerts:            []Errors{{Name: apiErr.Name}},
		DeprecationNotice: "«alerts» field is deprecated please use «errors» instead",
	}

	return json.Marshal(&res)
}

func NewAPIError(status int, name string) *APIError {
	return &APIError{
		Status: status,
		Name:   name,
	}
}

var internalServerError = &APIError{
	Status: http.StatusInternalServerError,
	Name:   "server.internal_error",
}

func serveAPIError(w http.ResponseWriter, err *APIError) {
	logger.Warnf("HTTP %d - %s", err.Status, err)

	w.WriteHeader(err.Status)
	if err := json.NewEncoder(w).Encode(err); err != nil {
		panic(err)
	}
}

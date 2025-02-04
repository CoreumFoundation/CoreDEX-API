package handler

import (
	"net/http"
)

var (
	ErrUnauthorized   = &APIError{Name: "unauthorized", Status: http.StatusUnauthorized}
	ErrInvalidInput   = &APIError{Name: "invalid-input", Status: http.StatusBadRequest}
	ErrNetworkInvalid = &APIError{Name: "invalid-network", Status: http.StatusBadRequest}
	ErrNotFound       = &APIError{Name: "not-found", Status: http.StatusNotFound}
	ErrAddressInvalid = &APIError{Name: "address-invalid", Status: http.StatusForbidden}
	ErrUnknown        = &APIError{Name: "unknown", Status: http.StatusUnprocessableEntity}
	ErrBE             = &APIError{Name: "backend-error", Status: http.StatusBadRequest}
)

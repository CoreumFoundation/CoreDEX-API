package handler

import (
	"context"
	"net/http"

	"github.com/google/uuid"
	"github.com/pkg/errors"

	"github.com/CoreumFoundation/CoreDEX-API/utils/logger"
)

const CO_SESSION_ID = "co-session-id"
const httpErrorContextCancelled = 570 // See https://en.wikipedia.org/wiki/List_of_HTTP_status_codes#5xx_server_errors: Value choosen as non-conflicting status code

// The Handler helps to handle errors in one place.
type Handler func(w http.ResponseWriter, r *http.Request) error
type HandlerFunc func(w http.ResponseWriter, r *http.Request)

// ServeHTTP allows our Handler type to satisfy http.Handler.
func (h Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	// Generate a uuid and make it recognizable in session:
	sessionID := uuid.New().String()
	w.Header().Set(CO_SESSION_ID, sessionID)
	// Set a sessionID in the request for the logger:
	r.Header.Set(CO_SESSION_ID, sessionID)

	// Replace context of the request to pass on the sg-session-id for log identification:
	// ignore go lint remark: Value has been deconflicted with sg- before the name of the identifier
	contextRequest := r.WithContext(context.WithValue(r.Context(), CO_SESSION_ID, sessionID))

	// Handles panic.
	defer func() {
		err := recover()
		if err != nil {
			serveAPIError(w, internalServerError)
		}
	}()

	// Handles errors returned by handlers.
	if err := h(w, contextRequest); err != nil {
		if r.Context().Err() == context.Canceled && errors.Is(err, context.Canceled) {
			// Internal context cancelled error code to distinguish between other possible errors created by infrastructure
			w.WriteHeader(httpErrorContextCancelled)
			w.Write([]byte("Context cancel occurred for session " + sessionID))
			return
		}

		switch e := err.(type) {
		case *APIError:
			// We can retrieve the status here and write out a specific HTTP status code.
			serveAPIError(w, e)
		default:
			// Any error types we don't expect results into HTTP 500.
			logger.Errorf("HTTP %d session %s, %v", internalServerError.Status, sessionID, e)
			serveAPIError(w, internalServerError)
		}
	}
}

package http

import (
	"io"
	"net/http"
)

const (
	// Required since a body, once read from the request, can not be read again, and we need to leave it somewhere for later use: This const indicates where it is stored on the request
	BODY = "body"
)

func Body(r *http.Request) []byte {
	return []byte(r.Header.Get(BODY))
}

// Reads the body from a post or put, sets the body as byte array on the request for later use
func ReadBody(r *http.Request) ([]byte, error) {
	b, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, nil
	}
	// Set the body on the request for later use
	r.Header.Set(BODY, string(b))
	return b, nil
}

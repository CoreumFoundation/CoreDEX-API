package handler

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHandler_ServeHTTP(t *testing.T) {
	t.Run("api_error", func(t *testing.T) {
		rr := httptest.NewRecorder()
		req, err := http.NewRequest(http.MethodGet, "/", nil)
		assert.NoError(t, err)

		h := Handler(func(w http.ResponseWriter, r *http.Request) error {
			return NewAPIError(http.StatusUnprocessableEntity, "invalid.params")
		})

		h.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusUnprocessableEntity, rr.Code)
		assert.Equal(t, "{\"alerts\":[{\"name\":\"invalid.params\"}]}\n", rr.Body.String())
	})

	t.Run("error", func(t *testing.T) {
		rr := httptest.NewRecorder()
		req, err := http.NewRequest(http.MethodGet, "/", nil)
		assert.NoError(t, err)

		h := Handler(func(w http.ResponseWriter, r *http.Request) error {
			return errors.New("unhandled error")
		})

		h.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusInternalServerError, rr.Code)
	})
}

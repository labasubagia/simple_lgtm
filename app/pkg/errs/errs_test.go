package errs

import (
	"errors"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMapHttp(t *testing.T) {
	t.Run("NilError", func(t *testing.T) {
		status, msg := MapHttp(nil)
		assert.Equal(t, http.StatusOK, status)
		assert.Empty(t, msg)
	})

	t.Run("InternalError", func(t *testing.T) {
		err := NewInternal(errors.New("internal failure"))
		status, msg := MapHttp(err)
		assert.Equal(t, http.StatusInternalServerError, status)
		assert.Contains(t, msg, "internal failure")
	})

	t.Run("NotFoundError", func(t *testing.T) {
		err := NewNotFound(errors.New("not found"))
		status, msg := MapHttp(err)
		assert.Equal(t, http.StatusNotFound, status)
		assert.Contains(t, msg, "not found")
	})

	t.Run("InvalidInputError", func(t *testing.T) {
		err := NewInvalidInput(errors.New("bad input"))
		status, msg := MapHttp(err)
		assert.Equal(t, http.StatusBadRequest, status)
		assert.Contains(t, msg, "bad input")
	})

	t.Run("UnknownAppErrorCode", func(t *testing.T) {
		unknownErr := &appError{Code: "UNKNOWN_CODE", Err: errors.New("unknown")}
		status, msg := MapHttp(unknownErr)
		assert.Equal(t, http.StatusInternalServerError, status)
		assert.Contains(t, msg, "Unknown app error")
	})

	t.Run("NonAppError", func(t *testing.T) {
		err := errors.New("plain error")
		status, msg := MapHttp(err)
		assert.Equal(t, http.StatusInternalServerError, status)
		assert.Contains(t, msg, "Unexpected error")
	})
}

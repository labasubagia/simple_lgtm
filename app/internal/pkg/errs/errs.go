package errs

import (
	"errors"
	"fmt"
	"net/http"
)

type ErrCode string

const (
	ErrInternal     ErrCode = "INTERNAL"
	ErrNotFound     ErrCode = "NOT_FOUND"
	ErrInvalidInput ErrCode = "INVALID_INPUT"
)

type appError struct {
	Code ErrCode
	Err  error
}

func (a *appError) Error() string {
	return fmt.Sprintf("code: %s, error: %s", a.Code, a.Err.Error())
}

func NewInternal(err error) error {
	if err == nil {
		return nil
	}
	return &appError{Code: ErrInternal, Err: err}
}

func NewInvalidInput(err error) error {
	if err == nil {
		return nil
	}
	return &appError{Code: ErrInvalidInput, Err: err}
}

func NewNotFound(err error) error {
	if err == nil {
		return nil
	}
	return &appError{Code: ErrNotFound, Err: err}
}

func MapHttp(err error) (statusCode int, message string) {
	if err == nil {
		return http.StatusOK, ""
	}

	var e *appError
	if errors.As(err, &e) {
		switch e.Code {
		case ErrInternal:
			return http.StatusInternalServerError, err.Error()
		case ErrNotFound:
			return http.StatusNotFound, err.Error()
		case ErrInvalidInput:
			return http.StatusBadRequest, err.Error()
		default:
			return http.StatusInternalServerError, fmt.Sprintf("Unknown error: %s", err.Error())
		}
	}

	return http.StatusInternalServerError, fmt.Sprintf("Unexpected error: %s", err.Error())
}

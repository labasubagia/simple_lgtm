package errs

import (
	"errors"
	"fmt"
	"net/http"
)

type errCode string

const (
	codeInternal     errCode = "INTERNAL"
	codeNotFound     errCode = "NOT_FOUND"
	codeInvalidInput errCode = "INVALID_INPUT"
)

type appError struct {
	Code errCode
	Err  error
}

func (a *appError) Error() string {
	return fmt.Sprintf("code: %s, error: %s", a.Code, a.Err.Error())
}

func NewInternal(err error) error {
	if err == nil {
		return nil
	}
	return &appError{Code: codeInternal, Err: err}
}

func NewInvalidInput(err error) error {
	if err == nil {
		return nil
	}
	return &appError{Code: codeInvalidInput, Err: err}
}

func NewNotFound(err error) error {
	if err == nil {
		return nil
	}
	return &appError{Code: codeNotFound, Err: err}
}

func MapHttp(err error) (statusCode int, message string) {
	if err == nil {
		return http.StatusOK, ""
	}

	var appErr *appError
	if errors.As(err, &appErr) {
		switch appErr.Code {
		case codeInternal:
			return http.StatusInternalServerError, err.Error()
		case codeNotFound:
			return http.StatusNotFound, err.Error()
		case codeInvalidInput:
			return http.StatusBadRequest, err.Error()
		default:
			return http.StatusInternalServerError, fmt.Sprintf("Unknown app error: %s", err.Error())
		}
	}

	return http.StatusInternalServerError, fmt.Sprintf("Unexpected error: %s", err.Error())
}

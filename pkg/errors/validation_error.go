package errors

import "net/http"

type ValidationError struct {
	message string
	errors  map[string]any
}

func NewValidationError(message string, errors map[string]any) CustomError {
	return &ValidationError{
		message: message,
		errors:  errors,
	}
}

func (e *ValidationError) Error() string {
	if e.message != "" {
		return e.message
	}

	return http.StatusText(http.StatusBadRequest)
}

func (e *ValidationError) StatusCode() int {
	return http.StatusBadRequest
}

func (e *ValidationError) Message() string {
	if e.message != "" {
		return e.message
	}

	return http.StatusText(http.StatusBadRequest)
}

func (e *ValidationError) Errors() map[string]any {
	return e.errors
}

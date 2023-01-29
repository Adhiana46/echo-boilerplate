package errors

import "net/http"

type ForbiddenError struct {
	message string
}

func NewForbiddenError(message string) CustomError {
	return &ForbiddenError{
		message: message,
	}
}

func (e *ForbiddenError) Error() string {
	if e.message != "" {
		return e.message
	}

	return http.StatusText(http.StatusForbidden)
}

func (e *ForbiddenError) StatusCode() int {
	return http.StatusForbidden
}

func (e *ForbiddenError) Message() string {
	if e.message != "" {
		return e.message
	}

	return http.StatusText(http.StatusForbidden)
}

func (e *ForbiddenError) Errors() map[string]any {
	return nil
}

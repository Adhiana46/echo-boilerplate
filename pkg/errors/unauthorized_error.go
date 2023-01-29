package errors

import "net/http"

type UnauthorizedError struct {
	message string
}

func NewUnauthorizedError(message string) CustomError {
	return &UnauthorizedError{
		message: message,
	}
}

func (e *UnauthorizedError) Error() string {
	if e.message != "" {
		return e.message
	}

	return http.StatusText(http.StatusUnauthorized)
}

func (e *UnauthorizedError) StatusCode() int {
	return http.StatusUnauthorized
}

func (e *UnauthorizedError) Message() string {
	if e.message != "" {
		return e.message
	}

	return http.StatusText(http.StatusUnauthorized)
}

func (e *UnauthorizedError) Errors() map[string]any {
	return nil
}

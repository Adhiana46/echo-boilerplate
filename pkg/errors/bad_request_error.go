package errors

import "net/http"

type BadRequestError struct {
	message string
}

func NewBadRequestError(message string) CustomError {
	return &BadRequestError{
		message: message,
	}
}

func (e *BadRequestError) Error() string {
	if e.message != "" {
		return e.message
	}

	return http.StatusText(http.StatusBadRequest)
}

func (e *BadRequestError) StatusCode() int {
	return http.StatusBadRequest
}

func (e *BadRequestError) Message() string {
	if e.message != "" {
		return e.message
	}

	return http.StatusText(http.StatusBadRequest)
}

func (e *BadRequestError) Errors() map[string]any {
	return nil
}

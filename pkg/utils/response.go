package utils

import (
	"net/http"
)

type JsonResponse struct {
	Status     int    `json:"status"`
	Message    string `json:"message"`
	Data       any    `json:"data,omitempty"`
	Errors     any    `json:"errors,omitempty"`
	Pagination any    `json:"pagination,omitempty"`
	Debug      any    `json:"debug,omitempty"`
}

type Pagination struct {
	Size        int `json:"size"`
	Total       int `json:"total"`
	TotalPages  int `json:"total_pages"`
	CurrentPage int `json:"current_page"`
}

func JsonSuccess(status int, message string, data any, pagination any) JsonResponse {
	if message == "" {
		message = http.StatusText(status)
	}

	return JsonResponse{
		Status:     status,
		Message:    message,
		Data:       data,
		Pagination: pagination,
	}
}

func JsonError(status int, message string, errs any, debug any) JsonResponse {
	if message == "" {
		message = http.StatusText(status)
	}

	return JsonResponse{
		Status:  status,
		Message: message,
		Errors:  errs,
		Debug:   debug,
	}
}

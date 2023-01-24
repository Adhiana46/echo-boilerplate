package errors

type CustomError interface {
	Error() string
	StatusCode() int
	Message() string
	Errors() map[string]any
}

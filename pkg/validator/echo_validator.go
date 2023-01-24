package validator

import "github.com/go-playground/validator/v10"

type EchoValidator struct {
	validator *validator.Validate
}

func NewEchoValidator(validator *validator.Validate) *EchoValidator {
	return &EchoValidator{
		validator: validator,
	}
}

func (s *EchoValidator) Validate(i interface{}) error {
	return s.validator.Struct(i)
}

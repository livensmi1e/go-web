package validator

import (
	"go-web/internal/core/ports"

	"github.com/go-playground/validator/v10"
)

type gpValidator struct {
	validator *validator.Validate
}

func NewValidator() ports.Validator {
	return &gpValidator{
		validator: validator.New(validator.WithRequiredStructEnabled()),
	}
}

func (cv *gpValidator) Validate(i interface{}) error {
	return cv.validator.Struct(i)
}

package appconfig

import (
	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

// ValidateStruct validates the given struct value using the configured validator.
func ValidateStruct(i interface{}) error {
	if validate == nil {
		validate = validator.New(validator.WithRequiredStructEnabled())
	}

	return validate.Struct(i)
}

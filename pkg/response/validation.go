package response

import (
	"errors"

	"github.com/go-playground/validator/v10"
)

// ParseValidationErrors converts validation errors into a slice of FieldError
func ParseValidationErrors(err error) []FieldError {
	var validationErrs validator.ValidationErrors
	if errors.As(err, &validationErrs) {
		var fieldErrors []FieldError
		for _, e := range validationErrs {
			fieldErrors = append(fieldErrors, FieldError{
				Field: e.Field(),
				Tag:   e.Tag(),
				Param: e.Param(),
			})
		}
		return fieldErrors
	}
	return nil
}

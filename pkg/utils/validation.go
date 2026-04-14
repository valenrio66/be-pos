package utils

import (
	"errors"

	"github.com/go-playground/validator/v10"
)

func FormatValidationError(err error) map[string]string {
	errMessages := make(map[string]string)

	var ve validator.ValidationErrors
	if errors.As(err, &ve) {
		for _, e := range ve {
			errMessages[e.Field()] = "This column does not meet the requirements: " + e.Tag()
		}
	} else {
		errMessages["general"] = err.Error()
	}

	return errMessages
}

package utils

import (
	"fmt"
	"slices"
	"strings"

	"github.com/go-playground/validator/v10"
)

var (
	Validate          *validator.Validate
	ValidContentTypes = []string{"blog", "project"}
)

func init() {
	Validate = validator.New(validator.WithRequiredStructEnabled())

	// No need to check error because error comes from empty tag and nil fn
	_ = Validate.RegisterValidation("content_type", func(fl validator.FieldLevel) bool {
		return slices.Contains(ValidContentTypes, fl.Field().String())
	})
}

func FormatValidationError(err error) (string, []string) {
	validationErrors, ok := err.(validator.ValidationErrors)

	if !ok {
		return "invalid input format", nil
	}

	var messages []string

	for _, fieldErr := range validationErrors {
		msg := fmt.Sprintf("field '%s' failed on '%s' validation", fieldErr.Field(), fieldErr.Tag())

		if param := fieldErr.Param(); param != "" {
			switch fieldErr.Tag() {
			case "oneof":
				msg += fmt.Sprintf(" (must be one of: %s)", strings.ReplaceAll(param, " ", ", "))
			case "min":
				msg += fmt.Sprintf(" (minimum: %s)", param)
			case "max":
				msg += fmt.Sprintf(" (maximum: %s)", param)
			default:
				msg += fmt.Sprintf(" (rule: %s)", param)
			}
		}

		// Handle custom hint for custom validation
		switch fieldErr.Tag() {
		case "content_type":
			msg += fmt.Sprintf(" (must be one of: %s)", strings.Join(ValidContentTypes, ", "))
		}

		messages = append(messages, msg)
	}

	return fmt.Sprintf("invalid input: %s", strings.Join(messages, ", ")), messages
}

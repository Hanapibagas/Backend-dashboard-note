package utils

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
)

var (
	// Validator instance
	validate *validator.Validate
)

func init() {
	validate = validator.New()

	// Register custom tag name function to use json tag names in error messages
	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})
}

// ValidationError represents a validation error for a specific field
type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

// Validator handles struct validation
type Validator struct {
	validator *validator.Validate
}

// NewValidator creates a new validator instance
func NewValidator() *Validator {
	return &Validator{
		validator: validate,
	}
}

// Validate validates a struct and returns validation errors if any
func (v *Validator) Validate(s interface{}) []ValidationError {
	err := v.validator.Struct(s)
	if err == nil {
		return nil
	}

	var validationErrors []ValidationError
	var errs validator.ValidationErrors

	if errors, ok := err.(validator.ValidationErrors); ok {
		errs = errors
	} else {
		return []ValidationError{
			{
				Field:   "general",
				Message: "validation failed",
			},
		}
	}

	for _, e := range errs {
		validationError := ValidationError{
			Field:   e.Field(),
			Message: v.getErrorMessage(e),
		}
		validationErrors = append(validationErrors, validationError)
	}

	return validationErrors
}

// getErrorMessage returns a user-friendly error message for a validation error
func (v *Validator) getErrorMessage(e validator.FieldError) string {
	field := e.Field()

	switch e.Tag() {
	case "required":
		return fmt.Sprintf("%s is required", field)
	case "email":
		return fmt.Sprintf("%s must be a valid email address", field)
	case "min":
		return fmt.Sprintf("%s must be at least %s characters", field, e.Param())
	case "max":
		return fmt.Sprintf("%s must be at most %s characters", field, e.Param())
	case "eqfield":
		return fmt.Sprintf("%s must match %s", field, e.Param())
	default:
		return fmt.Sprintf("%s is invalid", field)
	}
}

// FormatValidationErrors formats validation errors into a readable format
func FormatValidationErrors(errors []ValidationError) string {
	if len(errors) == 0 {
		return ""
	}

	var messages []string
	for _, e := range errors {
		messages = append(messages, fmt.Sprintf("%s: %s", e.Field, e.Message))
	}

	return strings.Join(messages, "; ")
}

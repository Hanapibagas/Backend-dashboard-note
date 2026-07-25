package utils

import (
	"fmt"
	"strings"
)

// ValidationError represents a validation error for a specific field
type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
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

package config

import (
	"fmt"
	"strings"
)

// FormatValidationErrors returns a human-readable string of validation errors
func FormatValidationErrors(errors []ValidationError) string {
	if len(errors) == 0 {
		return ""
	}

	var messages []string
	for _, err := range errors {
		messages = append(messages, fmt.Sprintf("- %s: %s", err.Field, err.Error))
	}

	return fmt.Sprintf(
		"Configuration validation failed:\n%s",
		strings.Join(messages, "\n"),
	)
}

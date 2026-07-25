package v1

import (
	"net/http"
)

// FieldError represents a validation error for a specific field
type FieldError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
	Key     string `json:"key"`
}

// ErrorResponse represents a standard API error response
type ErrorResponse struct {
	StatusCode int          `json:"status_code"`
	Status     string       `json:"status"`
	Message    string       `json:"message,omitempty"`
	Errors     []FieldError `json:"errors,omitempty"`

	Path      string `json:"path,omitempty"`
	Timestamp string `json:"timestamp"`
}

func (e *ErrorResponse) Error() string {
	switch {
	case e.Message != "":
		return e.Message
	case e.StatusCode != 0:
		return http.StatusText(e.StatusCode)
	default:
		return "unknown error"
	}
}

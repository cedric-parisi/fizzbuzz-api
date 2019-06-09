package errors

import (
	"net/http"
)

// APIError represents an api error.
type APIError struct {
	code    int
	Message string       `json:"message"`
	Fields  []FieldError `json:"fields,omitempty"`
}

// FieldError represents a validation error for a field.
type FieldError struct {
	Name    string `json:"name"`
	Message string `json:"message"`
}

func (a APIError) Error() string {
	return a.Message
}

// StatusCode return the error code.
func (a APIError) StatusCode() int {
	return a.code
}

// InvalidArgument creates a bad request error.
func InvalidArgument(msg string, fields []FieldError) error {
	return APIError{
		code:    http.StatusBadRequest,
		Message: msg,
		Fields:  fields,
	}
}

// NotFound creates a not found error.
func NotFound(err error) error {
	return APIError{
		code:    http.StatusNotFound,
		Message: err.Error(),
	}
}

// Internal creates an internal server error.
func Internal(err error) error {
	return APIError{
		code:    http.StatusInternalServerError,
		Message: err.Error(),
	}
}

package model

import (
	"errors"
	"fmt"
)

var ErrNotFound = errors.New("not found")

// APIError represents an API error with error code
type APIError struct {
	Code    string
	Message string
}

func (e *APIError) Error() string {
	return e.Message
}

// NewAPIError creates a new APIError
func NewAPIError(code, message string) *APIError {
	return &APIError{
		Code:    code,
		Message: message,
	}
}

// NewAPIErrorf creates a new APIError with formatted message
func NewAPIErrorf(code, format string, args ...interface{}) *APIError {
	return &APIError{
		Code:    code,
		Message: fmt.Sprintf(format, args...),
	}
}

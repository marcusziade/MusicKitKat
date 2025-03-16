// Package errors provides error types for the MusicKitKat SDK.
package errors

import (
	"fmt"
	"strings"
)

// ErrorType represents the type of error.
type ErrorType string

const (
	// ErrorTypeAuthentication represents an authentication error.
	ErrorTypeAuthentication ErrorType = "authentication"

	// ErrorTypeInvalidRequest represents an invalid request error.
	ErrorTypeInvalidRequest ErrorType = "invalid_request"

	// ErrorTypeRateLimit represents a rate limit error.
	ErrorTypeRateLimit ErrorType = "rate_limit"

	// ErrorTypeServer represents a server error.
	ErrorTypeServer ErrorType = "server"

	// ErrorTypeUnknown represents an unknown error.
	ErrorTypeUnknown ErrorType = "unknown"
)

// APIError represents an error returned by the Apple Music API.
type APIError struct {
	// HTTP status code
	StatusCode int `json:"-"`

	// Error type
	Type ErrorType `json:"-"`

	// Error message
	Message string `json:"-"`

	// Error details from the API
	Errors []struct {
		ID     string `json:"id"`
		Title  string `json:"title"`
		Detail string `json:"detail"`
		Status string `json:"status"`
		Code   string `json:"code"`
	} `json:"errors"`
}

// Error returns the error message.
func (e *APIError) Error() string {
	if len(e.Errors) == 0 {
		return fmt.Sprintf("API error (status code: %d)", e.StatusCode)
	}

	var messages []string
	for _, err := range e.Errors {
		messages = append(messages, fmt.Sprintf("%s: %s", err.Title, err.Detail))
	}

	return fmt.Sprintf("API error (status code: %d): %s", e.StatusCode, strings.Join(messages, "; "))
}

// GetType returns the error type based on the status code.
func (e *APIError) GetType() ErrorType {
	switch {
	case e.StatusCode == 401 || e.StatusCode == 403:
		return ErrorTypeAuthentication
	case e.StatusCode >= 400 && e.StatusCode < 500:
		return ErrorTypeInvalidRequest
	case e.StatusCode == 429:
		return ErrorTypeRateLimit
	case e.StatusCode >= 500:
		return ErrorTypeServer
	default:
		return ErrorTypeUnknown
	}
}

// IsAuthenticationError returns true if the error is an authentication error.
func IsAuthenticationError(err error) bool {
	if apiErr, ok := err.(*APIError); ok {
		return apiErr.GetType() == ErrorTypeAuthentication
	}
	return false
}

// IsInvalidRequestError returns true if the error is an invalid request error.
func IsInvalidRequestError(err error) bool {
	if apiErr, ok := err.(*APIError); ok {
		return apiErr.GetType() == ErrorTypeInvalidRequest
	}
	return false
}

// IsRateLimitError returns true if the error is a rate limit error.
func IsRateLimitError(err error) bool {
	if apiErr, ok := err.(*APIError); ok {
		return apiErr.GetType() == ErrorTypeRateLimit
	}
	return false
}

// IsServerError returns true if the error is a server error.
func IsServerError(err error) bool {
	if apiErr, ok := err.(*APIError); ok {
		return apiErr.GetType() == ErrorTypeServer
	}
	return false
}


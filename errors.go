package buchida

import "fmt"

// APIError represents a buchida API error.
type APIError struct {
	StatusCode int
	Message    string
	Code       string
}

func (e *APIError) Error() string {
	return fmt.Sprintf("buchida: %d %s", e.StatusCode, e.Message)
}

// AuthenticationError is returned for 401 responses.
type AuthenticationError struct {
	StatusCode int
	Message    string
}

func (e *AuthenticationError) Error() string {
	return fmt.Sprintf("buchida: %d %s", e.StatusCode, e.Message)
}

// NotFoundError is returned for 404 responses.
type NotFoundError struct {
	StatusCode int
	Message    string
}

func (e *NotFoundError) Error() string {
	return fmt.Sprintf("buchida: %d %s", e.StatusCode, e.Message)
}

// ValidationError is returned for 422 responses.
type ValidationError struct {
	StatusCode int
	Message    string
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("buchida: %d %s", e.StatusCode, e.Message)
}

// RateLimitError is returned for 429 responses.
type RateLimitError struct {
	StatusCode int
	Message    string
}

func (e *RateLimitError) Error() string {
	return fmt.Sprintf("buchida: %d %s", e.StatusCode, e.Message)
}

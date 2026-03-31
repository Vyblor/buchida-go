package buchida

import "fmt"

// Error represents a buchida API error.
type Error struct {
	StatusCode int
	Message    string
	Code       string
}

func (e *Error) Error() string {
	return fmt.Sprintf("buchida: %d %s", e.StatusCode, e.Message)
}

// AuthenticationError is returned for 401 responses.
type AuthenticationError struct{ Error }

// NotFoundError is returned for 404 responses.
type NotFoundError struct{ Error }

// ValidationError is returned for 422 responses.
type ValidationError struct{ Error }

// RateLimitError is returned for 429 responses.
type RateLimitError struct{ Error }

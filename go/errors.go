package screencraft

import (
	"errors"
	"fmt"
	"net/http"
	"time"
)

var (
	// ErrMissingAPIKey is returned when no API key is provided.
	ErrMissingAPIKey = errors.New("screencraft: API key is required")

	// ErrMissingURL is returned when no URL is provided for capture.
	ErrMissingURL = errors.New("screencraft: URL is required")

	// ErrInvalidFormat is returned when an invalid format is specified.
	ErrInvalidFormat = errors.New("screencraft: invalid format specified")

	// ErrInvalidQuality is returned when quality is out of range.
	ErrInvalidQuality = errors.New("screencraft: quality must be between 0 and 100")

	// ErrInvalidViewport is returned when viewport dimensions are invalid.
	ErrInvalidViewport = errors.New("screencraft: viewport dimensions must be positive")

	// ErrContextCanceled is returned when the context is canceled.
	ErrContextCanceled = errors.New("screencraft: context canceled")

	// ErrTimeout is returned when the operation times out.
	ErrTimeout = errors.New("screencraft: operation timed out")
)

// Error represents a ScreenCraft API error.
type Error struct {
	// StatusCode is the HTTP status code.
	StatusCode int

	// Code is the error code from the API.
	Code string

	// Message is the error message.
	Message string

	// Details contains additional error details.
	Details map[string]interface{}

	// RequestID is the unique request ID for debugging.
	RequestID string

	// Err is the underlying error, if any.
	Err error
}

// Error implements the error interface.
func (e *Error) Error() string {
	if e.Code != "" {
		return fmt.Sprintf("screencraft: %s (%s) - status %d", e.Message, e.Code, e.StatusCode)
	}
	if e.Message != "" {
		return fmt.Sprintf("screencraft: %s - status %d", e.Message, e.StatusCode)
	}
	if e.Err != nil {
		return fmt.Sprintf("screencraft: %v", e.Err)
	}
	return fmt.Sprintf("screencraft: HTTP %d", e.StatusCode)
}

// Unwrap returns the underlying error.
func (e *Error) Unwrap() error {
	return e.Err
}

// IsRetryable returns true if the error is retryable.
func (e *Error) IsRetryable() bool {
	switch e.StatusCode {
	case http.StatusTooManyRequests,
		http.StatusServiceUnavailable,
		http.StatusBadGateway,
		http.StatusGatewayTimeout,
		http.StatusInternalServerError:
		return true
	}
	return false
}

// AuthenticationError represents an authentication failure.
type AuthenticationError struct {
	*Error
}

// NewAuthenticationError creates a new AuthenticationError.
func NewAuthenticationError(message string) *AuthenticationError {
	return &AuthenticationError{
		Error: &Error{
			StatusCode: http.StatusUnauthorized,
			Code:       "AUTHENTICATION_ERROR",
			Message:    message,
		},
	}
}

// RateLimitError represents a rate limit exceeded error.
type RateLimitError struct {
	*Error

	// Limit is the rate limit.
	Limit int

	// Remaining is the remaining requests.
	Remaining int

	// ResetAt is when the rate limit resets.
	ResetAt time.Time

	// RetryAfter is the duration to wait before retrying.
	RetryAfter time.Duration
}

// NewRateLimitError creates a new RateLimitError.
func NewRateLimitError(limit, remaining int, resetAt time.Time, retryAfter time.Duration) *RateLimitError {
	return &RateLimitError{
		Error: &Error{
			StatusCode: http.StatusTooManyRequests,
			Code:       "RATE_LIMIT_EXCEEDED",
			Message:    "rate limit exceeded",
		},
		Limit:      limit,
		Remaining:  remaining,
		ResetAt:    resetAt,
		RetryAfter: retryAfter,
	}
}

// ValidationError represents a validation failure.
type ValidationError struct {
	*Error

	// Field is the field that failed validation.
	Field string

	// Constraint is the validation constraint that failed.
	Constraint string
}

// NewValidationError creates a new ValidationError.
func NewValidationError(field, message, constraint string) *ValidationError {
	return &ValidationError{
		Error: &Error{
			StatusCode: http.StatusBadRequest,
			Code:       "VALIDATION_ERROR",
			Message:    message,
		},
		Field:      field,
		Constraint: constraint,
	}
}

// TimeoutError represents a timeout error.
type TimeoutError struct {
	*Error

	// Duration is the timeout duration.
	Duration time.Duration
}

// NewTimeoutError creates a new TimeoutError.
func NewTimeoutError(duration time.Duration) *TimeoutError {
	return &TimeoutError{
		Error: &Error{
			StatusCode: http.StatusGatewayTimeout,
			Code:       "TIMEOUT",
			Message:    fmt.Sprintf("operation timed out after %s", duration),
		},
		Duration: duration,
	}
}

// NetworkError represents a network-related error.
type NetworkError struct {
	*Error
}

// NewNetworkError creates a new NetworkError.
func NewNetworkError(err error) *NetworkError {
	return &NetworkError{
		Error: &Error{
			StatusCode: 0,
			Code:       "NETWORK_ERROR",
			Message:    "network error occurred",
			Err:        err,
		},
	}
}

// ServerError represents a server-side error.
type ServerError struct {
	*Error
}

// NewServerError creates a new ServerError.
func NewServerError(statusCode int, message string) *ServerError {
	return &ServerError{
		Error: &Error{
			StatusCode: statusCode,
			Code:       "SERVER_ERROR",
			Message:    message,
		},
	}
}

// IsAuthenticationError checks if the error is an authentication error.
func IsAuthenticationError(err error) bool {
	var authErr *AuthenticationError
	return errors.As(err, &authErr)
}

// IsRateLimitError checks if the error is a rate limit error.
func IsRateLimitError(err error) bool {
	var rateErr *RateLimitError
	return errors.As(err, &rateErr)
}

// IsValidationError checks if the error is a validation error.
func IsValidationError(err error) bool {
	var valErr *ValidationError
	return errors.As(err, &valErr)
}

// IsTimeoutError checks if the error is a timeout error.
func IsTimeoutError(err error) bool {
	var timeoutErr *TimeoutError
	return errors.As(err, &timeoutErr)
}

// IsNetworkError checks if the error is a network error.
func IsNetworkError(err error) bool {
	var netErr *NetworkError
	return errors.As(err, &netErr)
}

// IsServerError checks if the error is a server error.
func IsServerError(err error) bool {
	var serverErr *ServerError
	return errors.As(err, &serverErr)
}

// IsRetryable checks if the error is retryable.
func IsRetryable(err error) bool {
	var scErr *Error
	if errors.As(err, &scErr) {
		return scErr.IsRetryable()
	}

	var rateErr *RateLimitError
	if errors.As(err, &rateErr) {
		return true
	}

	var netErr *NetworkError
	if errors.As(err, &netErr) {
		return true
	}

	var timeoutErr *TimeoutError
	if errors.As(err, &timeoutErr) {
		return true
	}

	return false
}

// GetRetryAfter returns the retry-after duration for retryable errors.
func GetRetryAfter(err error) time.Duration {
	var rateErr *RateLimitError
	if errors.As(err, &rateErr) {
		return rateErr.RetryAfter
	}
	return 0
}

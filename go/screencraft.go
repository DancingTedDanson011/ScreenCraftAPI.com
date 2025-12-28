// Package screencraft provides a Go SDK for the ScreenCraft API.
//
// ScreenCraft enables screenshot capture and PDF generation from web pages
// with extensive customization options including viewport settings, delays,
// cookie consent handling, and more.
//
// Basic usage:
//
//	client := screencraft.New("your-api-key")
//
//	result, err := client.Screenshot(context.Background(), &screencraft.ScreenshotOptions{
//	    URL:      "https://example.com",
//	    Format:   screencraft.FormatPNG,
//	    FullPage: true,
//	})
//	if err != nil {
//	    log.Fatal(err)
//	}
//
//	os.WriteFile("screenshot.png", result.Data, 0644)
package screencraft

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"math/rand"
	"net/http"
	"strconv"
	"sync"
	"time"
)

const (
	// DefaultBaseURL is the default ScreenCraft API base URL.
	DefaultBaseURL = "https://screencraftapi.com/v1"

	// DefaultTimeout is the default HTTP client timeout.
	DefaultTimeout = 60 * time.Second

	// DefaultMaxRetries is the default maximum number of retries.
	DefaultMaxRetries = 3

	// DefaultRetryWaitMin is the default minimum retry wait time.
	DefaultRetryWaitMin = 1 * time.Second

	// DefaultRetryWaitMax is the default maximum retry wait time.
	DefaultRetryWaitMax = 30 * time.Second

	// Version is the SDK version.
	Version = "1.0.0"
)

// Client is the ScreenCraft API client.
type Client struct {
	// apiKey is the API key for authentication.
	apiKey string

	// baseURL is the API base URL.
	baseURL string

	// httpClient is the HTTP client used for requests.
	httpClient *http.Client

	// maxRetries is the maximum number of retries for failed requests.
	maxRetries int

	// retryWaitMin is the minimum time to wait between retries.
	retryWaitMin time.Duration

	// retryWaitMax is the maximum time to wait between retries.
	retryWaitMax time.Duration

	// userAgent is the User-Agent header value.
	userAgent string

	// debug enables debug logging.
	debug bool

	// logger is the logger for debug output.
	logger Logger

	// mu protects concurrent access to client fields.
	mu sync.RWMutex

	// lastRateLimit stores the last rate limit info received.
	lastRateLimit *RateLimitInfo
}

// Logger is the interface for logging.
type Logger interface {
	Printf(format string, v ...interface{})
}

// Option is a functional option for configuring the Client.
type Option func(*Client)

// New creates a new ScreenCraft client with the given API key.
func New(apiKey string, opts ...Option) *Client {
	c := &Client{
		apiKey:       apiKey,
		baseURL:      DefaultBaseURL,
		maxRetries:   DefaultMaxRetries,
		retryWaitMin: DefaultRetryWaitMin,
		retryWaitMax: DefaultRetryWaitMax,
		userAgent:    fmt.Sprintf("screencraft-go/%s", Version),
		httpClient: &http.Client{
			Timeout: DefaultTimeout,
		},
	}

	for _, opt := range opts {
		opt(c)
	}

	return c
}

// WithBaseURL sets a custom base URL for the API.
func WithBaseURL(url string) Option {
	return func(c *Client) {
		c.baseURL = url
	}
}

// WithHTTPClient sets a custom HTTP client.
func WithHTTPClient(client *http.Client) Option {
	return func(c *Client) {
		c.httpClient = client
	}
}

// WithTimeout sets the HTTP client timeout.
func WithTimeout(timeout time.Duration) Option {
	return func(c *Client) {
		c.httpClient.Timeout = timeout
	}
}

// WithMaxRetries sets the maximum number of retries for failed requests.
func WithMaxRetries(maxRetries int) Option {
	return func(c *Client) {
		c.maxRetries = maxRetries
	}
}

// WithRetryWait sets the minimum and maximum retry wait times.
func WithRetryWait(min, max time.Duration) Option {
	return func(c *Client) {
		c.retryWaitMin = min
		c.retryWaitMax = max
	}
}

// WithUserAgent sets a custom User-Agent header.
func WithUserAgent(userAgent string) Option {
	return func(c *Client) {
		c.userAgent = userAgent
	}
}

// WithDebug enables debug logging.
func WithDebug(debug bool) Option {
	return func(c *Client) {
		c.debug = debug
	}
}

// WithLogger sets a custom logger for debug output.
func WithLogger(logger Logger) Option {
	return func(c *Client) {
		c.logger = logger
	}
}

// SetAPIKey updates the API key.
func (c *Client) SetAPIKey(apiKey string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.apiKey = apiKey
}

// GetRateLimitInfo returns the last rate limit information received.
func (c *Client) GetRateLimitInfo() *RateLimitInfo {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.lastRateLimit
}

// doRequest performs an HTTP request with retries.
func (c *Client) doRequest(ctx context.Context, method, endpoint string, body interface{}) (*http.Response, error) {
	if c.apiKey == "" {
		return nil, ErrMissingAPIKey
	}

	var bodyReader io.Reader
	if body != nil {
		jsonBody, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("screencraft: failed to marshal request body: %w", err)
		}
		bodyReader = bytes.NewReader(jsonBody)
	}

	url := c.baseURL + endpoint

	var lastErr error
	for attempt := 0; attempt <= c.maxRetries; attempt++ {
		if attempt > 0 {
			waitTime := c.calculateBackoff(attempt, lastErr)
			c.logf("Retrying request (attempt %d/%d) after %s", attempt+1, c.maxRetries+1, waitTime)

			select {
			case <-ctx.Done():
				return nil, ctx.Err()
			case <-time.After(waitTime):
			}

			// Reset body reader for retry
			if body != nil {
				jsonBody, _ := json.Marshal(body)
				bodyReader = bytes.NewReader(jsonBody)
			}
		}

		req, err := http.NewRequestWithContext(ctx, method, url, bodyReader)
		if err != nil {
			return nil, fmt.Errorf("screencraft: failed to create request: %w", err)
		}

		req.Header.Set("Authorization", "Bearer "+c.apiKey)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Accept", "application/json, image/*, application/pdf")
		req.Header.Set("User-Agent", c.userAgent)

		c.logf("Making %s request to %s", method, url)

		resp, err := c.httpClient.Do(req)
		if err != nil {
			lastErr = NewNetworkError(err)
			if !IsRetryable(lastErr) || attempt == c.maxRetries {
				return nil, lastErr
			}
			continue
		}

		// Parse rate limit headers
		c.parseRateLimitHeaders(resp)

		// Check for errors
		if resp.StatusCode >= 400 {
			lastErr = c.parseErrorResponse(resp)
			if !IsRetryable(lastErr) || attempt == c.maxRetries {
				return nil, lastErr
			}
			resp.Body.Close()
			continue
		}

		return resp, nil
	}

	return nil, lastErr
}

// calculateBackoff calculates the backoff duration for a retry.
func (c *Client) calculateBackoff(attempt int, lastErr error) time.Duration {
	// Check for Retry-After from rate limit errors
	if retryAfter := GetRetryAfter(lastErr); retryAfter > 0 {
		return retryAfter
	}

	// Exponential backoff with jitter
	backoff := float64(c.retryWaitMin) * math.Pow(2, float64(attempt-1))
	if backoff > float64(c.retryWaitMax) {
		backoff = float64(c.retryWaitMax)
	}

	// Add jitter (up to 25%)
	jitter := backoff * 0.25 * rand.Float64()
	return time.Duration(backoff + jitter)
}

// parseRateLimitHeaders parses rate limit information from response headers.
func (c *Client) parseRateLimitHeaders(resp *http.Response) {
	c.mu.Lock()
	defer c.mu.Unlock()

	limit, _ := strconv.Atoi(resp.Header.Get("X-RateLimit-Limit"))
	remaining, _ := strconv.Atoi(resp.Header.Get("X-RateLimit-Remaining"))
	resetUnix, _ := strconv.ParseInt(resp.Header.Get("X-RateLimit-Reset"), 10, 64)

	if limit > 0 || remaining > 0 || resetUnix > 0 {
		c.lastRateLimit = &RateLimitInfo{
			Limit:     limit,
			Remaining: remaining,
			Reset:     time.Unix(resetUnix, 0),
		}
	}
}

// parseErrorResponse parses an error response from the API.
func (c *Client) parseErrorResponse(resp *http.Response) error {
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return &Error{
			StatusCode: resp.StatusCode,
			Message:    "failed to read error response",
			Err:        err,
		}
	}

	var apiResp APIResponse
	if err := json.Unmarshal(body, &apiResp); err != nil {
		return &Error{
			StatusCode: resp.StatusCode,
			Message:    string(body),
		}
	}

	baseErr := &Error{
		StatusCode: resp.StatusCode,
		RequestID:  resp.Header.Get("X-Request-ID"),
	}

	if apiResp.Error != nil {
		baseErr.Code = apiResp.Error.Code
		baseErr.Message = apiResp.Error.Message
		baseErr.Details = apiResp.Error.Details
	} else if apiResp.Message != "" {
		baseErr.Message = apiResp.Message
	}

	// Handle specific error types
	switch resp.StatusCode {
	case http.StatusUnauthorized:
		return &AuthenticationError{Error: baseErr}

	case http.StatusTooManyRequests:
		retryAfter := time.Duration(0)
		if ra := resp.Header.Get("Retry-After"); ra != "" {
			if seconds, err := strconv.Atoi(ra); err == nil {
				retryAfter = time.Duration(seconds) * time.Second
			}
		}

		limit, _ := strconv.Atoi(resp.Header.Get("X-RateLimit-Limit"))
		remaining, _ := strconv.Atoi(resp.Header.Get("X-RateLimit-Remaining"))
		resetUnix, _ := strconv.ParseInt(resp.Header.Get("X-RateLimit-Reset"), 10, 64)

		return &RateLimitError{
			Error:      baseErr,
			Limit:      limit,
			Remaining:  remaining,
			ResetAt:    time.Unix(resetUnix, 0),
			RetryAfter: retryAfter,
		}

	case http.StatusBadRequest:
		field := ""
		constraint := ""
		if apiResp.Error != nil && apiResp.Error.Details != nil {
			if f, ok := apiResp.Error.Details["field"].(string); ok {
				field = f
			}
			if c, ok := apiResp.Error.Details["constraint"].(string); ok {
				constraint = c
			}
		}
		return &ValidationError{
			Error:      baseErr,
			Field:      field,
			Constraint: constraint,
		}

	case http.StatusInternalServerError,
		http.StatusBadGateway,
		http.StatusServiceUnavailable,
		http.StatusGatewayTimeout:
		return &ServerError{Error: baseErr}
	}

	return baseErr
}

// logf logs a message if debug mode is enabled.
func (c *Client) logf(format string, v ...interface{}) {
	if c.debug && c.logger != nil {
		c.logger.Printf(format, v...)
	}
}

// ValidateScreenshotOptions validates screenshot options.
func ValidateScreenshotOptions(opts *ScreenshotOptions) error {
	if opts == nil {
		return ErrMissingURL
	}

	if opts.URL == "" {
		return ErrMissingURL
	}

	if opts.Quality < 0 || opts.Quality > 100 {
		return ErrInvalidQuality
	}

	if opts.Viewport != nil {
		if opts.Viewport.Width < 0 || opts.Viewport.Height < 0 {
			return ErrInvalidViewport
		}
	}

	return nil
}

// ValidatePDFOptions validates PDF options.
func ValidatePDFOptions(opts *PDFOptions) error {
	if opts == nil {
		return ErrMissingURL
	}

	if opts.URL == "" {
		return ErrMissingURL
	}

	if opts.Scale != 0 && (opts.Scale < 0.1 || opts.Scale > 2.0) {
		return NewValidationError("scale", "scale must be between 0.1 and 2.0", "range").Error
	}

	if opts.Viewport != nil {
		if opts.Viewport.Width < 0 || opts.Viewport.Height < 0 {
			return ErrInvalidViewport
		}
	}

	return nil
}

// Bool returns a pointer to the given bool value.
// Useful for setting optional boolean fields.
func Bool(v bool) *bool {
	return &v
}

// Int returns a pointer to the given int value.
// Useful for setting optional int fields.
func Int(v int) *int {
	return &v
}

// String returns a pointer to the given string value.
// Useful for setting optional string fields.
func String(v string) *string {
	return &v
}

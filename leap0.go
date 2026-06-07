// Package leap0 provides a Go SDK for the Leap0 cloud sandbox platform.
package leap0

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/leap0dev/leap0-go/internal/transport"
)

const (
	Version                = "0.1.0"
	DefaultBaseURL         = "https://api.leap0.dev"
	DefaultSandboxDomain   = "sandbox.leap0.dev"
	DefaultTemplate        = "system/debian:bookworm"
	DefaultCodeInterpreter = "system/code-interpreter:v0.1.0"
	DefaultDesktop         = "system/desktop:v0.1.0"
	DefaultVCPU            = 1
	DefaultMemoryMiB       = 1024
	DefaultTimeoutSec      = 300

	sdkSource = "sdk-go"
)

// Sentinel errors.
var (
	ErrNotFound  = errors.New("leap0: not found")
	ErrForbidden = errors.New("leap0: forbidden")
	ErrConflict  = errors.New("leap0: conflict")
	ErrRateLimit = errors.New("leap0: rate limited")
	ErrTimeout   = errors.New("leap0: timeout")
)

// APIError is returned on non-success HTTP responses.
type APIError struct {
	StatusCode int
	Headers    http.Header
	Body       string
	sentinel   error
}

func (e *APIError) Error() string {
	if e.Body != "" {
		return fmt.Sprintf("leap0: %s (status %d): %s", http.StatusText(e.StatusCode), e.StatusCode, e.Body)
	}
	return fmt.Sprintf("leap0: %s (status %d)", http.StatusText(e.StatusCode), e.StatusCode)
}
func (e *APIError) Is(target error) bool { return e.sentinel == target }
func (e *APIError) Unwrap() error        { return e.sentinel }
func (e *APIError) Retryable() bool      { return e.StatusCode == 429 || e.StatusCode == 503 }

func wrapErr(msg string, err error) error {
	if err == nil {
		return nil
	}
	var te *transport.APIError
	if errors.As(err, &te) {
		apiErr := &APIError{StatusCode: te.Status, Headers: te.Headers, Body: te.Body}
		switch te.Status {
		case 403:
			apiErr.sentinel = ErrForbidden
		case 404:
			apiErr.sentinel = ErrNotFound
		case 409:
			apiErr.sentinel = ErrConflict
		case 429:
			apiErr.sentinel = ErrRateLimit
		}
		return fmt.Errorf("%s: %w", msg, apiErr)
	}
	return fmt.Errorf("%s: %w", msg, err)
}

package agentstudio

import (
	"errors"
	"fmt"
	"net/http"
)

// Sentinel errors callers can match with errors.Is.
//
// These mirror the most common HTTP outcomes from the Agent Studio backend
// and let CLI commands give actionable hints without parsing status codes.
var (
	// ErrUnauthorized is returned when the API key is invalid or rejected
	// (HTTP 401).
	ErrUnauthorized = errors.New("agent studio: unauthorized — check your application id and api key")

	// ErrForbidden is returned for 403 responses that aren't a feature gate.
	// Typical cause: the API key is missing the required ACL.
	ErrForbidden = errors.New("agent studio: forbidden — the api key is missing a required ACL")

	// ErrFeatureDisabled is returned when the application does not have
	// Agent Studio enabled (the backend's gen_ai.agent_studio_enabled gate).
	// Detected from a 403 with a body matching the feature-disabled marker.
	ErrFeatureDisabled = errors.New("agent studio: feature is not enabled for this application — contact your Algolia account team or enable it from the Dashboard")

	// ErrNotFound is returned for 404 responses.
	ErrNotFound = errors.New("agent studio: resource not found")

	// ErrServer wraps any 5xx response.
	ErrServer = errors.New("agent studio: server error")
)

// APIError is the structured error returned for any non-2xx response.
//
// It always wraps one of the sentinel errors above (ErrUnauthorized,
// ErrForbidden, ErrFeatureDisabled, ErrNotFound, or ErrServer for 5xx;
// nil for other 4xx) so callers can match with errors.Is.
type APIError struct {
	StatusCode int
	// Detail is the parsed "detail" / "message" field from the response body
	// when present, otherwise the raw body truncated to 512 bytes.
	Detail string
	// Body is the raw response body, kept for debugging.
	Body []byte
	// Sentinel is the wrapped sentinel error (see above), or nil.
	Sentinel error
}

func (e *APIError) Error() string {
	if e.Detail != "" {
		return fmt.Sprintf("agent studio: %d %s: %s", e.StatusCode, http.StatusText(e.StatusCode), e.Detail)
	}
	return fmt.Sprintf("agent studio: %d %s", e.StatusCode, http.StatusText(e.StatusCode))
}

func (e *APIError) Unwrap() error { return e.Sentinel }

package agentstudio

import (
	"errors"
	"fmt"
	"net/http"
)

// Sentinel errors callers can match with errors.Is.
var (
	ErrUnauthorized    = errors.New("agent studio: unauthorized — check your application id and api key")
	ErrForbidden       = errors.New("agent studio: forbidden — the api key is missing a required ACL")
	ErrFeatureDisabled = errors.New(
		"agent studio: feature is not enabled for this application — contact your Algolia account team or enable it from the Dashboard",
	)
	ErrNotFound = errors.New("agent studio: resource not found")
	ErrServer   = errors.New("agent studio: server error")
)

// APIError is the structured error returned for any non-2xx response.
// Sentinel wraps one of the package-level sentinels (or nil for 4xx
// not in the table above).
type APIError struct {
	StatusCode int
	Detail     string
	Body       []byte
	Sentinel   error
}

func (e *APIError) Error() string {
	if e.Detail != "" {
		return fmt.Sprintf("agent studio: %d %s: %s", e.StatusCode, http.StatusText(e.StatusCode), e.Detail)
	}
	return fmt.Sprintf("agent studio: %d %s", e.StatusCode, http.StatusText(e.StatusCode))
}

func (e *APIError) Unwrap() error { return e.Sentinel }

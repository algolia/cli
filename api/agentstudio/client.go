package agentstudio

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// Standard Algolia auth headers used by Agent Studio.
//
// HeaderUserID is a CLEARTEXT label used by the backend for telemetry
// and rate-limiting only — it is NOT an authorization signal and MUST
// NOT be used for access decisions. The backend's signed equivalent
// (X-Algolia-Secure-User-Token, see common/models/secure_user_token.py
// in algolia/conversational-ai) is wired into the streaming
// /completions endpoint via CompletionOptions.SecureUserToken.
const (
	HeaderApplicationID = "X-Algolia-Application-Id"
	HeaderAPIKey        = "X-Algolia-API-Key" //nolint:gosec // header name, not a credential
	HeaderUserID        = "X-Algolia-User-ID"
)

// TODO(yuki): replace this hand-written client with the generated client
// once algolia/api-clients-automation publishes a Go module from
// algolia/conversational-ai's specs/agent-studio/spec.yml. Same pipeline
// that produces our Search SDK; tracked separately.

// Config configures a Client. All fields except ApplicationID, APIKey, and
// BaseURL are optional.
type Config struct {
	// BaseURL is the Agent Studio base URL without trailing slash and
	// without the /1 suffix (use ResolveHost to build it).
	BaseURL string

	// ApplicationID and APIKey are the standard Algolia credentials.
	// Required.
	ApplicationID string
	APIKey        string

	// UserID is sent as X-Algolia-User-ID. The backend defaults missing
	// values to "default" outside production but enforces presence in prod.
	// Recommended pattern from the CLI: "cli-<profile-name>".
	UserID string

	// UserAgent is sent as User-Agent. If empty, a minimal default is used.
	UserAgent string

	// HTTPClient overrides the default http.Client. Mainly for tests.
	HTTPClient *http.Client
}

// Client talks to the Agent Studio backend.
//
// Methods are organised by API tag (one source file per tag):
//
//   - agents.go        — Agents tag (CRUD, lifecycle, cache invalidation)
//   - completions.go   — Completions tag (streaming + buffered)
//   - providers.go     — Providers tag (CRUD + model discovery)
//   - configuration.go — Configurations tag (app-wide settings)
//
// This file only carries the cross-cutting pieces (Config, Client,
// NewClient, header injection, error mapping) so adding a new resource
// is a single new file plus a single new test file.
//
// All methods accept a context for cancellation and propagate request
// failures as *APIError (which wraps the appropriate sentinel error so
// errors.Is works).
type Client struct {
	cfg        Config
	httpClient *http.Client
}

// NewClient validates cfg and returns a ready-to-use Client.
func NewClient(cfg Config) (*Client, error) {
	if strings.TrimSpace(cfg.BaseURL) == "" {
		return nil, fmt.Errorf("agent studio: base url is required")
	}
	if strings.TrimSpace(cfg.ApplicationID) == "" {
		return nil, fmt.Errorf("agent studio: application id is required")
	}
	if strings.TrimSpace(cfg.APIKey) == "" {
		return nil, fmt.Errorf("agent studio: api key is required")
	}

	cfg.BaseURL = strings.TrimRight(cfg.BaseURL, "/")
	if cfg.HTTPClient == nil {
		cfg.HTTPClient = http.DefaultClient
	}
	if cfg.UserAgent == "" {
		cfg.UserAgent = "algolia-cli/agentstudio"
	}

	return &Client{cfg: cfg, httpClient: cfg.HTTPClient}, nil
}

func (c *Client) setHeaders(req *http.Request) {
	req.Header.Set(HeaderApplicationID, c.cfg.ApplicationID)
	req.Header.Set(HeaderAPIKey, c.cfg.APIKey)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", c.cfg.UserAgent)
	if c.cfg.UserID != "" {
		req.Header.Set(HeaderUserID, c.cfg.UserID)
	}
}

// checkResponse returns nil for 2xx and an *APIError otherwise.
func checkResponse(resp *http.Response) error {
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		return nil
	}

	body, _ := io.ReadAll(io.LimitReader(resp.Body, 64*1024))

	apiErr := &APIError{
		StatusCode: resp.StatusCode,
		Body:       body,
		Detail:     extractDetail(body),
		Sentinel:   sentinelFor(resp.StatusCode, body),
	}
	return apiErr
}

// extractDetail pulls a human-readable message from the response body.
//
// The backend returns errors in three observed shapes:
//   - FastAPI validation: {"detail":[{"msg":"...","loc":[...]}, ...]}
//     (often paired with a generic message like "Input is invalid, see
//     detail/body:" — we must prefer the structured detail).
//   - FastAPI default:    {"detail":"..."}
//   - Algolia ClientError: {"message":"..."}
//
// Priority is structured detail > string detail > message > raw body, so
// we never return a "see detail/body:" pointer when the actual detail is
// right there.
func extractDetail(body []byte) string {
	if len(body) == 0 {
		return ""
	}

	var generic map[string]any
	if err := json.Unmarshal(body, &generic); err != nil {
		s := strings.TrimSpace(string(body))
		if len(s) > 512 {
			return s[:512] + "…"
		}
		return s
	}

	switch d := generic["detail"].(type) {
	case []any:
		if len(d) > 0 {
			if first, ok := d[0].(map[string]any); ok {
				if msg, ok := first["msg"].(string); ok && msg != "" {
					return msg
				}
			}
		}
	case string:
		if d != "" {
			return d
		}
	}

	if msg, ok := generic["message"].(string); ok && msg != "" {
		return msg
	}

	return ""
}

// sentinelFor maps a status code (and body markers) to one of the package-level
// sentinel errors so callers can match with errors.Is.
func sentinelFor(status int, body []byte) error {
	switch {
	case status == http.StatusUnauthorized:
		return ErrUnauthorized
	case status == http.StatusForbidden:
		// The backend uses this exact phrase when the GenAI feature flag is
		// off for the app (see rag/dependencies/auth.py: "This feature is
		// not enabled for this application.").
		if strings.Contains(strings.ToLower(string(body)), "feature is not enabled") {
			return ErrFeatureDisabled
		}
		return ErrForbidden
	case status == http.StatusNotFound:
		return ErrNotFound
	case status >= 500:
		return ErrServer
	}
	return nil
}

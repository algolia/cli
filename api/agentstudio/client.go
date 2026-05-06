package agentstudio

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

// Standard Algolia auth headers used by Agent Studio.
//
// HeaderUserID is a CLEARTEXT label used by the backend for telemetry
// and rate-limiting only — it is NOT an authorization signal and MUST
// NOT be used for access decisions. The backend's signed equivalent
// (X-Algolia-Secure-User-Token, see common/models/secure_user_token.py
// in algolia/conversational-ai) is required by the streaming
// /completions endpoint and will be added in a later phase once the
// minting endpoint is reachable from the CLI's OAuth identity (Track A
// in tmp/agent_studio_plan.md).
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

// ListAgents calls GET /1/agents with optional pagination/filter params.
func (c *Client) ListAgents(ctx context.Context, params ListAgentsParams) (*PaginatedAgentsResponse, error) {
	q := url.Values{}
	if params.Page > 0 {
		q.Set("page", strconv.Itoa(params.Page))
	}
	if params.Limit > 0 {
		q.Set("limit", strconv.Itoa(params.Limit))
	}
	if params.ProviderID != "" {
		q.Set("providerId", params.ProviderID)
	}

	endpoint := c.cfg.BaseURL + "/1/agents"
	if encoded := q.Encode(); encoded != "" {
		endpoint += "?" + encoded
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, err
	}
	c.setHeaders(req)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("agent studio: list agents: %w", err)
	}
	defer resp.Body.Close()

	if err := checkResponse(resp); err != nil {
		return nil, err
	}

	var out PaginatedAgentsResponse
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return nil, fmt.Errorf("agent studio: decode list agents response: %w", err)
	}
	return &out, nil
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
// The backend returns errors in two main shapes:
//   - FastAPI default: {"detail": "..."} or {"detail": [{"msg":"..."}, ...]}
//   - common.exceptions.ClientError: {"message": "..."}
func extractDetail(body []byte) string {
	if len(body) == 0 {
		return ""
	}

	var generic map[string]any
	if err := json.Unmarshal(body, &generic); err != nil {
		// Not JSON; return the raw body trimmed.
		s := strings.TrimSpace(string(body))
		if len(s) > 512 {
			return s[:512] + "…"
		}
		return s
	}

	if msg, ok := generic["message"].(string); ok && msg != "" {
		return msg
	}

	switch d := generic["detail"].(type) {
	case string:
		return d
	case []any:
		if len(d) > 0 {
			if first, ok := d[0].(map[string]any); ok {
				if msg, ok := first["msg"].(string); ok {
					return msg
				}
			}
		}
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

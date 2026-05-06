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

// CreateAgent calls POST /1/agents with the supplied request body.
//
// The body is sent as opaque JSON: the AgentConfigCreate schema in
// algolia/conversational-ai is large, deeply validated, and evolves
// frequently (8+ fields, nested ToolConfig, free-form config dict). The
// CLI is a pass-through — it lets users supply the JSON, the backend
// validates, our 422 surfacing makes errors actionable. Mirroring the
// schema in Go would lie about parity and force a release every time the
// backend adds a field.
func (c *Client) CreateAgent(ctx context.Context, body json.RawMessage) (*Agent, error) {
	if len(body) == 0 {
		return nil, fmt.Errorf("agent studio: create agent: body is required")
	}
	return c.doAgentMutation(ctx, http.MethodPost, c.cfg.BaseURL+"/1/agents", body, "create agent")
}

// UpdateAgent calls PATCH /1/agents/{id} with the supplied partial body.
// See CreateAgent for the rationale behind json.RawMessage.
func (c *Client) UpdateAgent(ctx context.Context, id string, body json.RawMessage) (*Agent, error) {
	if strings.TrimSpace(id) == "" {
		return nil, fmt.Errorf("agent studio: agent id is required")
	}
	if len(body) == 0 {
		return nil, fmt.Errorf("agent studio: update agent: body is required")
	}
	endpoint := c.cfg.BaseURL + "/1/agents/" + url.PathEscape(id)
	return c.doAgentMutation(ctx, http.MethodPatch, endpoint, body, "update agent")
}

// DeleteAgent calls DELETE /1/agents/{id}.
//
// Returns nil on the backend's HTTP 204 No Content. The backend
// soft-deletes (the row stays in the DB with a deleted flag), so this is
// reversible at the platform level — but the CLI exposes it as a
// terminal action; recovery is a backend-side ops concern.
func (c *Client) DeleteAgent(ctx context.Context, id string) error {
	if strings.TrimSpace(id) == "" {
		return fmt.Errorf("agent studio: agent id is required")
	}
	endpoint := c.cfg.BaseURL + "/1/agents/" + url.PathEscape(id)

	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, endpoint, nil)
	if err != nil {
		return err
	}
	c.setHeaders(req)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("agent studio: delete agent: %w", err)
	}
	defer resp.Body.Close()

	return checkResponse(resp)
}

// PublishAgent calls POST /1/agents/{id}/publish. Backend transitions the
// agent's status from draft to published; returns the updated Agent.
func (c *Client) PublishAgent(ctx context.Context, id string) (*Agent, error) {
	return c.doAgentLifecycle(ctx, id, "publish")
}

// UnpublishAgent calls POST /1/agents/{id}/unpublish.
func (c *Client) UnpublishAgent(ctx context.Context, id string) (*Agent, error) {
	return c.doAgentLifecycle(ctx, id, "unpublish")
}

// DuplicateAgent calls POST /1/agents/{id}/duplicate. Returns the newly
// created Agent (which has its own ID).
func (c *Client) DuplicateAgent(ctx context.Context, id string) (*Agent, error) {
	return c.doAgentLifecycle(ctx, id, "duplicate")
}

// doAgentLifecycle is the shared implementation for publish/unpublish/duplicate.
// All three are POST /1/agents/{id}/<verb> with no body, returning Agent.
func (c *Client) doAgentLifecycle(ctx context.Context, id, verb string) (*Agent, error) {
	if strings.TrimSpace(id) == "" {
		return nil, fmt.Errorf("agent studio: agent id is required")
	}
	endpoint := c.cfg.BaseURL + "/1/agents/" + url.PathEscape(id) + "/" + verb
	return c.doAgentMutation(ctx, http.MethodPost, endpoint, nil, verb+" agent")
}

// doAgentMutation issues an HTTP request with optional JSON body, expects
// a JSON Agent response, and is shared by Create/Update/Publish/Unpublish/Duplicate.
// errLabel is used to scope error messages (e.g., "create agent").
func (c *Client) doAgentMutation(
	ctx context.Context,
	method, endpoint string,
	body json.RawMessage,
	errLabel string,
) (*Agent, error) {
	var reqBody io.Reader
	if body != nil {
		reqBody = strings.NewReader(string(body))
	}

	req, err := http.NewRequestWithContext(ctx, method, endpoint, reqBody)
	if err != nil {
		return nil, err
	}
	c.setHeaders(req)
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("agent studio: %s: %w", errLabel, err)
	}
	defer resp.Body.Close()

	if err := checkResponse(resp); err != nil {
		return nil, err
	}

	var out Agent
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return nil, fmt.Errorf("agent studio: decode %s response: %w", errLabel, err)
	}
	return &out, nil
}

// GetAgent calls GET /1/agents/{id}.
func (c *Client) GetAgent(ctx context.Context, id string) (*Agent, error) {
	if strings.TrimSpace(id) == "" {
		return nil, fmt.Errorf("agent studio: agent id is required")
	}

	endpoint := c.cfg.BaseURL + "/1/agents/" + url.PathEscape(id)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, err
	}
	c.setHeaders(req)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("agent studio: get agent: %w", err)
	}
	defer resp.Body.Close()

	if err := checkResponse(resp); err != nil {
		return nil, err
	}

	var out Agent
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return nil, fmt.Errorf("agent studio: decode get agent response: %w", err)
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

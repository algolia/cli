package agentstudio

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

const (
	HeaderApplicationID = "X-Algolia-Application-Id"
	HeaderAPIKey        = "X-Algolia-API-Key" //nolint:gosec // header name, not a credential
	HeaderUserID        = "X-Algolia-User-ID"
)

// Config configures a Client. ApplicationID, APIKey, and BaseURL are
// required; everything else is optional.
type Config struct {
	BaseURL       string
	ApplicationID string
	APIKey        string
	UserID        string
	UserAgent     string
	HTTPClient    *http.Client
}

// Client talks to the Agent Studio backend. Methods are organised by API
// tag — one source file per tag (agents.go, completions.go, …). This
// file carries only Config, NewClient, header injection, and error
// mapping.
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

	return &APIError{
		StatusCode: resp.StatusCode,
		Body:       body,
		Detail:     extractDetail(body),
		Sentinel:   sentinelFor(resp.StatusCode, body),
	}
}

// extractDetail pulls a human-readable message from the response body.
// Priority: structured FastAPI detail[].msg > string detail > Algolia
// {message:...} > raw body.
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

// sentinelFor maps a status code (and body markers) to a sentinel error.
func sentinelFor(status int, body []byte) error {
	switch {
	case status == http.StatusUnauthorized:
		return ErrUnauthorized
	case status == http.StatusForbidden:
		// Backend uses this exact phrase when the GenAI feature flag is off.
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

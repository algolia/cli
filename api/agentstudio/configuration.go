package agentstudio

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

// GetConfiguration calls GET /1/configuration. ACL is `logs` (not
// `settings`) because the only field today governs log retention.
func (c *Client) GetConfiguration(ctx context.Context) (*ApplicationConfig, error) {
	endpoint := c.cfg.BaseURL + "/1/configuration"

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, err
	}
	c.setHeaders(req)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("agent studio: get configuration: %w", err)
	}
	defer resp.Body.Close()

	if err := checkResponse(resp); err != nil {
		return nil, err
	}

	var out ApplicationConfig
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return nil, fmt.Errorf("agent studio: decode get configuration response: %w", err)
	}
	return &out, nil
}

// UpdateConfiguration calls PATCH /1/configuration. Body is opaque JSON
// for symmetry with the rest of the agents tree.
func (c *Client) UpdateConfiguration(ctx context.Context, body json.RawMessage) (*ApplicationConfig, error) {
	if len(body) == 0 {
		return nil, fmt.Errorf("agent studio: update configuration: body is required")
	}

	endpoint := c.cfg.BaseURL + "/1/configuration"

	req, err := http.NewRequestWithContext(ctx, http.MethodPatch, endpoint, strings.NewReader(string(body)))
	if err != nil {
		return nil, err
	}
	c.setHeaders(req)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("agent studio: update configuration: %w", err)
	}
	defer resp.Body.Close()

	if err := checkResponse(resp); err != nil {
		return nil, err
	}

	var out ApplicationConfig
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return nil, fmt.Errorf("agent studio: decode update configuration response: %w", err)
	}
	return &out, nil
}

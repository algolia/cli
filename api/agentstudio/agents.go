package agentstudio

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

// The standard agent CRUD/lifecycle surface now lives in the official SDK
// (algoliasearch-client-go/.../agent-studio). The methods kept here back
// features the SDK can't serve cleanly: DuplicateAgent (no SDK endpoint), and
// GetAgent + UpdateAgent (used by `agents tools add-search-index`, which
// merges a tool entry into the agent's raw tools array and PATCHes it
// verbatim — the SDK's typed AgentConfigUpdate.Tools union can't round-trip
// that arbitrary tool JSON).

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

// UpdateAgent calls PATCH /1/agents/{id}.
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

// DuplicateAgent calls POST /1/agents/{id}/duplicate.
func (c *Client) DuplicateAgent(ctx context.Context, id string) (*Agent, error) {
	return c.doAgentLifecycle(ctx, id, "duplicate")
}

func (c *Client) doAgentLifecycle(ctx context.Context, id, verb string) (*Agent, error) {
	if strings.TrimSpace(id) == "" {
		return nil, fmt.Errorf("agent studio: agent id is required")
	}
	endpoint := c.cfg.BaseURL + "/1/agents/" + url.PathEscape(id) + "/" + verb
	return c.doAgentMutation(ctx, http.MethodPost, endpoint, nil, verb+" agent")
}

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

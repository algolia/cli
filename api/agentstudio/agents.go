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

// ListAgents calls GET /1/agents.
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

// CreateAgent calls POST /1/agents. Body is opaque JSON; see docs/agents.md.
func (c *Client) CreateAgent(ctx context.Context, body json.RawMessage) (*Agent, error) {
	if len(body) == 0 {
		return nil, fmt.Errorf("agent studio: create agent: body is required")
	}
	return c.doAgentMutation(ctx, http.MethodPost, c.cfg.BaseURL+"/1/agents", body, "create agent")
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

// DeleteAgent calls DELETE /1/agents/{id}. Backend soft-deletes;
// recovery is a backend ops concern.
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

// PublishAgent calls POST /1/agents/{id}/publish.
func (c *Client) PublishAgent(ctx context.Context, id string) (*Agent, error) {
	return c.doAgentLifecycle(ctx, id, "publish")
}

// UnpublishAgent calls POST /1/agents/{id}/unpublish.
func (c *Client) UnpublishAgent(ctx context.Context, id string) (*Agent, error) {
	return c.doAgentLifecycle(ctx, id, "unpublish")
}

// DuplicateAgent calls POST /1/agents/{id}/duplicate.
func (c *Client) DuplicateAgent(ctx context.Context, id string) (*Agent, error) {
	return c.doAgentLifecycle(ctx, id, "duplicate")
}

// InvalidateAgentCache calls DELETE /1/agents/{id}/cache?before=YYYY-MM-DD.
// Empty before = wipe all cache entries. Date format is validated
// server-side (see docs/agents.md gotchas).
func (c *Client) InvalidateAgentCache(ctx context.Context, id, before string) error {
	if strings.TrimSpace(id) == "" {
		return fmt.Errorf("agent studio: agent id is required")
	}

	endpoint := c.cfg.BaseURL + "/1/agents/" + url.PathEscape(id) + "/cache"
	if before != "" {
		q := url.Values{}
		q.Set("before", before)
		endpoint += "?" + q.Encode()
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, endpoint, nil)
	if err != nil {
		return err
	}
	c.setHeaders(req)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("agent studio: invalidate agent cache: %w", err)
	}
	defer resp.Body.Close()

	return checkResponse(resp)
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

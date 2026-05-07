// Hidden/internal endpoints (x-hidden in the OpenAPI spec). Surfaced
// in the CLI behind hidden commands for diagnostics. Memory ops use
// the doubled /1/agents/agents/ path — see docs/agents.md gotchas.

package agentstudio

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

// GetStatus calls GET /status. No auth headers needed; useful as a
// liveness probe and to read the build version + migration revision.
func (c *Client) GetStatus(ctx context.Context) (StatusResponse, error) {
	endpoint := c.cfg.BaseURL + "/status"
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, err
	}
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("agent studio: get status: %w", err)
	}
	defer resp.Body.Close()
	if err := checkResponse(resp); err != nil {
		return nil, err
	}
	var out StatusResponse
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return nil, fmt.Errorf("agent studio: decode status response: %w", err)
	}
	return out, nil
}

// GetProviderModelDefaults calls GET /1/providers/models/defaults and
// returns provider-type → recommended-model-name.
func (c *Client) GetProviderModelDefaults(ctx context.Context) (ModelDefaults, error) {
	endpoint := c.cfg.BaseURL + "/1/providers/models/defaults"
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, err
	}
	c.setHeaders(req)
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("agent studio: get model defaults: %w", err)
	}
	defer resp.Body.Close()
	if err := checkResponse(resp); err != nil {
		return nil, err
	}
	var out ModelDefaults
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return nil, fmt.Errorf("agent studio: decode model defaults: %w", err)
	}
	return out, nil
}

// AgentMemorize calls POST /1/agents/agents/{id}/memorize (doubled path).
func (c *Client) AgentMemorize(ctx context.Context, agentID string, body json.RawMessage) (json.RawMessage, error) {
	return c.doAgentMemoryOp(ctx, agentID, "memorize", body)
}

func (c *Client) AgentPonder(ctx context.Context, agentID string, body json.RawMessage) (json.RawMessage, error) {
	return c.doAgentMemoryOp(ctx, agentID, "ponder", body)
}

func (c *Client) AgentConsolidate(ctx context.Context, agentID string, body json.RawMessage) (json.RawMessage, error) {
	return c.doAgentMemoryOp(ctx, agentID, "consolidate", body)
}

func (c *Client) doAgentMemoryOp(
	ctx context.Context,
	agentID, verb string,
	body json.RawMessage,
) (json.RawMessage, error) {
	if strings.TrimSpace(agentID) == "" {
		return nil, fmt.Errorf("agent studio: %s: agent id is required", verb)
	}
	if len(body) == 0 {
		return nil, fmt.Errorf("agent studio: %s: body is required", verb)
	}
	endpoint := c.cfg.BaseURL + "/1/agents/agents/" + url.PathEscape(agentID) + "/" + verb
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	c.setHeaders(req)
	req.Header.Set("Content-Type", "application/json")
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("agent studio: %s: %w", verb, err)
	}
	defer resp.Body.Close()
	if err := checkResponse(resp); err != nil {
		return nil, err
	}
	out, err := readAll(resp)
	if err != nil {
		return nil, fmt.Errorf("agent studio: read %s response: %w", verb, err)
	}
	return out, nil
}

func readAll(resp *http.Response) ([]byte, error) {
	var buf bytes.Buffer
	if _, err := buf.ReadFrom(resp.Body); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

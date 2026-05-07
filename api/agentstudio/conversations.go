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

// ListConversations calls GET /1/agents/{agent_id}/conversations.
// All query params are optional; empty values are omitted from the wire.
func (c *Client) ListConversations(
	ctx context.Context,
	agentID string,
	params ListConversationsParams,
) (*PaginatedConversationsResponse, error) {
	if strings.TrimSpace(agentID) == "" {
		return nil, fmt.Errorf("agent studio: list conversations: agent id is required")
	}

	q := url.Values{}
	if params.Page > 0 {
		q.Set("page", strconv.Itoa(params.Page))
	}
	if params.Limit > 0 {
		q.Set("limit", strconv.Itoa(params.Limit))
	}
	if params.StartDate != "" {
		q.Set("startDate", params.StartDate)
	}
	if params.EndDate != "" {
		q.Set("endDate", params.EndDate)
	}
	if params.IncludeFeedback {
		q.Set("includeFeedback", "true")
	}
	if params.FeedbackVote != nil {
		q.Set("feedbackVote", strconv.Itoa(*params.FeedbackVote))
	}

	endpoint := c.cfg.BaseURL + "/1/agents/" + url.PathEscape(agentID) + "/conversations"
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
		return nil, fmt.Errorf("agent studio: list conversations: %w", err)
	}
	defer resp.Body.Close()

	if err := checkResponse(resp); err != nil {
		return nil, err
	}

	var out PaginatedConversationsResponse
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return nil, fmt.Errorf("agent studio: decode list conversations response: %w", err)
	}
	return &out, nil
}

// GetConversation calls GET /1/agents/{agent_id}/conversations/{conversation_id}.
// Returns raw JSON — `messages` is a discriminated role union; see docs/agents.md.
func (c *Client) GetConversation(
	ctx context.Context,
	agentID, conversationID string,
	includeFeedback bool,
) (json.RawMessage, error) {
	if strings.TrimSpace(agentID) == "" {
		return nil, fmt.Errorf("agent studio: get conversation: agent id is required")
	}
	if strings.TrimSpace(conversationID) == "" {
		return nil, fmt.Errorf("agent studio: get conversation: conversation id is required")
	}

	endpoint := c.cfg.BaseURL + "/1/agents/" + url.PathEscape(agentID) +
		"/conversations/" + url.PathEscape(conversationID)
	if includeFeedback {
		endpoint += "?includeFeedback=true"
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, err
	}
	c.setHeaders(req)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("agent studio: get conversation: %w", err)
	}
	defer resp.Body.Close()

	if err := checkResponse(resp); err != nil {
		return nil, err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("agent studio: read get conversation response: %w", err)
	}
	return json.RawMessage(body), nil
}

// DeleteConversation calls DELETE /1/agents/{agent_id}/conversations/{conversation_id}.
// Returns nil on the backend's HTTP 204.
func (c *Client) DeleteConversation(ctx context.Context, agentID, conversationID string) error {
	if strings.TrimSpace(agentID) == "" {
		return fmt.Errorf("agent studio: delete conversation: agent id is required")
	}
	if strings.TrimSpace(conversationID) == "" {
		return fmt.Errorf("agent studio: delete conversation: conversation id is required")
	}

	endpoint := c.cfg.BaseURL + "/1/agents/" + url.PathEscape(agentID) +
		"/conversations/" + url.PathEscape(conversationID)

	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, endpoint, nil)
	if err != nil {
		return err
	}
	c.setHeaders(req)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("agent studio: delete conversation: %w", err)
	}
	defer resp.Body.Close()

	return checkResponse(resp)
}

// PurgeConversations calls DELETE /1/agents/{agent_id}/conversations.
// Backend rejects dateless purges; CLI enforces this at the flag layer.
// See docs/agents.md gotchas.
func (c *Client) PurgeConversations(ctx context.Context, agentID string, params PurgeConversationsParams) error {
	if strings.TrimSpace(agentID) == "" {
		return fmt.Errorf("agent studio: purge conversations: agent id is required")
	}

	endpoint := c.cfg.BaseURL + "/1/agents/" + url.PathEscape(agentID) + "/conversations"
	q := url.Values{}
	if params.StartDate != "" {
		q.Set("startDate", params.StartDate)
	}
	if params.EndDate != "" {
		q.Set("endDate", params.EndDate)
	}
	if encoded := q.Encode(); encoded != "" {
		endpoint += "?" + encoded
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, endpoint, nil)
	if err != nil {
		return err
	}
	c.setHeaders(req)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("agent studio: purge conversations: %w", err)
	}
	defer resp.Body.Close()

	return checkResponse(resp)
}

// ExportConversations calls GET /1/agents/{agent_id}/conversations/export.
// Returns raw JSON — the spec leaves the response body unspecified.
func (c *Client) ExportConversations(
	ctx context.Context,
	agentID string,
	params ExportConversationsParams,
) (json.RawMessage, error) {
	if strings.TrimSpace(agentID) == "" {
		return nil, fmt.Errorf("agent studio: export conversations: agent id is required")
	}

	endpoint := c.cfg.BaseURL + "/1/agents/" + url.PathEscape(agentID) + "/conversations/export"
	q := url.Values{}
	if params.StartDate != "" {
		q.Set("startDate", params.StartDate)
	}
	if params.EndDate != "" {
		q.Set("endDate", params.EndDate)
	}
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
		return nil, fmt.Errorf("agent studio: export conversations: %w", err)
	}
	defer resp.Body.Close()

	if err := checkResponse(resp); err != nil {
		return nil, err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("agent studio: read export response: %w", err)
	}
	return json.RawMessage(body), nil
}

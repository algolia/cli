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
//
// All query params are optional. Empty values are omitted from the wire
// (matches the backend's "use default" semantics for page/limit and
// "no filter" for date/feedback knobs).
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
//
// Returns the body as raw JSON because ConversationFullResponse embeds
// `messages: []MessageResponse-Output` which is a discriminated union
// over message roles (system/user/assistant/tool) with per-role nested
// content arrays. Same passthrough rationale as Agent.Config / Provider.Input
// — the CLI prints, the user's `jq` extracts.
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
//
// Backend behaviour:
//   - Both StartDate and EndDate empty → ALL conversations for this
//     agent are deleted. This is intentional but destructive enough
//     that the CLI requires an explicit `--all` flag (enforced at
//     the cmd layer, not here — the client mirrors the wire shape).
//   - Either present → range filter applied (inclusive on both ends,
//     per backend convention).
//
// Returns nil on the backend's HTTP 204.
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

// ExportConversations calls GET /1/agents/{agent_id}/conversations/export
// and returns the body as raw JSON.
//
// The OpenAPI spec leaves the response body unspecified (the operation
// has a 200 with no schema). Empirically the backend returns a JSON
// document; the CLI prints it as-is and lets users pipe through `jq`
// or write to a file with `--output-file`. Pinning a Go type here would
// silently break the day the backend switches to NDJSON or CSV.
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

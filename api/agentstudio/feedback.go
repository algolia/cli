package agentstudio

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

func (c *Client) CreateFeedback(ctx context.Context, body FeedbackCreate) (*Feedback, error) {
	if strings.TrimSpace(body.MessageID) == "" {
		return nil, fmt.Errorf("agent studio: create feedback: messageId is required")
	}
	if strings.TrimSpace(body.AgentID) == "" {
		return nil, fmt.Errorf("agent studio: create feedback: agentId is required")
	}
	if body.Vote != 0 && body.Vote != 1 {
		return nil, fmt.Errorf("agent studio: create feedback: vote must be 0 or 1")
	}
	raw, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.cfg.BaseURL+"/1/feedback", bytes.NewReader(raw))
	if err != nil {
		return nil, err
	}
	c.setHeaders(req)
	req.Header.Set("Content-Type", "application/json")
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("agent studio: create feedback: %w", err)
	}
	defer resp.Body.Close()
	if err := checkResponse(resp); err != nil {
		return nil, err
	}
	var out Feedback
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return nil, fmt.Errorf("agent studio: decode create feedback response: %w", err)
	}
	return &out, nil
}

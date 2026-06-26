package agentstudio

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

// CompletionOptions configures Completions(...) query params and headers.
// No* fields are inverted from the wire (backend defaults all three to
// true; only the negative case is interesting). See docs/agents.md
// "Completion runtime knobs".
type CompletionOptions struct {
	Stream          bool
	Compatibility   CompatibilityMode
	NoCache         bool
	NoMemory        bool
	NoAnalytics     bool
	SecureUserToken string
}

// Completions calls POST /1/agents/{agentID}/completions and returns the
// raw HTTP response. Caller closes resp.Body on all paths and inspects
// Content-Type to decide between ParseStream (text/event-stream) and a
// single json.Decode. agentID may be a UUID or the literal "test".
//
// Callers should pass a ctx that supports cancellation; for streaming,
// also consider an application-level deadline if the remote end stalls
// mid-body (the default HTTP client sets ResponseHeaderTimeout but does
// not cap total stream duration).
func (c *Client) Completions(
	ctx context.Context,
	agentID string,
	body json.RawMessage,
	opts CompletionOptions,
) (*http.Response, error) {
	if strings.TrimSpace(agentID) == "" {
		return nil, fmt.Errorf("agent studio: completions: agent id is required (or use \"test\")")
	}
	if len(body) == 0 {
		return nil, fmt.Errorf("agent studio: completions: body is required")
	}
	if !json.Valid(body) {
		return nil, fmt.Errorf("agent studio: completions: body is not valid JSON")
	}

	mode := opts.Compatibility
	if mode == "" {
		mode = CompatV5
	}

	q := url.Values{}
	q.Set("stream", boolToWire(opts.Stream))
	q.Set("compatibilityMode", string(mode))
	if opts.NoCache {
		q.Set("cache", "false")
	}
	if opts.NoMemory {
		q.Set("memory", "false")
	}
	if opts.NoAnalytics {
		q.Set("analytics", "false")
	}

	endpoint := c.cfg.BaseURL + "/1/agents/" + url.PathEscape(agentID) + "/completions?" + q.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, strings.NewReader(string(body)))
	if err != nil {
		return nil, err
	}
	c.setHeaders(req)
	req.Header.Set("Content-Type", "application/json")
	if opts.SecureUserToken != "" {
		req.Header.Set("X-Algolia-Secure-User-Token", opts.SecureUserToken)
	}
	req.Header.Set("Accept", "text/event-stream, application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("agent studio: completions: %w", err)
	}

	if err := checkResponse(resp); err != nil {
		// checkResponse only drains a 64 KiB prefix; close to release the conn.
		_ = resp.Body.Close()
		return nil, err
	}
	return resp, nil
}

func boolToWire(b bool) string {
	if b {
		return "true"
	}
	return "false"
}

package agentstudio

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

// CompletionOptions configures Completions(...) query parameters and
// per-request headers.
//
// Stream maps to ?stream=true|false; the default zero value (false) gives
// a buffered single-JSON response. Set explicitly via the command layer
// (`agents try` / `agents run` set Stream=true unless --no-stream).
//
// Compatibility maps to ?compatibilityMode=ai-sdk-4|ai-sdk-5. The backend
// requires this query param (no server-side default), so empty here is
// promoted to CompatV5 — its frames are standard SSE with [DONE], easier
// to parse defensively than v4's `<type>:<json>\n` line format.
//
// NoCache, NoMemory, and NoAnalytics are inverted from the backend's
// query-param polarity for two reasons:
//
//   - The backend defaults all three to true; only the negated case is
//     interesting from the CLI surface.
//   - The flag layer ships them as `--no-cache`/`--no-memory`/`--no-analytics`
//     so the option fields keep that polarity end-to-end.
//
// When a No*-field is false (the zero value) the corresponding query
// param is omitted entirely, which matches the backend's "default ON"
// behavior. The `memory` schema in particular is `anyOf [{const: false},
// {type: null}]` — false is the ONLY valid passable value, so always
// emitting `memory=true` would be a server-side validation error.
//
// SecureUserToken populates the X-Algolia-Secure-User-Token header when
// non-empty. It carries a signed JWT that scopes the conversation /
// memory / analytics partition to a specific end-user; required by the
// backend whenever a feature behind SecureUserTokenDep is enabled (see
// rag/dependencies/secure_user_token.py in algolia/conversational-ai).
// Empty here means no header is sent — the existing X-Algolia-User-ID
// fallback applies.
type CompletionOptions struct {
	Stream          bool
	Compatibility   CompatibilityMode
	NoCache         bool
	NoMemory        bool
	NoAnalytics     bool
	SecureUserToken string
}

// Completions calls POST /1/agents/{agentID}/completions and returns the
// raw HTTP response. The caller is responsible for:
//
//   - Closing resp.Body in all paths.
//   - Inspecting resp.Header.Get("Content-Type") to decide whether to
//     stream-parse via ParseStream (Content-Type: text/event-stream) or
//     to json.Decode the body once (any other Content-Type).
//
// agentID is either a real UUID or the literal string "test" — the
// backend special-cases "test" to mean "no agent is persisted; use the
// AgentTestConfiguration in the request body" (rag/routers/v1/
// agents_completion.py uses Union[uuid.UUID, Literal["test"]]).
//
// body must be a valid AgentCompletionRequest JSON document. The CLI is
// a pass-through (same rationale as CreateAgent/UpdateAgent): the
// request schema includes a discriminated `messages` union and a
// vendored `algolia.searchParameters` shape that evolves often. Server
// 422s surface the structured FastAPI detail via extractDetail.
//
// Cancellation is the caller's job: pass a ctx that is cancelled on
// SIGINT and the underlying transport will tear down the request mid-
// stream cleanly.
//
// This is the sole user-facing endpoint of the Completions API tag —
// both the buffered and streaming variants share the same handler;
// `?stream=` is what flips the response shape.
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
	// Only emit the negative cases — backend defaults match the omitted
	// state, so adding `cache=true`/`analytics=true` would be wire noise,
	// and `memory=true` would actually be a 422 (the schema only allows
	// `false` or null). See CompletionOptions godoc for the full reasoning.
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
	// Preferred Accept: streaming responses come back as text/event-stream
	// (both v4 and v5); buffered ones as application/json. Listing both
	// is safe — the server picks based on ?stream and we inspect the
	// Content-Type on return.
	req.Header.Set("Accept", "text/event-stream, application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("agent studio: completions: %w", err)
	}

	if err := checkResponse(resp); err != nil {
		// checkResponse only drains a 64 KiB prefix for the error detail;
		// it does NOT close the body. We have to do it here to release
		// the underlying connection back to the transport pool.
		_ = resp.Body.Close()
		return nil, err
	}
	return resp, nil
}

// boolToWire renders Go bools as the lowercase strings the FastAPI
// Query() bool coercion expects (it accepts case-insensitively but the
// canonical form is lowercase). Lives here because Completions is the
// only caller — if a future endpoint needs the same coercion, lift to
// client.go.
func boolToWire(b bool) string {
	if b {
		return "true"
	}
	return "false"
}

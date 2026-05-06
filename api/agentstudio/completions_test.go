package agentstudio

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCompletions_StreamingV5_DefaultMode(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/1/agents/test/completions", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, "true", r.URL.Query().Get("stream"))
		// Empty CompatibilityMode in opts must be promoted to ai-sdk-5.
		assert.Equal(t, "ai-sdk-5", r.URL.Query().Get("compatibilityMode"))
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))
		assert.Contains(t, r.Header.Get("Accept"), "text/event-stream")

		body, _ := io.ReadAll(r.Body)
		assert.JSONEq(t, `{"messages":[{"role":"user","content":"hi"}]}`, string(body))

		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("x-vercel-ai-ui-message-stream", "v1")
		_, _ = w.Write([]byte(`data: {"type":"text-delta","delta":"hello"}` + "\n\n"))
		_, _ = w.Write([]byte(`data: [DONE]` + "\n\n"))
	})

	_, c := newTestClient(t, mux)

	resp, err := c.Completions(context.Background(), "test",
		json.RawMessage(`{"messages":[{"role":"user","content":"hi"}]}`),
		CompletionOptions{Stream: true})
	require.NoError(t, err)
	t.Cleanup(func() { _ = resp.Body.Close() })

	assert.Equal(t, "text/event-stream", resp.Header.Get("Content-Type"))

	var got []StreamEvent
	require.NoError(t, ParseStream(resp.Body, func(e StreamEvent) error {
		got = append(got, e)
		return nil
	}))
	require.Len(t, got, 1)
	assert.Equal(t, "text-delta", got[0].Type)
}

func TestCompletions_BufferedNonStream(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/1/agents/abc-123/completions", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "false", r.URL.Query().Get("stream"))
		assert.Equal(t, "ai-sdk-4", r.URL.Query().Get("compatibilityMode"))

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"role":"assistant","content":"hello"}`))
	})

	_, c := newTestClient(t, mux)

	resp, err := c.Completions(context.Background(), "abc-123",
		json.RawMessage(`{"messages":[{"role":"user","content":"x"}]}`),
		CompletionOptions{Stream: false, Compatibility: CompatV4})
	require.NoError(t, err)
	t.Cleanup(func() { _ = resp.Body.Close() })

	assert.Equal(t, "application/json", resp.Header.Get("Content-Type"))
	body, _ := io.ReadAll(resp.Body)
	assert.JSONEq(t, `{"role":"assistant","content":"hello"}`, string(body))
}

func TestCompletions_RejectsEmptyAgentID(t *testing.T) {
	_, c := newTestClient(t, http.NewServeMux())
	_, err := c.Completions(context.Background(), "  ",
		json.RawMessage(`{}`), CompletionOptions{})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "agent id is required")
}

func TestCompletions_RejectsEmptyOrInvalidBody(t *testing.T) {
	_, c := newTestClient(t, http.NewServeMux())

	_, err := c.Completions(context.Background(), "test", nil, CompletionOptions{})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "body is required")

	_, err = c.Completions(context.Background(), "test",
		json.RawMessage(`{not json`), CompletionOptions{})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "not valid JSON")
}

func TestCompletions_PropagatesAPIError(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/1/agents/test/completions", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusUnprocessableEntity)
		_, _ = w.Write([]byte(`{"detail":[{"msg":"messages: cannot be empty","loc":["body","messages"]}]}`))
	})

	_, c := newTestClient(t, mux)

	resp, err := c.Completions(context.Background(), "test",
		json.RawMessage(`{"messages":[]}`),
		CompletionOptions{Stream: true})
	require.Error(t, err)
	assert.Nil(t, resp, "resp must be nil on error so callers don't accidentally close a closed body")

	var apiErr *APIError
	require.True(t, errors.As(err, &apiErr))
	assert.Equal(t, http.StatusUnprocessableEntity, apiErr.StatusCode)
	assert.Equal(t, "messages: cannot be empty", apiErr.Detail)
}

func TestCompletions_PathEscapesAgentID(t *testing.T) {
	// Belt-and-suspenders: a real backend should never see weird IDs
	// (cobra positional args + UUID validation upstream), but if a
	// caller passes one we mustn't break URL parsing.
	mux := http.NewServeMux()
	hit := false
	mux.HandleFunc("/1/agents/weird%20id/completions", func(_ http.ResponseWriter, _ *http.Request) {
		hit = true
	})
	_, c := newTestClient(t, mux)
	resp, err := c.Completions(context.Background(), "weird id",
		json.RawMessage(`{}`), CompletionOptions{Stream: true})
	require.NoError(t, err)
	defer resp.Body.Close()
	assert.True(t, hit, "request must hit the URL-escaped path")
}

// Sanity check on the "test" literal — backend pattern-matches on it.
func TestCompletions_TestLiteralIsNotEscaped(t *testing.T) {
	mux := http.NewServeMux()
	hit := false
	mux.HandleFunc("/1/agents/test/completions", func(_ http.ResponseWriter, _ *http.Request) {
		hit = true
	})
	_, c := newTestClient(t, mux)
	resp, err := c.Completions(context.Background(), "test",
		json.RawMessage(`{}`), CompletionOptions{})
	require.NoError(t, err)
	defer resp.Body.Close()
	assert.True(t, hit)
}

// readAllString is a small helper used only here; kept local to avoid
// polluting other test files that don't need it.
func readAllString(t *testing.T, r io.Reader) string {
	t.Helper()
	b, err := io.ReadAll(r)
	require.NoError(t, err)
	return strings.TrimSpace(string(b))
}

func TestCompletions_BodyContentRoundTrip(t *testing.T) {
	// Confirms we POST exactly the bytes we were handed (no re-encode).
	wire := `{"messages":[{"role":"user","content":"x"}],"id":"conv-1"}`
	mux := http.NewServeMux()
	mux.HandleFunc("/1/agents/test/completions", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, wire, readAllString(t, r.Body))
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{}`))
	})
	_, c := newTestClient(t, mux)
	resp, err := c.Completions(context.Background(), "test",
		json.RawMessage(wire), CompletionOptions{})
	require.NoError(t, err)
	defer resp.Body.Close()
}

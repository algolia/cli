package agentstudio

// Tests for the cross-cutting infrastructure in client.go: NewClient
// validation, header injection, error mapping (checkResponse +
// extractDetail + sentinelFor), and context cancellation. Per-tag
// method tests live in <tag>_test.go (agents_test.go,
// completions_test.go, providers_test.go, configuration_test.go).
//
// The error-mapping and ctx-cancellation tests use ListAgents as a
// vehicle: it's the simplest GET endpoint in the package and exercising
// it keeps the assertions concrete. They are infra tests, not method
// tests — moving them to agents_test.go would obscure their intent.

import (
	"context"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// writeTestJSONResponse writes fixture JSON to w without calling
// http.ResponseWriter.Write directly. Static scanners flag that pattern as a
// potential HTML XSS sink even for JSON test doubles.
func writeTestJSONResponse(w http.ResponseWriter, body []byte) {
	var out io.Writer = w
	_, _ = out.Write(body)
}

// newTestClient is the shared httptest harness for every *_test.go in
// this package. Lives here because client.go owns Client construction.
func newTestClient(t *testing.T, handler http.Handler) (*httptest.Server, *Client) {
	t.Helper()
	ts := httptest.NewServer(handler)
	t.Cleanup(ts.Close)

	c, err := NewClient(Config{
		BaseURL:       ts.URL,
		ApplicationID: "APP123",
		APIKey:        "key-abc",
		UserID:        "cli-test",
		HTTPClient:    ts.Client(),
	})
	require.NoError(t, err)
	return ts, c
}

func TestNewClient_Validation(t *testing.T) {
	tests := []struct {
		name    string
		cfg     Config
		wantErr string
	}{
		{"missing baseURL", Config{ApplicationID: "x", APIKey: "y"}, "base url is required"},
		{"missing appID", Config{BaseURL: "http://x", APIKey: "y"}, "application id is required"},
		{"missing apiKey", Config{BaseURL: "http://x", ApplicationID: "y"}, "api key is required"},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			_, err := NewClient(tc.cfg)
			require.Error(t, err)
			assert.Contains(t, err.Error(), tc.wantErr)
		})
	}
}

func TestNewClient_TrimsTrailingSlashAndDefaults(t *testing.T) {
	c, err := NewClient(Config{
		BaseURL:       "https://x.example.com/",
		ApplicationID: "APP",
		APIKey:        "k",
	})
	require.NoError(t, err)
	assert.Equal(t, "https://x.example.com", c.cfg.BaseURL)
	assert.Equal(t, "algolia-cli/agentstudio", c.cfg.UserAgent)
	require.NotNil(t, c.httpClient)
	assert.Zero(t, c.httpClient.Timeout)
	_, ok := c.httpClient.Transport.(*http.Transport)
	assert.True(t, ok, "expected default *http.Transport for timeouts without killing SSE streams")
}

func TestCheckResponse_ErrorMapping(t *testing.T) {
	tests := []struct {
		name         string
		status       int
		body         string
		wantSentinel error
		wantDetail   string
	}{
		{
			name:         "401 → ErrUnauthorized",
			status:       http.StatusUnauthorized,
			body:         `{"detail":"Invalid API key"}`,
			wantSentinel: ErrUnauthorized,
			wantDetail:   "Invalid API key",
		},
		{
			name:         "403 missing ACL → ErrForbidden",
			status:       http.StatusForbidden,
			body:         `{"message":"API key is missing the following ACLs: settings."}`,
			wantSentinel: ErrForbidden,
			wantDetail:   "API key is missing the following ACLs: settings.",
		},
		{
			name:         "403 feature disabled → ErrFeatureDisabled",
			status:       http.StatusForbidden,
			body:         `{"message":"This feature is not enabled for this application."}`,
			wantSentinel: ErrFeatureDisabled,
			wantDetail:   "This feature is not enabled for this application.",
		},
		{
			name:         "404 → ErrNotFound",
			status:       http.StatusNotFound,
			body:         `{"detail":"Agent not found"}`,
			wantSentinel: ErrNotFound,
			wantDetail:   "Agent not found",
		},
		{
			name:         "500 → ErrServer",
			status:       http.StatusInternalServerError,
			body:         `{"detail":"oops"}`,
			wantSentinel: ErrServer,
			wantDetail:   "oops",
		},
		{
			name:         "422 (validation) wraps no sentinel but preserves detail",
			status:       http.StatusUnprocessableEntity,
			body:         `{"detail":[{"msg":"name is required"}]}`,
			wantSentinel: nil,
			wantDetail:   "name is required",
		},
		{
			name:         "non-JSON body falls back to raw text",
			status:       http.StatusBadGateway,
			body:         `<html>upstream broke</html>`,
			wantSentinel: ErrServer,
			wantDetail:   "<html>upstream broke</html>",
		},
		{
			// Regression for live behaviour: when the backend pairs a generic
			// "Input is invalid, see detail/body:" message with a structured
			// detail array, the structured msg wins.
			name:   "422 with both message and detail prefers structured detail.msg",
			status: http.StatusUnprocessableEntity,
			body: `{
				"message":"Input is invalid, see detail/body:",
				"detail":[{"loc":["path","agent_id"],"msg":"Input should be a valid UUID","type":"uuid_parsing"}]
			}`,
			wantSentinel: nil,
			wantDetail:   "Input should be a valid UUID",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mux := http.NewServeMux()
			mux.HandleFunc("/1/agents", func(w http.ResponseWriter, _ *http.Request) {
				w.WriteHeader(tc.status)
				writeTestJSONResponse(w, []byte(tc.body))
			})
			_, c := newTestClient(t, mux)

			_, err := c.ListAgents(context.Background(), ListAgentsParams{})
			require.Error(t, err)

			var apiErr *APIError
			require.True(t, errors.As(err, &apiErr), "expected *APIError, got %T", err)
			assert.Equal(t, tc.status, apiErr.StatusCode)
			assert.Equal(t, tc.wantDetail, apiErr.Detail)

			if tc.wantSentinel != nil {
				assert.True(t, errors.Is(err, tc.wantSentinel),
					"expected errors.Is(err, %v); got %v", tc.wantSentinel, err)
			}
		})
	}
}

func TestRequest_ContextCancellation(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/1/agents", func(_ http.ResponseWriter, r *http.Request) {
		<-r.Context().Done()
	})
	_, c := newTestClient(t, mux)

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	_, err := c.ListAgents(ctx, ListAgentsParams{})
	require.Error(t, err)
	assert.True(t, errors.Is(err, context.Canceled), "got %v", err)
}

func TestSetHeaders_OmitsUserIDWhenEmpty(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/1/agents", func(w http.ResponseWriter, r *http.Request) {
		assert.Empty(t, r.Header.Get(HeaderUserID))
		writeTestJSONResponse(w, []byte(`{"data":[],"pagination":{"page":1,"limit":10,"totalCount":0,"totalPages":0}}`))
	})

	ts := httptest.NewServer(mux)
	t.Cleanup(ts.Close)
	c, err := NewClient(Config{
		BaseURL:       ts.URL,
		ApplicationID: "APP",
		APIKey:        "k",
		HTTPClient:    ts.Client(),
		// no UserID
	})
	require.NoError(t, err)

	_, err = c.ListAgents(context.Background(), ListAgentsParams{})
	require.NoError(t, err)
}

package agentstudio

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

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
	assert.Equal(t, http.DefaultClient, c.httpClient)
}

func TestListAgents_Success(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/1/agents", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Equal(t, "APP123", r.Header.Get(HeaderApplicationID))
		assert.Equal(t, "key-abc", r.Header.Get(HeaderAPIKey))
		assert.Equal(t, "cli-test", r.Header.Get(HeaderUserID))
		assert.Equal(t, "application/json", r.Header.Get("Accept"))

		assert.Equal(t, "2", r.URL.Query().Get("page"))
		assert.Equal(t, "25", r.URL.Query().Get("limit"))
		assert.Equal(t, "prov-1", r.URL.Query().Get("providerId"))

		require.NoError(t, json.NewEncoder(w).Encode(PaginatedAgentsResponse{
			Data: []Agent{{
				ID:           "11111111-1111-1111-1111-111111111111",
				Name:         "Concierge",
				Status:       StatusDraft,
				Instructions: "Be helpful.",
			}},
			Pagination: PaginationMetadata{
				Page: 2, Limit: 25, TotalCount: 1, TotalPages: 1,
			},
		}))
	})

	_, c := newTestClient(t, mux)

	got, err := c.ListAgents(context.Background(), ListAgentsParams{
		Page:       2,
		Limit:      25,
		ProviderID: "prov-1",
	})
	require.NoError(t, err)
	require.Len(t, got.Data, 1)
	assert.Equal(t, "Concierge", got.Data[0].Name)
	assert.Equal(t, StatusDraft, got.Data[0].Status)
	assert.Equal(t, 25, got.Pagination.Limit)
}

func TestListAgents_OmitsZeroParams(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/1/agents", func(w http.ResponseWriter, r *http.Request) {
		assert.Empty(t, r.URL.RawQuery, "expected no query params for zero-valued params")
		_, _ = w.Write([]byte(`{"data":[],"pagination":{"page":1,"limit":10,"totalCount":0,"totalPages":0}}`))
	})

	_, c := newTestClient(t, mux)
	_, err := c.ListAgents(context.Background(), ListAgentsParams{})
	require.NoError(t, err)
}

func TestListAgents_ErrorMapping(t *testing.T) {
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
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mux := http.NewServeMux()
			mux.HandleFunc("/1/agents", func(w http.ResponseWriter, _ *http.Request) {
				w.WriteHeader(tc.status)
				_, _ = w.Write([]byte(tc.body))
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

func TestListAgents_ContextCancellation(t *testing.T) {
	// Server intentionally never responds; we cancel client-side.
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

func TestListAgents_OmitsUserIDHeaderWhenEmpty(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/1/agents", func(w http.ResponseWriter, r *http.Request) {
		assert.Empty(t, r.Header.Get(HeaderUserID))
		_, _ = w.Write([]byte(`{"data":[],"pagination":{"page":1,"limit":10,"totalCount":0,"totalPages":0}}`))
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

package agentstudio

// Tests for the provider methods that remain local. Provider CRUD moved to
// the official SDK; UpdateProvider (verbatim PATCH) and the model-catalog
// endpoints stay here.

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUpdateProvider_PartialPatch(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/1/providers/p1", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPatch, r.Method)
		body, _ := io.ReadAll(r.Body)
		assert.JSONEq(t, `{"name":"renamed"}`, string(body))
		_, _ = w.Write([]byte(`{
			"id":"p1","name":"renamed","providerName":"openai",
			"input":{"apiKey":"sk-XXX"},
			"createdAt":"2026-01-01T00:00:00Z","updatedAt":"2026-01-03T00:00:00Z"
		}`))
	})
	_, c := newTestClient(t, mux)
	got, err := c.UpdateProvider(context.Background(), "p1", json.RawMessage(`{"name":"renamed"}`))
	require.NoError(t, err)
	assert.Equal(t, "renamed", got.Name)
}

func TestUpdateProvider_RejectsEmptyIDOrBody(t *testing.T) {
	_, c := newTestClient(t, http.NewServeMux())

	_, err := c.UpdateProvider(context.Background(), "", json.RawMessage(`{}`))
	require.Error(t, err)
	assert.Contains(t, err.Error(), "provider id is required")

	_, err = c.UpdateProvider(context.Background(), "p1", nil)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "body is required")
}

func TestListProviderModels_Success(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/1/providers/models", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		_, _ = w.Write([]byte(`{"openai":["gpt-4o","gpt-4o-mini"],"anthropic":["claude-3-5-sonnet"]}`))
	})
	_, c := newTestClient(t, mux)
	got, err := c.ListProviderModels(context.Background())
	require.NoError(t, err)
	assert.Equal(t, []string{"gpt-4o", "gpt-4o-mini"}, got["openai"])
	assert.Equal(t, []string{"claude-3-5-sonnet"}, got["anthropic"])
}

func TestListModelsForProvider_PassesThroughRawJSON(t *testing.T) {
	// Spec leaves the response shape unspecified ([]string is likely
	// but not pinned). Keep this test loose: assert it round-trips
	// whatever shape arrives, and that the result is valid JSON.
	cases := []struct {
		name string
		body string
	}{
		{"array of strings", `["gpt-4o","gpt-4o-mini"]`},
		{"object envelope", `{"models":["gpt-4o"]}`},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			mux := http.NewServeMux()
			mux.HandleFunc("/1/providers/p1/models", func(w http.ResponseWriter, _ *http.Request) {
				writeTestJSONResponse(w, []byte(tc.body))
			})
			_, c := newTestClient(t, mux)
			got, err := c.ListModelsForProvider(context.Background(), "p1")
			require.NoError(t, err)
			assert.JSONEq(t, tc.body, string(got))
		})
	}
}

func TestListModelsForProvider_RejectsEmptyID(t *testing.T) {
	_, c := newTestClient(t, http.NewServeMux())
	_, err := c.ListModelsForProvider(context.Background(), "")
	require.Error(t, err)
}

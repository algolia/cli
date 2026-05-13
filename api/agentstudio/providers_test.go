package agentstudio

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestListProviders_Pagination(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/1/providers", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Equal(t, "2", r.URL.Query().Get("page"))
		assert.Equal(t, "5", r.URL.Query().Get("limit"))
		require.NoError(t, json.NewEncoder(w).Encode(PaginatedProvidersResponse{
			Data: []Provider{
				{ID: "p1", Name: "openai-prod", ProviderName: "openai", Input: json.RawMessage(`{"apiKey":"sk-..."}`)},
			},
			Pagination: PaginationMetadata{Page: 2, Limit: 5, TotalCount: 11, TotalPages: 3},
		}))
	})

	_, c := newTestClient(t, mux)

	got, err := c.ListProviders(context.Background(), ListProvidersParams{Page: 2, Limit: 5})
	require.NoError(t, err)
	require.Len(t, got.Data, 1)
	assert.Equal(t, "openai-prod", got.Data[0].Name)
	assert.Equal(t, "openai", got.Data[0].ProviderName)
	assert.Equal(t, 11, got.Pagination.TotalCount)
}

func TestListProviders_OmitsZeroPaginationParams(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/1/providers", func(w http.ResponseWriter, r *http.Request) {
		// Zero values must NOT be sent — backend would treat them as
		// explicit values and (in practice) 422 on `page=0`.
		assert.Equal(t, "", r.URL.RawQuery)
		require.NoError(t, json.NewEncoder(w).Encode(PaginatedProvidersResponse{Data: []Provider{}}))
	})
	_, c := newTestClient(t, mux)
	_, err := c.ListProviders(context.Background(), ListProvidersParams{})
	require.NoError(t, err)
}

func TestGetProvider_Success(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/1/providers/p1", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		_, _ = w.Write([]byte(`{
			"id":"p1",
			"name":"openai-prod",
			"providerName":"openai",
			"input":{"apiKey":"sk-XXX","baseUrl":null},
			"createdAt":"2026-01-01T00:00:00Z",
			"updatedAt":"2026-01-02T00:00:00Z"
		}`))
	})
	_, c := newTestClient(t, mux)

	got, err := c.GetProvider(context.Background(), "p1")
	require.NoError(t, err)
	assert.Equal(t, "p1", got.ID)
	assert.Equal(t, "openai", got.ProviderName)
	assert.JSONEq(t, `{"apiKey":"sk-XXX","baseUrl":null}`, string(got.Input))
}

func TestGetProvider_RejectsEmptyID(t *testing.T) {
	_, c := newTestClient(t, http.NewServeMux())
	_, err := c.GetProvider(context.Background(), "  ")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "provider id is required")
}

func TestGetProvider_NotFound(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/1/providers/missing", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte(`{"detail":"Provider not found"}`))
	})
	_, c := newTestClient(t, mux)
	_, err := c.GetProvider(context.Background(), "missing")
	require.Error(t, err)
	assert.True(t, errors.Is(err, ErrNotFound))
}

func TestCreateProvider_RoundTripsBody(t *testing.T) {
	wire := `{"name":"my-openai","providerName":"openai","input":{"apiKey":"sk-XXX"}}`
	mux := http.NewServeMux()
	mux.HandleFunc("/1/providers", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))
		body, _ := io.ReadAll(r.Body)
		assert.JSONEq(t, wire, string(body))
		w.WriteHeader(http.StatusCreated)
		_, _ = w.Write([]byte(`{
			"id":"p1","name":"my-openai","providerName":"openai",
			"input":{"apiKey":"sk-XXX"},
			"createdAt":"2026-01-01T00:00:00Z","updatedAt":"2026-01-01T00:00:00Z"
		}`))
	})
	_, c := newTestClient(t, mux)

	got, err := c.CreateProvider(context.Background(), json.RawMessage(wire))
	require.NoError(t, err)
	assert.Equal(t, "p1", got.ID)
}

func TestCreateProvider_RejectsEmptyBody(t *testing.T) {
	_, c := newTestClient(t, http.NewServeMux())
	_, err := c.CreateProvider(context.Background(), nil)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "body is required")
}

func TestCreateProvider_PropagatesValidationError(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/1/providers", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusUnprocessableEntity)
		_, _ = w.Write([]byte(`{"detail":[{"msg":"input.apiKey: field required","loc":["body","input","apiKey"]}]}`))
	})
	_, c := newTestClient(t, mux)
	_, err := c.CreateProvider(context.Background(), json.RawMessage(`{"name":"x","providerName":"openai","input":{}}`))
	require.Error(t, err)
	var apiErr *APIError
	require.True(t, errors.As(err, &apiErr))
	assert.Equal(t, http.StatusUnprocessableEntity, apiErr.StatusCode)
	assert.Equal(t, "input.apiKey: field required", apiErr.Detail)
}

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

func TestDeleteProvider_NoContentSuccess(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/1/providers/p1", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodDelete, r.Method)
		w.WriteHeader(http.StatusNoContent)
	})
	_, c := newTestClient(t, mux)
	require.NoError(t, c.DeleteProvider(context.Background(), "p1"))
}

func TestDeleteProvider_RejectsEmptyID(t *testing.T) {
	_, c := newTestClient(t, http.NewServeMux())
	require.Error(t, c.DeleteProvider(context.Background(), ""))
}

func TestDeleteProvider_PropagatesConflict(t *testing.T) {
	// If a provider has agents pointing at it, the backend may 409.
	// CLI surfaces the structured detail unchanged.
	mux := http.NewServeMux()
	mux.HandleFunc("/1/providers/p1", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusConflict)
		_, _ = w.Write([]byte(`{"detail":"Provider is in use by 3 agent(s); detach first"}`))
	})
	_, c := newTestClient(t, mux)
	err := c.DeleteProvider(context.Background(), "p1")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "in use by 3 agent")
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

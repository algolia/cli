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

func TestListSecretKeys_PaginationAndQuery(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/1/secret-keys", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "2", r.URL.Query().Get("page"))
		assert.Equal(t, "50", r.URL.Query().Get("limit"))
		_, _ = w.Write([]byte(`{"data":[],"pagination":{"page":2,"limit":50,"totalCount":0,"totalPages":0}}`))
	})
	_, c := newTestClient(t, mux)
	got, err := c.ListSecretKeys(context.Background(), ListSecretKeysParams{Page: 2, Limit: 50})
	require.NoError(t, err)
	assert.Equal(t, 2, got.Pagination.Page)
}

func TestListSecretKeys_OmitsZeroParams(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/1/secret-keys", func(w http.ResponseWriter, r *http.Request) {
		assert.Empty(t, r.URL.RawQuery)
		_, _ = w.Write([]byte(`{"data":[],"pagination":{"page":1,"limit":10,"totalCount":0,"totalPages":0}}`))
	})
	_, c := newTestClient(t, mux)
	_, err := c.ListSecretKeys(context.Background(), ListSecretKeysParams{})
	require.NoError(t, err)
}

func TestListSecretKeys_ForbiddenMaps(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/1/secret-keys", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusForbidden)
		_, _ = w.Write([]byte(`{"message":"Admin API key required."}`))
	})
	_, c := newTestClient(t, mux)
	_, err := c.ListSecretKeys(context.Background(), ListSecretKeysParams{})
	require.Error(t, err)
	assert.True(t, errors.Is(err, ErrForbidden))
}

func TestCreateSecretKey_RoundTrip(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/1/secret-keys", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		body, _ := io.ReadAll(r.Body)
		var got SecretKeyCreate
		require.NoError(t, json.Unmarshal(body, &got))
		assert.Equal(t, "k1", got.Name)
		assert.Equal(t, []string{"a1"}, got.AgentIDs)
		w.WriteHeader(http.StatusCreated)
		_, _ = w.Write([]byte(`{"id":"id1","name":"k1","value":"sk-real","createdAt":"2026-01-01T00:00:00Z","updatedAt":"2026-01-01T00:00:00Z","lastUsedAt":null,"isDefault":false,"agentIds":["a1"]}`))
	})
	_, c := newTestClient(t, mux)
	got, err := c.CreateSecretKey(context.Background(), SecretKeyCreate{Name: "k1", AgentIDs: []string{"a1"}})
	require.NoError(t, err)
	assert.Equal(t, "sk-real", got.Value)
	assert.Nil(t, got.LastUsedAt)
}

func TestCreateSecretKey_RejectsEmptyName(t *testing.T) {
	_, c := newTestClient(t, http.NewServeMux())
	_, err := c.CreateSecretKey(context.Background(), SecretKeyCreate{Name: " "})
	require.Error(t, err)
}

func TestCreateSecretKey_OmitsAgentIDsWhenEmpty(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/1/secret-keys", func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		var got map[string]any
		require.NoError(t, json.Unmarshal(body, &got))
		_, hasAgentIDs := got["agentIds"]
		assert.False(t, hasAgentIDs, "agentIds should be omitted when empty (got body=%s)", string(body))
		w.WriteHeader(http.StatusCreated)
		_, _ = w.Write([]byte(`{"id":"id1","name":"k1","value":"v","createdAt":"2026-01-01T00:00:00Z","updatedAt":"2026-01-01T00:00:00Z","lastUsedAt":null,"isDefault":false,"agentIds":[]}`))
	})
	_, c := newTestClient(t, mux)
	_, err := c.CreateSecretKey(context.Background(), SecretKeyCreate{Name: "k1"})
	require.NoError(t, err)
}

func TestUpdateSecretKey_RejectsNoChanges(t *testing.T) {
	_, c := newTestClient(t, http.NewServeMux())
	_, err := c.UpdateSecretKey(context.Background(), "id1", SecretKeyPatch{})
	require.Error(t, err)
}

func TestUpdateSecretKey_NameOnly(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/1/secret-keys/id1", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPatch, r.Method)
		body, _ := io.ReadAll(r.Body)
		var got map[string]any
		require.NoError(t, json.Unmarshal(body, &got))
		assert.Equal(t, "renamed", got["name"])
		_, hasAgentIDs := got["agentIds"]
		assert.False(t, hasAgentIDs, "agentIds must be omitted when nil (body=%s)", string(body))
		_, _ = w.Write([]byte(`{"id":"id1","name":"renamed","value":"v","createdAt":"2026-01-01T00:00:00Z","updatedAt":"2026-01-01T00:00:00Z","lastUsedAt":null,"isDefault":false,"agentIds":[]}`))
	})
	_, c := newTestClient(t, mux)
	name := "renamed"
	_, err := c.UpdateSecretKey(context.Background(), "id1", SecretKeyPatch{Name: &name})
	require.NoError(t, err)
}

func TestUpdateSecretKey_EmptyAgentIDsClearsList(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/1/secret-keys/id1", func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		var got map[string]any
		require.NoError(t, json.Unmarshal(body, &got))
		v, hasAgentIDs := got["agentIds"]
		assert.True(t, hasAgentIDs)
		assert.Equal(t, []any{}, v)
		_, _ = w.Write([]byte(`{"id":"id1","name":"n","value":"v","createdAt":"2026-01-01T00:00:00Z","updatedAt":"2026-01-01T00:00:00Z","lastUsedAt":null,"isDefault":false,"agentIds":[]}`))
	})
	_, c := newTestClient(t, mux)
	empty := []string{}
	_, err := c.UpdateSecretKey(context.Background(), "id1", SecretKeyPatch{AgentIDs: &empty})
	require.NoError(t, err)
}

func TestDeleteSecretKey_NoContent(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/1/secret-keys/id1", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodDelete, r.Method)
		w.WriteHeader(http.StatusNoContent)
	})
	_, c := newTestClient(t, mux)
	require.NoError(t, c.DeleteSecretKey(context.Background(), "id1"))
}

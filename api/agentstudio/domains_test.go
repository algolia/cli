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

func TestListAllowedDomains_Success(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/1/agents/agent-1/allowed-domains", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		_, _ = w.Write([]byte(`{"domains":[{
			"id":"d1","appId":"APP","agentId":"agent-1","domain":"https://x.test",
			"createdAt":"2026-01-01T00:00:00Z","updatedAt":"2026-01-01T00:00:00Z"
		}]}`))
	})
	_, c := newTestClient(t, mux)

	got, err := c.ListAllowedDomains(context.Background(), "agent-1")
	require.NoError(t, err)
	require.Len(t, got.Domains, 1)
	assert.Equal(t, "https://x.test", got.Domains[0].Domain)
}

func TestListAllowedDomains_RejectsEmptyAgentID(t *testing.T) {
	_, c := newTestClient(t, http.NewServeMux())
	_, err := c.ListAllowedDomains(context.Background(), " ")
	require.Error(t, err)
}

func TestGetAllowedDomain_Success(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/1/agents/agent-1/allowed-domains/d1", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		_, _ = w.Write([]byte(`{"id":"d1","appId":"APP","agentId":"agent-1","domain":"x","createdAt":"2026-01-01T00:00:00Z","updatedAt":"2026-01-01T00:00:00Z"}`))
	})
	_, c := newTestClient(t, mux)
	got, err := c.GetAllowedDomain(context.Background(), "agent-1", "d1")
	require.NoError(t, err)
	assert.Equal(t, "d1", got.ID)
}

func TestGetAllowedDomain_NotFound(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/1/agents/agent-1/allowed-domains/missing", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte(`{"detail":"Allowed domain not found"}`))
	})
	_, c := newTestClient(t, mux)
	_, err := c.GetAllowedDomain(context.Background(), "agent-1", "missing")
	require.Error(t, err)
	assert.True(t, errors.Is(err, ErrNotFound))
}

func TestCreateAllowedDomain_SendsBody(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/1/agents/agent-1/allowed-domains", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))
		body, _ := io.ReadAll(r.Body)
		assert.JSONEq(t, `{"domain":"https://x.test"}`, string(body))
		w.WriteHeader(http.StatusCreated)
		_, _ = w.Write([]byte(`{"id":"d1","appId":"APP","agentId":"agent-1","domain":"https://x.test","createdAt":"2026-01-01T00:00:00Z","updatedAt":"2026-01-01T00:00:00Z"}`))
	})
	_, c := newTestClient(t, mux)
	got, err := c.CreateAllowedDomain(context.Background(), "agent-1", "https://x.test")
	require.NoError(t, err)
	assert.Equal(t, "d1", got.ID)
}

func TestCreateAllowedDomain_RejectsEmptyDomain(t *testing.T) {
	_, c := newTestClient(t, http.NewServeMux())
	_, err := c.CreateAllowedDomain(context.Background(), "agent-1", "")
	require.Error(t, err)
}

func TestDeleteAllowedDomain_NoContent(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/1/agents/agent-1/allowed-domains/d1", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodDelete, r.Method)
		w.WriteHeader(http.StatusNoContent)
	})
	_, c := newTestClient(t, mux)
	require.NoError(t, c.DeleteAllowedDomain(context.Background(), "agent-1", "d1"))
}

func TestBulkInsertAllowedDomains_RoundTrip(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/1/agents/agent-1/allowed-domains/bulk", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		var got map[string][]string
		require.NoError(t, json.NewDecoder(r.Body).Decode(&got))
		assert.Equal(t, []string{"a", "b"}, got["domains"])
		w.WriteHeader(http.StatusCreated)
		_, _ = w.Write([]byte(`{"domains":[
			{"id":"1","appId":"APP","agentId":"agent-1","domain":"a","createdAt":"2026-01-01T00:00:00Z","updatedAt":"2026-01-01T00:00:00Z"},
			{"id":"2","appId":"APP","agentId":"agent-1","domain":"b","createdAt":"2026-01-01T00:00:00Z","updatedAt":"2026-01-01T00:00:00Z"}
		]}`))
	})
	_, c := newTestClient(t, mux)
	got, err := c.BulkInsertAllowedDomains(context.Background(), "agent-1", []string{"a", "b"})
	require.NoError(t, err)
	require.Len(t, got.Domains, 2)
}

func TestBulkInsertAllowedDomains_RejectsEmptyList(t *testing.T) {
	_, c := newTestClient(t, http.NewServeMux())
	_, err := c.BulkInsertAllowedDomains(context.Background(), "agent-1", nil)
	require.Error(t, err)
}

func TestBulkDeleteAllowedDomains_SendsIDsBody(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/1/agents/agent-1/allowed-domains/bulk", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodDelete, r.Method)
		var got map[string][]string
		require.NoError(t, json.NewDecoder(r.Body).Decode(&got))
		assert.Equal(t, []string{"d1", "d2"}, got["domainIds"])
		w.WriteHeader(http.StatusNoContent)
	})
	_, c := newTestClient(t, mux)
	require.NoError(t, c.BulkDeleteAllowedDomains(context.Background(), "agent-1", []string{"d1", "d2"}))
}

func TestBulkDeleteAllowedDomains_RejectsEmptyList(t *testing.T) {
	_, c := newTestClient(t, http.NewServeMux())
	require.Error(t, c.BulkDeleteAllowedDomains(context.Background(), "agent-1", nil))
}

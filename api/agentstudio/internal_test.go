package agentstudio

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetStatus_NoAuthHeaders(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/status", func(w http.ResponseWriter, r *http.Request) {
		assert.Empty(t, r.Header.Get("X-Algolia-Application-Id"))
		assert.Empty(t, r.Header.Get("X-Algolia-API-Key"))
		_, _ = w.Write([]byte(`{"status":"ok","version":"abc123","migration_revision":"r1"}`))
	})
	_, c := newTestClient(t, mux)
	got, err := c.GetStatus(context.Background())
	require.NoError(t, err)
	require.NotNil(t, got["status"])
	assert.Equal(t, "ok", *got["status"])
}

func TestGetProviderModelDefaults_Success(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/1/providers/models/defaults", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "APP123", r.Header.Get("X-Algolia-Application-Id"))
		_, _ = w.Write([]byte(`{"openai":"gpt-4.1-mini","anthropic":"claude-haiku-4-5"}`))
	})
	_, c := newTestClient(t, mux)
	got, err := c.GetProviderModelDefaults(context.Background())
	require.NoError(t, err)
	assert.Equal(t, "gpt-4.1-mini", got["openai"])
}

func TestAgentMemorize_DoubledPath(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/1/agents/agents/agent-1/memorize", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))
		body, _ := io.ReadAll(r.Body)
		assert.JSONEq(t, `{"providerID":"p","model":"m","messages":[]}`, string(body))
		_, _ = w.Write([]byte(`{"savedMemories":[],"deletedIds":[],"deleteTaskIds":[],"message":"ok"}`))
	})
	_, c := newTestClient(t, mux)
	out, err := c.AgentMemorize(
		context.Background(),
		"agent-1",
		json.RawMessage(`{"providerID":"p","model":"m","messages":[]}`),
	)
	require.NoError(t, err)
	assert.Contains(t, string(out), `"message":"ok"`)
}

func TestAgentPonder_RequiresAgentIDAndBody(t *testing.T) {
	_, c := newTestClient(t, http.NewServeMux())
	_, err := c.AgentPonder(context.Background(), "", json.RawMessage(`{}`))
	require.Error(t, err)
	_, err = c.AgentPonder(context.Background(), "agent-1", nil)
	require.Error(t, err)
}

func TestAgentConsolidate_PassesThroughErrors(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/1/agents/agents/agent-1/consolidate", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte(`{"message":"Provider not found or invalid"}`))
	})
	_, c := newTestClient(t, mux)
	_, err := c.AgentConsolidate(
		context.Background(),
		"agent-1",
		json.RawMessage(`{"providerID":"x","model":"m","messages":[]}`),
	)
	require.Error(t, err)
}

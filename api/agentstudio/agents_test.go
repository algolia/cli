package agentstudio

// Tests for the Agents API tag — the methods defined in agents.go.
//
// Only DuplicateAgent remains local (the rest of the agent CRUD/lifecycle
// surface moved to the official SDK). Cross-cutting concerns (NewClient
// validation, error mapping via checkResponse, header injection) live in
// client_test.go and use DuplicateAgent as a vehicle.

import (
	"context"
	"errors"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDuplicateAgent_Success(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/1/agents/abc-123/duplicate", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		// No request body for the lifecycle endpoint.
		assert.Empty(t, r.Header.Get("Content-Type"))
		body, _ := io.ReadAll(r.Body)
		assert.Empty(t, body)

		writeTestJSONResponse(w, []byte(`{
			"id":"new-id",
			"name":"Copy",
			"status":"draft",
			"instructions":"x",
			"createdAt":"2025-01-01T00:00:00Z"
		}`))
	})

	_, c := newTestClient(t, mux)

	got, err := c.DuplicateAgent(context.Background(), "abc-123")
	require.NoError(t, err)
	assert.Equal(t, "new-id", got.ID)
	assert.Equal(t, StatusDraft, got.Status)
}

func TestDuplicateAgent_RejectsEmptyID(t *testing.T) {
	_, c := newTestClient(t, http.NewServeMux())
	_, err := c.DuplicateAgent(context.Background(), "  ")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "agent id is required")
}

func TestDuplicateAgent_PathEscapesID(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/1/agents/weird%2Fid/duplicate", func(w http.ResponseWriter, _ *http.Request) {
		writeTestJSONResponse(w, []byte(`{
			"id":"weird/id","name":"x","status":"draft",
			"instructions":"","createdAt":"2025-01-01T00:00:00Z"
		}`))
	})

	_, c := newTestClient(t, mux)

	got, err := c.DuplicateAgent(context.Background(), "weird/id")
	require.NoError(t, err)
	assert.Equal(t, "weird/id", got.ID)
}

func TestDuplicateAgent_NotFound(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/1/agents/missing/duplicate", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte(`{"detail":"Agent not found"}`))
	})

	_, c := newTestClient(t, mux)
	_, err := c.DuplicateAgent(context.Background(), "missing")
	require.Error(t, err)
	assert.True(t, errors.Is(err, ErrNotFound))
}

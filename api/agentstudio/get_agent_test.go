package agentstudio

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetAgent_Success(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc(
		"/1/agents/abc-123",
		func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, http.MethodGet, r.Method)
			assert.Equal(t, "APP123", r.Header.Get(HeaderApplicationID))
			_, _ = w.Write([]byte(`{
				"id":"abc-123",
				"name":"Concierge",
				"status":"published",
				"instructions":"Be helpful.",
				"createdAt":"2025-01-02T03:04:05Z"
			}`))
		},
	)

	_, c := newTestClient(t, mux)

	got, err := c.GetAgent(context.Background(), "abc-123")
	require.NoError(t, err)
	assert.Equal(t, "abc-123", got.ID)
	assert.Equal(t, "Concierge", got.Name)
	assert.Equal(t, StatusPublished, got.Status)
}

func TestGetAgent_EmptyIDRejected(t *testing.T) {
	// No server hit because validation fails before transport.
	_, c := newTestClient(t, http.NewServeMux())

	_, err := c.GetAgent(context.Background(), "  ")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "agent id is required")
}

func TestGetAgent_NotFound(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/1/agents/missing", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte(`{"detail":"Agent not found"}`))
	})

	_, c := newTestClient(t, mux)

	_, err := c.GetAgent(context.Background(), "missing")
	require.Error(t, err)
	assert.True(t, errors.Is(err, ErrNotFound), "expected ErrNotFound, got %v", err)
}

func TestGetAgent_PathEscapesID(t *testing.T) {
	// Confirms IDs with reserved chars are encoded into the path so they
	// can't poison the URL (e.g., "../" attempts).
	mux := http.NewServeMux()
	mux.HandleFunc(
		"/1/agents/weird%2Fid",
		func(w http.ResponseWriter, _ *http.Request) {
			_, _ = w.Write([]byte(`{
				"id":"weird/id",
				"name":"x",
				"status":"draft",
				"instructions":"",
				"createdAt":"2025-01-01T00:00:00Z"
			}`))
		},
	)

	_, c := newTestClient(t, mux)

	got, err := c.GetAgent(context.Background(), "weird/id")
	require.NoError(t, err)
	assert.Equal(t, "weird/id", got.ID)
}

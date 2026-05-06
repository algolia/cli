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

func TestCreateFeedback_RoundTrip(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/1/feedback", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		body, _ := io.ReadAll(r.Body)
		var got map[string]any
		require.NoError(t, json.Unmarshal(body, &got))
		assert.Equal(t, "msg-1", got["messageId"])
		assert.Equal(t, "agent-1", got["agentId"])
		assert.EqualValues(t, 1, got["vote"])
		w.WriteHeader(http.StatusCreated)
		_, _ = w.Write([]byte(`{"id":"fb1","agentId":"agent-1","messageId":"msg-1","vote":1,"tags":["x"],"notes":null,"model":null,"createdAt":"2026-01-01T00:00:00Z","updatedAt":"2026-01-01T00:00:00Z"}`))
	})
	_, c := newTestClient(t, mux)
	got, err := c.CreateFeedback(context.Background(), FeedbackCreate{
		MessageID: "msg-1",
		AgentID:   "agent-1",
		Vote:      1,
	})
	require.NoError(t, err)
	assert.Equal(t, "fb1", got.ID)
}

func TestCreateFeedback_RejectsBadVote(t *testing.T) {
	_, c := newTestClient(t, http.NewServeMux())
	_, err := c.CreateFeedback(context.Background(), FeedbackCreate{
		MessageID: "msg-1", AgentID: "a", Vote: 5,
	})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "0 or 1")
}

func TestCreateFeedback_OmitsEmptyOptionals(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/1/feedback", func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		var got map[string]any
		require.NoError(t, json.Unmarshal(body, &got))
		_, hasTags := got["tags"]
		_, hasNotes := got["notes"]
		assert.False(t, hasTags, "tags should be omitted when empty (body=%s)", string(body))
		assert.False(t, hasNotes, "notes should be omitted when empty (body=%s)", string(body))
		w.WriteHeader(http.StatusCreated)
		_, _ = w.Write([]byte(`{"id":"fb1","agentId":"a","messageId":"m","vote":0,"tags":[],"notes":null,"model":null,"createdAt":"2026-01-01T00:00:00Z","updatedAt":"2026-01-01T00:00:00Z"}`))
	})
	_, c := newTestClient(t, mux)
	_, err := c.CreateFeedback(context.Background(), FeedbackCreate{
		MessageID: "m", AgentID: "a", Vote: 0,
	})
	require.NoError(t, err)
}

func TestCreateFeedback_NotFoundMaps(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/1/feedback", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte(`{"message":"Message with ID nonexistent-msg not found."}`))
	})
	_, c := newTestClient(t, mux)
	_, err := c.CreateFeedback(context.Background(), FeedbackCreate{
		MessageID: "nope", AgentID: "a", Vote: 1,
	})
	require.Error(t, err)
	assert.True(t, errors.Is(err, ErrNotFound))
}

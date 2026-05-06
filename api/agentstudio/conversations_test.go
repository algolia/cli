package agentstudio

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Tests for the Conversations API tag — the methods defined in
// conversations.go. Cross-cutting infra is covered in client_test.go.

func TestListConversations_Success(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc(
		"/1/agents/agent-1/conversations",
		func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, http.MethodGet, r.Method)
			q := r.URL.Query()
			assert.Equal(t, "2", q.Get("page"))
			assert.Equal(t, "5", q.Get("limit"))
			assert.Equal(t, "2026-01-01", q.Get("startDate"))
			assert.Equal(t, "2026-01-31", q.Get("endDate"))
			assert.Equal(t, "true", q.Get("includeFeedback"))
			assert.Equal(t, "1", q.Get("feedbackVote"))

			_, _ = w.Write([]byte(`{
				"data":[{
					"id":"c1","agentId":"agent-1","title":"Test conv",
					"createdAt":"2026-01-15T00:00:00Z","updatedAt":"2026-01-15T00:01:00Z",
					"messageCount":4,"totalInputTokens":120,"totalOutputTokens":340,"totalTokens":460,
					"isFromDashboard":false
				}],
				"pagination":{"page":2,"limit":5,"totalCount":1,"totalPages":1}
			}`))
		},
	)
	_, c := newTestClient(t, mux)

	vote := 1
	got, err := c.ListConversations(context.Background(), "agent-1", ListConversationsParams{
		Page:            2,
		Limit:           5,
		StartDate:       "2026-01-01",
		EndDate:         "2026-01-31",
		IncludeFeedback: true,
		FeedbackVote:    &vote,
	})
	require.NoError(t, err)
	require.Len(t, got.Data, 1)
	assert.Equal(t, "c1", got.Data[0].ID)
	assert.Equal(t, "agent-1", got.Data[0].AgentID)
	require.NotNil(t, got.Data[0].Title)
	assert.Equal(t, "Test conv", *got.Data[0].Title)
	assert.Equal(t, 460, got.Data[0].TotalTokens)
}

func TestListConversations_OmitsZeroParams(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc(
		"/1/agents/agent-1/conversations",
		func(w http.ResponseWriter, r *http.Request) {
			assert.Empty(t, r.URL.RawQuery, "expected no query params for zero values")
			_, _ = w.Write([]byte(`{"data":[],"pagination":{"page":1,"limit":20,"totalCount":0,"totalPages":0}}`))
		},
	)
	_, c := newTestClient(t, mux)
	_, err := c.ListConversations(context.Background(), "agent-1", ListConversationsParams{})
	require.NoError(t, err)
}

func TestListConversations_EmitsZeroVoteAsExplicitFilter(t *testing.T) {
	// Regression: FeedbackVote=&0 must reach the wire because 0
	// (downvote) is a meaningful filter, distinct from "no filter".
	mux := http.NewServeMux()
	mux.HandleFunc(
		"/1/agents/agent-1/conversations",
		func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "0", r.URL.Query().Get("feedbackVote"))
			_, _ = w.Write([]byte(`{"data":[],"pagination":{"page":1,"limit":20,"totalCount":0,"totalPages":0}}`))
		},
	)
	_, c := newTestClient(t, mux)

	zero := 0
	_, err := c.ListConversations(context.Background(), "agent-1", ListConversationsParams{
		FeedbackVote: &zero,
	})
	require.NoError(t, err)
}

func TestListConversations_RejectsEmptyAgentID(t *testing.T) {
	_, c := newTestClient(t, http.NewServeMux())
	_, err := c.ListConversations(context.Background(), " ", ListConversationsParams{})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "agent id is required")
}

func TestGetConversation_PassThroughBody(t *testing.T) {
	body := `{
		"id":"c1","agentId":"agent-1","title":"x",
		"createdAt":"2026-01-15T00:00:00Z","updatedAt":"2026-01-15T00:01:00Z",
		"messages":[{"role":"user","content":[{"type":"text","text":"hi"}]}]
	}`
	mux := http.NewServeMux()
	mux.HandleFunc(
		"/1/agents/agent-1/conversations/c1",
		func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "true", r.URL.Query().Get("includeFeedback"))
			_, _ = w.Write([]byte(body))
		},
	)
	_, c := newTestClient(t, mux)

	raw, err := c.GetConversation(context.Background(), "agent-1", "c1", true)
	require.NoError(t, err)

	// The bytes round-trip — no unmarshal/remarshal in the client path.
	var roundTrip map[string]any
	require.NoError(t, json.Unmarshal(raw, &roundTrip))
	assert.Equal(t, "c1", roundTrip["id"])
	assert.Contains(t, roundTrip, "messages")
}

func TestGetConversation_OmitsIncludeFeedbackWhenFalse(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc(
		"/1/agents/agent-1/conversations/c1",
		func(w http.ResponseWriter, r *http.Request) {
			assert.Empty(t, r.URL.RawQuery)
			_, _ = w.Write([]byte(`{"id":"c1","agentId":"agent-1","createdAt":"2026-01-15T00:00:00Z","updatedAt":"2026-01-15T00:01:00Z","messages":[]}`))
		},
	)
	_, c := newTestClient(t, mux)

	_, err := c.GetConversation(context.Background(), "agent-1", "c1", false)
	require.NoError(t, err)
}

func TestGetConversation_RejectsEmptyIDs(t *testing.T) {
	_, c := newTestClient(t, http.NewServeMux())

	_, err := c.GetConversation(context.Background(), "", "c1", false)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "agent id is required")

	_, err = c.GetConversation(context.Background(), "agent-1", "", false)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "conversation id is required")
}

func TestGetConversation_NotFound(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc(
		"/1/agents/agent-1/conversations/missing",
		func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusNotFound)
			_, _ = w.Write([]byte(`{"detail":"Conversation not found"}`))
		},
	)
	_, c := newTestClient(t, mux)

	_, err := c.GetConversation(context.Background(), "agent-1", "missing", false)
	require.Error(t, err)
	assert.True(t, errors.Is(err, ErrNotFound), "expected ErrNotFound, got %v", err)
}

func TestDeleteConversation_NoContentSuccess(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc(
		"/1/agents/agent-1/conversations/c1",
		func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, http.MethodDelete, r.Method)
			w.WriteHeader(http.StatusNoContent)
		},
	)
	_, c := newTestClient(t, mux)

	require.NoError(t, c.DeleteConversation(context.Background(), "agent-1", "c1"))
}

func TestDeleteConversation_RejectsEmptyIDs(t *testing.T) {
	_, c := newTestClient(t, http.NewServeMux())
	require.Error(t, c.DeleteConversation(context.Background(), "", "c1"))
	require.Error(t, c.DeleteConversation(context.Background(), "agent-1", ""))
}

func TestPurgeConversations_SendsDateRange(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc(
		"/1/agents/agent-1/conversations",
		func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, http.MethodDelete, r.Method)
			q := r.URL.Query()
			assert.Equal(t, "2026-01-01", q.Get("startDate"))
			assert.Equal(t, "2026-01-31", q.Get("endDate"))
			w.WriteHeader(http.StatusNoContent)
		},
	)
	_, c := newTestClient(t, mux)

	err := c.PurgeConversations(context.Background(), "agent-1", PurgeConversationsParams{
		StartDate: "2026-01-01",
		EndDate:   "2026-01-31",
	})
	require.NoError(t, err)
}

func TestPurgeConversations_DatelessSendsEmptyQuery(t *testing.T) {
	// Wire-level: empty params = empty query. The cmd layer enforces
	// the --all guardrail; the client mirrors the wire shape.
	mux := http.NewServeMux()
	mux.HandleFunc(
		"/1/agents/agent-1/conversations",
		func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, http.MethodDelete, r.Method)
			assert.Empty(t, r.URL.RawQuery)
			w.WriteHeader(http.StatusNoContent)
		},
	)
	_, c := newTestClient(t, mux)

	require.NoError(t, c.PurgeConversations(context.Background(), "agent-1", PurgeConversationsParams{}))
}

func TestExportConversations_PassThroughBody(t *testing.T) {
	body := `[{"id":"c1","messages":[]},{"id":"c2","messages":[]}]`
	mux := http.NewServeMux()
	mux.HandleFunc(
		"/1/agents/agent-1/conversations/export",
		func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, http.MethodGet, r.Method)
			assert.Equal(t, "2026-01-01", r.URL.Query().Get("startDate"))
			_, _ = w.Write([]byte(body))
		},
	)
	_, c := newTestClient(t, mux)

	raw, err := c.ExportConversations(context.Background(), "agent-1", ExportConversationsParams{
		StartDate: "2026-01-01",
	})
	require.NoError(t, err)
	assert.Equal(t, body, strings.TrimSpace(string(raw)))
}

func TestExportConversations_FeatureDisabledSurfaces(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc(
		"/1/agents/agent-1/conversations/export",
		func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusForbidden)
			_, _ = w.Write([]byte(`{"message":"This feature is not enabled for this application."}`))
		},
	)
	_, c := newTestClient(t, mux)

	_, err := c.ExportConversations(context.Background(), "agent-1", ExportConversationsParams{})
	require.Error(t, err)
	assert.True(t, errors.Is(err, ErrFeatureDisabled), "expected ErrFeatureDisabled, got %v", err)
}

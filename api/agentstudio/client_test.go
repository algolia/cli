package agentstudio

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newTestClient(handler http.Handler) (*httptest.Server, *Client) {
	ts := httptest.NewServer(handler)
	c := NewClientWithHTTPClient("APP", "KEY", ts.Client())
	c.BaseURL = ts.URL + "/"
	return ts, c
}

func TestRequest_SetsAuthHeaders(t *testing.T) {
	ts, c := newTestClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "APP", r.Header.Get("X-Algolia-Application-Id"))
		assert.Equal(t, "KEY", r.Header.Get("X-Algolia-API-Key"))
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"data":[],"pagination":{"page":1,"limit":10,"totalCount":0,"totalPages":0}}`))
	}))
	defer ts.Close()

	_, err := c.ListAgents(0, 0)
	require.NoError(t, err)
}

func TestListAgents_Success(t *testing.T) {
	ts, c := newTestClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Equal(t, "/1/agents", r.URL.Path)
		assert.Equal(t, "2", r.URL.Query().Get("page"))
		assert.Equal(t, "5", r.URL.Query().Get("limit"))
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{
			"data": [{"id":"a1","name":"Agent 1","status":"published","instructions":"hi","createdAt":"2026-01-01T00:00:00Z","updatedAt":"2026-01-01T00:00:00Z"}],
			"pagination": {"page":2,"limit":5,"totalCount":1,"totalPages":1}
		}`))
	}))
	defer ts.Close()

	res, err := c.ListAgents(2, 5)
	require.NoError(t, err)
	assert.Len(t, res.Data, 1)
	assert.Equal(t, "a1", res.Data[0].ID)
	assert.Equal(t, 2, res.Pagination.Page)
}

func TestGetAgent_NotFound(t *testing.T) {
	ts, c := newTestClient(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte(`{"detail":"Agent not found"}`))
	}))
	defer ts.Close()

	_, err := c.GetAgent("missing")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "404")
	assert.Contains(t, err.Error(), "Agent not found")
}

func TestCreateAgent_Success(t *testing.T) {
	ts, c := newTestClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, "/1/agents", r.URL.Path)
		var body AgentConfigCreate
		require.NoError(t, json.NewDecoder(r.Body).Decode(&body))
		assert.Equal(t, "Helper", body.Name)
		assert.Equal(t, "be helpful", body.Instructions)
		w.WriteHeader(http.StatusCreated)
		_, _ = w.Write([]byte(`{"id":"a2","name":"Helper","status":"draft","instructions":"be helpful","createdAt":"2026-01-01T00:00:00Z","updatedAt":"2026-01-01T00:00:00Z"}`))
	}))
	defer ts.Close()

	got, err := c.CreateAgent(AgentConfigCreate{Name: "Helper", Instructions: "be helpful"})
	require.NoError(t, err)
	assert.Equal(t, "a2", got.ID)
}

func TestDeleteAgent_NoContent(t *testing.T) {
	ts, c := newTestClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodDelete, r.Method)
		assert.Equal(t, "/1/agents/a3", r.URL.Path)
		w.WriteHeader(http.StatusNoContent)
	}))
	defer ts.Close()

	require.NoError(t, c.DeleteAgent("a3"))
}

func TestCreateCompletion_RawBody(t *testing.T) {
	ts, c := newTestClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, "/1/agents/a1/completions", r.URL.Path)
		assert.Equal(t, "ai-sdk-5", r.URL.Query().Get("compatibilityMode"))
		assert.Equal(t, "false", r.URL.Query().Get("stream"))
		var req AgentCompletionRequest
		require.NoError(t, json.NewDecoder(r.Body).Decode(&req))
		require.NotNil(t, req.ID)
		assert.Equal(t, "conv-1", *req.ID)
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"answer":"hello"}`))
	}))
	defer ts.Close()

	convID := "conv-1"
	stream := false
	rc, err := c.CreateCompletion("a1", AgentCompletionRequest{ID: &convID}, CompletionParams{
		CompatibilityMode: "ai-sdk-5",
		Stream:            &stream,
	})
	require.NoError(t, err)
	defer rc.Close()
	body, err := io.ReadAll(rc)
	require.NoError(t, err)
	assert.JSONEq(t, `{"answer":"hello"}`, string(body))
}

func TestFormatErr_DetailString(t *testing.T) {
	out := formatErr(ErrResponse{Detail: json.RawMessage(`"plain message"`)})
	assert.Equal(t, "plain message", out)
}

func TestFormatErr_DetailArray(t *testing.T) {
	out := formatErr(ErrResponse{Detail: json.RawMessage(`[{"loc":["body","name"],"msg":"required","type":"missing"}]`)})
	assert.Contains(t, out, "required")
}

func TestFormatErr_MessagePreferred(t *testing.T) {
	// /completions returns {"message": "..."} instead of {"detail": ...}
	out := formatErr(ErrResponse{Message: "Authentication failed for OpenAI"})
	assert.Equal(t, "Authentication failed for OpenAI", out)
}

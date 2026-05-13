package agentstudio

// Tests for the Agents API tag — the methods defined in agents.go.
//
// Cross-cutting concerns (NewClient validation, error mapping via
// checkResponse, header injection) live in client_test.go and use
// ListAgents as a vehicle since it's the simplest endpoint shape.

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

func TestListAgents_Success(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/1/agents", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Equal(t, "APP123", r.Header.Get(HeaderApplicationID))
		assert.Equal(t, "key-abc", r.Header.Get(HeaderAPIKey))
		assert.Equal(t, "cli-test", r.Header.Get(HeaderUserID))
		assert.Equal(t, "application/json", r.Header.Get("Accept"))

		assert.Equal(t, "2", r.URL.Query().Get("page"))
		assert.Equal(t, "25", r.URL.Query().Get("limit"))
		assert.Equal(t, "prov-1", r.URL.Query().Get("providerId"))

		require.NoError(t, json.NewEncoder(w).Encode(PaginatedAgentsResponse{
			Data: []Agent{{
				ID:           "11111111-1111-1111-1111-111111111111",
				Name:         "Concierge",
				Status:       StatusDraft,
				Instructions: "Be helpful.",
			}},
			Pagination: PaginationMetadata{
				Page: 2, Limit: 25, TotalCount: 1, TotalPages: 1,
			},
		}))
	})

	_, c := newTestClient(t, mux)

	got, err := c.ListAgents(context.Background(), ListAgentsParams{
		Page:       2,
		Limit:      25,
		ProviderID: "prov-1",
	})
	require.NoError(t, err)
	require.Len(t, got.Data, 1)
	assert.Equal(t, "Concierge", got.Data[0].Name)
	assert.Equal(t, StatusDraft, got.Data[0].Status)
	assert.Equal(t, 25, got.Pagination.Limit)
}

func TestListAgents_OmitsZeroParams(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/1/agents", func(w http.ResponseWriter, r *http.Request) {
		assert.Empty(t, r.URL.RawQuery, "expected no query params for zero-valued params")
		_, _ = w.Write([]byte(`{"data":[],"pagination":{"page":1,"limit":10,"totalCount":0,"totalPages":0}}`))
	})

	_, c := newTestClient(t, mux)
	_, err := c.ListAgents(context.Background(), ListAgentsParams{})
	require.NoError(t, err)
}

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

func TestCreateAgent_Success(t *testing.T) {
	wantBody := `{"name":"Concierge","instructions":"Be helpful."}`

	mux := http.NewServeMux()
	mux.HandleFunc("/1/agents", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))
		assert.Equal(t, "APP123", r.Header.Get(HeaderApplicationID))

		body, err := io.ReadAll(r.Body)
		require.NoError(t, err)
		assert.JSONEq(t, wantBody, string(body))

		w.WriteHeader(http.StatusCreated)
		_, _ = w.Write([]byte(`{
			"id":"abc-123",
			"name":"Concierge",
			"status":"draft",
			"instructions":"Be helpful.",
			"createdAt":"2025-01-01T00:00:00Z"
		}`))
	})

	_, c := newTestClient(t, mux)

	got, err := c.CreateAgent(context.Background(), json.RawMessage(wantBody))
	require.NoError(t, err)
	assert.Equal(t, "abc-123", got.ID)
	assert.Equal(t, "Concierge", got.Name)
	assert.Equal(t, StatusDraft, got.Status)
}

func TestCreateAgent_RejectsEmptyBody(t *testing.T) {
	_, c := newTestClient(t, http.NewServeMux())

	_, err := c.CreateAgent(context.Background(), nil)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "body is required")

	_, err = c.CreateAgent(context.Background(), json.RawMessage(""))
	require.Error(t, err)
	assert.Contains(t, err.Error(), "body is required")
}

func TestCreateAgent_PropagatesValidationError(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/1/agents", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusUnprocessableEntity)
		_, _ = w.Write([]byte(`{"detail":[{"msg":"name is required","loc":["body","name"]}]}`))
	})

	_, c := newTestClient(t, mux)

	_, err := c.CreateAgent(context.Background(), json.RawMessage(`{"instructions":"x"}`))
	require.Error(t, err)
	var apiErr *APIError
	require.True(t, errors.As(err, &apiErr))
	assert.Equal(t, http.StatusUnprocessableEntity, apiErr.StatusCode)
	assert.Equal(t, "name is required", apiErr.Detail)
}

func TestUpdateAgent_Success(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/1/agents/abc-123", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPatch, r.Method)
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))

		body, _ := io.ReadAll(r.Body)
		assert.JSONEq(t, `{"name":"Renamed"}`, string(body))

		_, _ = w.Write([]byte(`{
			"id":"abc-123",
			"name":"Renamed",
			"status":"draft",
			"instructions":"Be helpful.",
			"createdAt":"2025-01-01T00:00:00Z"
		}`))
	})

	_, c := newTestClient(t, mux)

	got, err := c.UpdateAgent(context.Background(), "abc-123", json.RawMessage(`{"name":"Renamed"}`))
	require.NoError(t, err)
	assert.Equal(t, "Renamed", got.Name)
}

func TestUpdateAgent_RejectsEmptyID(t *testing.T) {
	_, c := newTestClient(t, http.NewServeMux())
	_, err := c.UpdateAgent(context.Background(), "  ", json.RawMessage(`{"name":"x"}`))
	require.Error(t, err)
	assert.Contains(t, err.Error(), "agent id is required")
}

func TestUpdateAgent_PropagatesNotFound(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/1/agents/missing", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte(`{"detail":"Agent not found"}`))
	})

	_, c := newTestClient(t, mux)
	_, err := c.UpdateAgent(context.Background(), "missing", json.RawMessage(`{"name":"x"}`))
	require.Error(t, err)
	assert.True(t, errors.Is(err, ErrNotFound))
}

func TestDeleteAgent_NoContentSuccess(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/1/agents/abc-123", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodDelete, r.Method)
		w.WriteHeader(http.StatusNoContent)
	})

	_, c := newTestClient(t, mux)
	require.NoError(t, c.DeleteAgent(context.Background(), "abc-123"))
}

func TestDeleteAgent_NotFound(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/1/agents/missing", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte(`{"detail":"Agent not found"}`))
	})

	_, c := newTestClient(t, mux)
	err := c.DeleteAgent(context.Background(), "missing")
	require.Error(t, err)
	assert.True(t, errors.Is(err, ErrNotFound))
}

func TestDeleteAgent_RejectsEmptyID(t *testing.T) {
	_, c := newTestClient(t, http.NewServeMux())
	err := c.DeleteAgent(context.Background(), "")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "agent id is required")
}

func TestPublishUnpublishDuplicate(t *testing.T) {
	cases := []struct {
		name       string
		verb       string
		fn         func(*Client, context.Context, string) (*Agent, error)
		wantStatus AgentStatus
	}{
		{
			name:       "publish",
			verb:       "publish",
			fn:         (*Client).PublishAgent,
			wantStatus: StatusPublished,
		},
		{
			name:       "unpublish",
			verb:       "unpublish",
			fn:         (*Client).UnpublishAgent,
			wantStatus: StatusDraft,
		},
		{
			name:       "duplicate",
			verb:       "duplicate",
			fn:         (*Client).DuplicateAgent,
			wantStatus: StatusDraft,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			mux := http.NewServeMux()
			mux.HandleFunc(
				"/1/agents/abc-123/"+tc.verb,
				func(w http.ResponseWriter, r *http.Request) {
					assert.Equal(t, http.MethodPost, r.Method)
					// No request body for any of the lifecycle endpoints.
					assert.Empty(t, r.Header.Get("Content-Type"))
					body, _ := io.ReadAll(r.Body)
					assert.Empty(t, body)

					writeTestJSONResponse(w, []byte(`{
						"id":"new-id-or-same",
						"name":"X",
						"status":"`+string(tc.wantStatus)+`",
						"instructions":"x",
						"createdAt":"2025-01-01T00:00:00Z"
					}`))
				},
			)

			_, c := newTestClient(t, mux)

			got, err := tc.fn(c, context.Background(), "abc-123")
			require.NoError(t, err)
			assert.Equal(t, tc.wantStatus, got.Status)
		})
	}
}

func TestLifecycle_RejectsEmptyID(t *testing.T) {
	_, c := newTestClient(t, http.NewServeMux())

	for name, fn := range map[string]func(*Client, context.Context, string) (*Agent, error){
		"publish":   (*Client).PublishAgent,
		"unpublish": (*Client).UnpublishAgent,
		"duplicate": (*Client).DuplicateAgent,
	} {
		t.Run(name, func(t *testing.T) {
			_, err := fn(c, context.Background(), "")
			require.Error(t, err)
			assert.Contains(t, err.Error(), "agent id is required")
		})
	}
}

func TestInvalidateAgentCache(t *testing.T) {
	cases := []struct {
		name       string
		id         string
		before     string
		serverFn   func(t *testing.T) http.HandlerFunc
		wantErr    string // substring; "" = expect success
		isSentinel error
	}{
		{
			name:   "no before -> DELETE without query",
			id:     "abc-123",
			before: "",
			serverFn: func(t *testing.T) http.HandlerFunc {
				return func(w http.ResponseWriter, r *http.Request) {
					assert.Equal(t, http.MethodDelete, r.Method)
					assert.Equal(t, "", r.URL.RawQuery, "no before -> no query string")
					w.WriteHeader(http.StatusNoContent)
				}
			},
		},
		{
			name:   "with before -> DELETE with ?before",
			id:     "abc-123",
			before: "2026-01-15",
			serverFn: func(t *testing.T) http.HandlerFunc {
				return func(w http.ResponseWriter, r *http.Request) {
					assert.Equal(t, http.MethodDelete, r.Method)
					assert.Equal(t, "2026-01-15", r.URL.Query().Get("before"))
					w.WriteHeader(http.StatusNoContent)
				}
			},
		},
		{
			name:   "404 from backend surfaces as ErrNotFound",
			id:     "missing",
			before: "",
			serverFn: func(t *testing.T) http.HandlerFunc {
				return func(w http.ResponseWriter, _ *http.Request) {
					w.WriteHeader(http.StatusNotFound)
					_, _ = w.Write([]byte(`{"detail":"Agent not found"}`))
				}
			},
			wantErr:    "Agent not found",
			isSentinel: ErrNotFound,
		},
		{
			name:   "422 with structured detail (e.g. malformed before) surfaces backend message verbatim",
			id:     "abc-123",
			before: "not-a-date",
			serverFn: func(t *testing.T) http.HandlerFunc {
				return func(w http.ResponseWriter, _ *http.Request) {
					w.WriteHeader(http.StatusUnprocessableEntity)
					_, _ = w.Write(
						[]byte(
							`{"detail":[{"msg":"Input should be a valid date in YYYY-MM-DD format","loc":["query","before"]}]}`,
						),
					)
				}
			},
			wantErr: "valid date",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			mux := http.NewServeMux()
			mux.HandleFunc("/1/agents/"+tc.id+"/cache", tc.serverFn(t))
			_, c := newTestClient(t, mux)

			err := c.InvalidateAgentCache(context.Background(), tc.id, tc.before)
			if tc.wantErr == "" {
				require.NoError(t, err)
				return
			}
			require.Error(t, err)
			assert.Contains(t, err.Error(), tc.wantErr)
			if tc.isSentinel != nil {
				assert.True(t, errors.Is(err, tc.isSentinel))
			}
		})
	}
}

func TestInvalidateAgentCache_RejectsEmptyID(t *testing.T) {
	_, c := newTestClient(t, http.NewServeMux())
	err := c.InvalidateAgentCache(context.Background(), "  ", "")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "agent id is required")
}

func TestLifecycle_NotFound(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/1/agents/missing/publish", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte(`{"detail":"Agent not found"}`))
	})
	_, c := newTestClient(t, mux)
	_, err := c.PublishAgent(context.Background(), "missing")
	require.Error(t, err)
	assert.True(t, errors.Is(err, ErrNotFound))
}

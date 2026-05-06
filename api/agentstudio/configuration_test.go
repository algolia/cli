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

func TestGetConfiguration_Success(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/1/configuration", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		_, _ = w.Write([]byte(`{"maxRetentionDays":90}`))
	})
	_, c := newTestClient(t, mux)
	got, err := c.GetConfiguration(context.Background())
	require.NoError(t, err)
	assert.Equal(t, 90, got.MaxRetentionDays)
}

func TestGetConfiguration_PropagatesAPIError(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/1/configuration", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusForbidden)
		_, _ = w.Write([]byte(`{"detail":"insufficient permissions"}`))
	})
	_, c := newTestClient(t, mux)
	_, err := c.GetConfiguration(context.Background())
	require.Error(t, err)
	assert.True(t, errors.Is(err, ErrForbidden))
}

func TestUpdateConfiguration_RoundTripsBody(t *testing.T) {
	wire := `{"maxRetentionDays":30}`
	mux := http.NewServeMux()
	mux.HandleFunc("/1/configuration", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPatch, r.Method)
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))
		body, _ := io.ReadAll(r.Body)
		assert.JSONEq(t, wire, string(body))
		_, _ = w.Write([]byte(`{"maxRetentionDays":30}`))
	})
	_, c := newTestClient(t, mux)
	got, err := c.UpdateConfiguration(context.Background(), json.RawMessage(wire))
	require.NoError(t, err)
	assert.Equal(t, 30, got.MaxRetentionDays)
}

func TestUpdateConfiguration_RejectsEmptyBody(t *testing.T) {
	_, c := newTestClient(t, http.NewServeMux())
	_, err := c.UpdateConfiguration(context.Background(), nil)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "body is required")
}

func TestUpdateConfiguration_PropagatesValidationError(t *testing.T) {
	// The backend documents valid values [0, 30, 60, 90]. Anything
	// else returns 422 with structured detail.
	mux := http.NewServeMux()
	mux.HandleFunc("/1/configuration", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusUnprocessableEntity)
		_, _ = w.Write([]byte(`{"detail":[{"msg":"maxRetentionDays must be one of [0, 30, 60, 90]","loc":["body","maxRetentionDays"]}]}`))
	})
	_, c := newTestClient(t, mux)
	_, err := c.UpdateConfiguration(context.Background(), json.RawMessage(`{"maxRetentionDays":45}`))
	require.Error(t, err)
	var apiErr *APIError
	require.True(t, errors.As(err, &apiErr))
	assert.Equal(t, "maxRetentionDays must be one of [0, 30, 60, 90]", apiErr.Detail)
}

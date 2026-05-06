package agentstudio

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetUserData_Success(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/1/user-data/tok1", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		_, _ = w.Write([]byte(`{"conversations":[{"id":"c1"}],"memories":[{"id":"m1"}]}`))
	})
	_, c := newTestClient(t, mux)
	got, err := c.GetUserData(context.Background(), "tok1")
	require.NoError(t, err)
	assert.Len(t, got.Conversations, 1)
	assert.Len(t, got.Memories, 1)
}

func TestGetUserData_RejectsEmptyToken(t *testing.T) {
	_, c := newTestClient(t, http.NewServeMux())
	_, err := c.GetUserData(context.Background(), " ")
	require.Error(t, err)
}

func TestGetUserData_PathEscaped(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/1/user-data/tok%20with%20space", func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte(`{"conversations":[],"memories":[]}`))
	})
	_, c := newTestClient(t, mux)
	_, err := c.GetUserData(context.Background(), "tok with space")
	require.NoError(t, err)
}

func TestDeleteUserData_NoContent(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/1/user-data/tok1", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodDelete, r.Method)
		w.WriteHeader(http.StatusNoContent)
	})
	_, c := newTestClient(t, mux)
	require.NoError(t, c.DeleteUserData(context.Background(), "tok1"))
}

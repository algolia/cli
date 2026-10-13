package duplicate

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/algolia/cli/pkg/cmd/agents/sharedtest"
	"github.com/algolia/cli/test"
)

func Test_runDuplicateCmd(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/1/agents/abc-123/duplicate", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		_, _ = w.Write([]byte(`{
			"id":"new-id-456",
			"name":"Concierge (copy)",
			"status":"draft",
			"instructions":"x",
			"createdAt":"2025-01-02T00:00:00Z"
		}`))
	})

	ts := httptest.NewServer(mux)
	t.Cleanup(ts.Close)

	f, out := test.NewFactory(false, nil, nil, "")
	f.AgentStudioClient = sharedtest.NewClient(t, ts)

	cmd := NewDuplicateCmd(f, nil)
	result, err := test.Execute(cmd, "abc-123", out)
	require.NoError(t, err)
	assert.Contains(t, result.String(), `"id":"new-id-456"`)
	assert.Contains(t, result.String(), `"name":"Concierge (copy)"`)
}

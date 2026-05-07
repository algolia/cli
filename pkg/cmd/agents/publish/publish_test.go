package publish

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/algolia/cli/pkg/cmd/agents/sharedtest"
	"github.com/algolia/cli/test"
)

func Test_runPublishCmd(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/1/agents/abc-123/publish", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		_, _ = w.Write([]byte(`{
			"id":"abc-123",
			"name":"Concierge",
			"status":"published",
			"instructions":"x",
			"createdAt":"2025-01-01T00:00:00Z"
		}`))
	})

	ts := httptest.NewServer(mux)
	t.Cleanup(ts.Close)

	f, out := test.NewFactory(false, nil, nil, "")
	f.AgentStudioClient = sharedtest.NewClient(t, ts)

	cmd := NewPublishCmd(f, nil)
	result, err := test.Execute(cmd, "abc-123", out)
	require.NoError(t, err)
	assert.Contains(t, result.String(), `"status":"published"`)
}

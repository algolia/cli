package update

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/algolia/cli/pkg/cmd/agents/sharedtest"
	"github.com/algolia/cli/test"
)

func Test_runUpdateCmd_Success(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/1/agents/abc-123", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPatch, r.Method)
		body, _ := io.ReadAll(r.Body)
		assert.JSONEq(t, `{"name":"Renamed"}`, string(body))

		_, _ = w.Write([]byte(`{
			"id":"abc-123","name":"Renamed","status":"draft",
			"instructions":"x","createdAt":"2025-01-01T00:00:00Z"
		}`))
	})
	ts := httptest.NewServer(mux)
	t.Cleanup(ts.Close)

	f, out := test.NewFactory(false, nil, nil, `{"name":"Renamed"}`)
	f.AgentStudioClient = sharedtest.NewClient(t, ts)

	cmd := NewUpdateCmd(f, nil)
	result, err := test.Execute(cmd, "abc-123 -F -", out)
	require.NoError(t, err)
	assert.Contains(t, result.String(), `"name":"Renamed"`)
}

func Test_runUpdateCmd_RequiresAgentID(t *testing.T) {
	f, out := test.NewFactory(false, nil, nil, `{"name":"X"}`)
	cmd := NewUpdateCmd(f, nil)
	_, err := test.Execute(cmd, "-F -", out)
	require.Error(t, err)
	// validators.ExactArgs(1) message — not pinning exact wording, just that it errors.
}

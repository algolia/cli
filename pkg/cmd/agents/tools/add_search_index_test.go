package tools

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

func Test_runAddSearchIndexCmd_createsTool(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/1/agents/agent-1", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			_, _ = w.Write([]byte(`{
				"id":"agent-1","name":"x","status":"draft","instructions":"hi",
				"createdAt":"2025-01-01T00:00:00Z"
			}`))
		case http.MethodPatch:
			body, _ := io.ReadAll(r.Body)
			assert.JSONEq(t,
				`{"tools":[{"indices":[{"description":"Catalog","index":"PRODUCTS"}],"name":"products","type":"algolia_search_index"}]}`,
				string(body))
			_, _ = w.Write([]byte(`{
				"id":"agent-1","name":"x","status":"draft","instructions":"hi",
				"tools":[{"name":"products","type":"algolia_search_index","indices":[{"index":"PRODUCTS","description":"Catalog"}]}],
				"createdAt":"2025-01-01T00:00:00Z"
			}`))
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})
	ts := httptest.NewServer(mux)
	t.Cleanup(ts.Close)

	f, out := test.NewFactory(false, nil, nil, "")
	f.AgentStudioClient = sharedtest.NewClient(t, ts)

	cmd := newAddSearchIndexCmd(f, nil)
	result, err := test.Execute(cmd, `agent-1 --index PRODUCTS --description Catalog`, out)
	require.NoError(t, err)
	assert.Contains(t, result.String(), `"type":"algolia_search_index"`)
}

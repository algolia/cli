package list

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/algolia/cli/pkg/cmd/agents/sharedtest"
	"github.com/algolia/cli/test"
)

// newCmdAgainst wires a fresh agentstudio client (pointed at a
// httptest server) onto the standard test factory and returns an
// executor for `algolia agents list <args>`.
func newCmdAgainst(
	t *testing.T,
	isTTY bool,
	handler http.Handler,
) func(args string) (*test.CmdInOut, error) {
	t.Helper()

	ts := httptest.NewServer(handler)
	t.Cleanup(ts.Close)

	f, out := test.NewFactory(isTTY, nil, nil, "")
	f.AgentStudioClient = sharedtest.NewClient(t, ts)

	return func(args string) (*test.CmdInOut, error) {
		cmd := NewListCmd(f, nil)
		return test.Execute(cmd, args, out)
	}
}

func Test_runListCmd(t *testing.T) {
	// Freeze time so "ago" formatting is deterministic.
	oldNowFn := nowFn
	nowFn = func() time.Time { return time.Unix(1735689600, 0) } // 2025-01-01T00:00:00Z
	t.Cleanup(func() { nowFn = oldNowFn })

	tests := []struct {
		name      string
		isTTY     bool
		args      string
		wantOut   string
		wantQuery string
	}{
		{
			name:      "non-tty defaults",
			isTTY:     false,
			args:      "",
			wantOut:   "abc-123\tConcierge\tpublished\tprov-1\t1 year ago\n",
			wantQuery: "",
		},
		{
			name:      "non-tty with paging + provider filter",
			isTTY:     false,
			args:      "--page 2 --per-page 25 --provider-id prov-1",
			wantQuery: "limit=25&page=2&providerId=prov-1",
			wantOut:   "abc-123\tConcierge\tpublished\tprov-1\t1 year ago\n",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			updated := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)

			handler := http.NewServeMux()
			handler.HandleFunc("/1/agents", func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, tc.wantQuery, r.URL.RawQuery)
				_, _ = w.Write([]byte(`{
					"data":[{
						"id":"abc-123",
						"name":"Concierge",
						"status":"published",
						"providerId":"prov-1",
						"instructions":"Be helpful.",
						"createdAt":"2023-01-01T00:00:00Z",
						"updatedAt":"` + updated.Format(time.RFC3339) + `"
					}],
					"pagination":{"page":1,"limit":10,"totalCount":1,"totalPages":1}
				}`))
			})

			exec := newCmdAgainst(t, tc.isTTY, handler)
			result, err := exec(tc.args)
			require.NoError(t, err)
			assert.Equal(t, tc.wantOut, result.String())
		})
	}
}

func Test_runListCmd_outputJSON(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/1/agents", func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte(`{
			"data":[{
				"id":"abc-123",
				"name":"Concierge",
				"status":"draft",
				"instructions":"Be helpful.",
				"createdAt":"2025-01-01T00:00:00Z"
			}],
			"pagination":{"page":1,"limit":10,"totalCount":1,"totalPages":1}
		}`))
	})

	exec := newCmdAgainst(t, false, mux)
	result, err := exec("--output json")
	require.NoError(t, err)
	assert.Contains(t, result.String(), `"id":"abc-123"`)
	assert.Contains(t, result.String(), `"name":"Concierge"`)
	assert.Contains(t, result.String(), `"pagination"`)
}

func Test_runListCmd_PropagatesAPIError(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/1/agents", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusForbidden)
		_, _ = w.Write([]byte(`{"message":"This feature is not enabled for this application."}`))
	})

	exec := newCmdAgainst(t, false, mux)
	_, err := exec("")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "feature is not enabled")
}

package conversations

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/algolia/cli/test"
)

func Test_runExportCmd_StdoutPathPrettyPrints(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/1/agents/agent-1/conversations/export", func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte(`[{"id":"c1"},{"id":"c2"}]`))
	})
	ts := httptest.NewServer(mux)
	t.Cleanup(ts.Close)

	f, out := test.NewFactory(false, nil, nil, "")
	f.AgentStudioClient = newClientForServer(t, ts)

	cmd := NewConversationsCmd(f)
	result, err := test.Execute(cmd, "export agent-1", out)
	require.NoError(t, err)
	got := result.String()
	// Pretty-printed: 2-space indent, one element per line.
	assert.Contains(t, got, "  {")
	assert.Contains(t, got, `"id": "c1"`)
}

func Test_runExportCmd_OutputFileWritesCompact(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/1/agents/agent-1/conversations/export", func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte(`[{"id":"c1"}]`))
	})
	ts := httptest.NewServer(mux)
	t.Cleanup(ts.Close)

	dir := t.TempDir()
	outPath := filepath.Join(dir, "export.json")

	f, out := test.NewFactory(false, nil, nil, "")
	f.AgentStudioClient = newClientForServer(t, ts)

	cmd := NewConversationsCmd(f)
	_, err := test.Execute(cmd, "export agent-1 -O "+outPath, out)
	require.NoError(t, err)

	contents, err := os.ReadFile(outPath)
	require.NoError(t, err)
	// Files get the raw bytes (jq-friendly), not the pretty form.
	assert.Equal(t, `[{"id":"c1"}]`, string(contents))
}

func Test_runExportCmd_PassesDateRange(t *testing.T) {
	mux := http.NewServeMux()
	hit := false
	mux.HandleFunc("/1/agents/agent-1/conversations/export", func(w http.ResponseWriter, r *http.Request) {
		hit = true
		assert.Equal(t, "2026-01-01", r.URL.Query().Get("startDate"))
		assert.Equal(t, "2026-01-31", r.URL.Query().Get("endDate"))
		_, _ = w.Write([]byte(`[]`))
	})
	ts := httptest.NewServer(mux)
	t.Cleanup(ts.Close)

	f, out := test.NewFactory(false, nil, nil, "")
	f.AgentStudioClient = newClientForServer(t, ts)

	cmd := NewConversationsCmd(f)
	_, err := test.Execute(cmd, "export agent-1 --start-date 2026-01-01 --end-date 2026-01-31", out)
	require.NoError(t, err)
	assert.True(t, hit)
}

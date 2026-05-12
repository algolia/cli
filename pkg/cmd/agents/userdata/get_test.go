package userdata

import (
	"io/fs"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/algolia/cli/pkg/cmd/agents/sharedtest"
	"github.com/algolia/cli/test"
)

func Test_runGetCmd_Stdout(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/1/user-data/tok1", func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte(`{"conversations":[{"id":"c1"}],"memories":[{"id":"m1"}]}`))
	})
	ts := httptest.NewServer(mux)
	t.Cleanup(ts.Close)

	f, out := test.NewFactory(false, nil, nil, "")
	f.AgentStudioClient = sharedtest.NewClient(t, ts)
	cmd := NewUserDataCmd(f)
	result, err := test.Execute(cmd, "get tok1", out)
	require.NoError(t, err)
	got := result.String()
	assert.Contains(t, got, `"conversations"`)
	assert.Contains(t, got, `"c1"`)
	assert.Contains(t, got, `"m1"`)
}

func Test_runGetCmd_OutputFile(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/1/user-data/tok1", func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte(`{"conversations":[],"memories":[]}`))
	})
	ts := httptest.NewServer(mux)
	t.Cleanup(ts.Close)

	const outFile = "out.json"
	dir := t.TempDir()
	dst := filepath.Join(dir, outFile)
	f, out := test.NewFactory(false, nil, nil, "")
	f.AgentStudioClient = sharedtest.NewClient(t, ts)
	cmd := NewUserDataCmd(f)
	_, err := test.Execute(cmd, "get tok1 -o "+dst, out)
	require.NoError(t, err)
	body, err := fs.ReadFile(os.DirFS(dir), outFile)
	require.NoError(t, err)
	assert.Contains(t, string(body), `"conversations"`)
}

func Test_runGetCmd_RequiresToken(t *testing.T) {
	f, out := test.NewFactory(false, nil, nil, "")
	cmd := NewUserDataCmd(f)
	_, err := test.Execute(cmd, "get", out)
	require.Error(t, err)
}

func Test_runGetCmd_RejectsSlashInToken(t *testing.T) {
	f, out := test.NewFactory(false, nil, nil, "")
	cmd := NewUserDataCmd(f)
	_, err := test.Execute(cmd, `get ab/cd`, out)
	require.Error(t, err)
	assert.Contains(t, err.Error(), `contains "/"`)
}

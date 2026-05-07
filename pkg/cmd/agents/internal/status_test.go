package internal

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/algolia/cli/pkg/cmd/agents/sharedtest"
	"github.com/algolia/cli/test"
)

func Test_runStatusCmd_Success(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/status", func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte(`{"status":"ok","version":"abc","migration_revision":"r1"}`))
	})
	ts := httptest.NewServer(mux)
	t.Cleanup(ts.Close)

	f, out := test.NewFactory(false, nil, nil, "")
	f.AgentStudioClient = sharedtest.NewClient(t, ts)
	cmd := NewInternalCmd(f)
	result, err := test.Execute(cmd, "status", out)
	require.NoError(t, err)
	got := result.String()
	assert.Contains(t, got, `"status"`)
	assert.Contains(t, got, `"version"`)
}

func Test_InternalCmd_ParentHiddenButSubsVisible(t *testing.T) {
	// Parent is hidden so `algolia agents --help` doesn't list "internal";
	// subcommands stay visible so `algolia agents internal --help`
	// surfaces what's available to anyone who knows to look.
	f, _ := test.NewFactory(false, nil, nil, "")
	cmd := NewInternalCmd(f)
	assert.True(t, cmd.Hidden, "internal parent should be hidden")
	for _, sub := range cmd.Commands() {
		assert.False(t, sub.Hidden, "subcommand %q should NOT be hidden (parent is)", sub.Name())
	}
}

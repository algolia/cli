package keys

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/algolia/cli/pkg/cmd/agents/sharedtest"
	"github.com/algolia/cli/test"
)

func Test_runDeleteCmd_NonTTYWithoutConfirmFails(t *testing.T) {
	f, out := test.NewFactory(false, nil, nil, "")
	cmd := NewKeysCmd(f)
	_, err := test.Execute(cmd, "delete id1", out)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "--confirm required")
}

func Test_runDeleteCmd_Live(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/1/secret-keys/id1", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	})
	ts := httptest.NewServer(mux)
	t.Cleanup(ts.Close)

	f, out := test.NewFactory(false, nil, nil, "")
	f.AgentStudioClient = sharedtest.NewClient(t, ts)
	cmd := NewKeysCmd(f)
	_, err := test.Execute(cmd, "delete id1 -y", out)
	require.NoError(t, err)
}

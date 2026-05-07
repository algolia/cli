package keys

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/algolia/cli/test"
)

func Test_runListCmd_MasksValueByDefault(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/1/secret-keys", func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte(`{"data":[{"id":"id1","name":"k1","value":"sk-real","createdAt":"2026-01-01T00:00:00Z","updatedAt":"2026-01-01T00:00:00Z","lastUsedAt":null,"isDefault":false,"agentIds":[]}],"pagination":{"page":1,"limit":10,"totalCount":1,"totalPages":1}}`))
	})
	ts := httptest.NewServer(mux)
	t.Cleanup(ts.Close)

	f, out := test.NewFactory(false, nil, nil, "")
	f.AgentStudioClient = newClientForServer(t, ts)
	cmd := NewKeysCmd(f)
	result, err := test.Execute(cmd, "list --output json", out)
	require.NoError(t, err)
	got := result.String()
	assert.Contains(t, got, `"***"`)
	assert.NotContains(t, got, "sk-real")
}

func Test_runListCmd_ShowSecretRevealsValue(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/1/secret-keys", func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte(`{"data":[{"id":"id1","name":"k1","value":"sk-real","createdAt":"2026-01-01T00:00:00Z","updatedAt":"2026-01-01T00:00:00Z","lastUsedAt":null,"isDefault":false,"agentIds":[]}],"pagination":{"page":1,"limit":10,"totalCount":1,"totalPages":1}}`))
	})
	ts := httptest.NewServer(mux)
	t.Cleanup(ts.Close)

	f, out := test.NewFactory(false, nil, nil, "")
	f.AgentStudioClient = newClientForServer(t, ts)
	cmd := NewKeysCmd(f)
	result, err := test.Execute(cmd, "list --show-secret --output json", out)
	require.NoError(t, err)
	assert.Contains(t, result.String(), "sk-real")
}

// TTY render path: the table writes k.Value verbatim, so masking must
// happen before the row is added (it does, in maskKey). Locks that in.
func Test_runListCmd_TableMasksValueByDefault(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/1/secret-keys", func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte(`{"data":[{"id":"id1","name":"k1","value":"sk-real","createdAt":"2026-01-01T00:00:00Z","updatedAt":"2026-01-01T00:00:00Z","lastUsedAt":null,"isDefault":false,"agentIds":[]}],"pagination":{"page":1,"limit":10,"totalCount":1,"totalPages":1}}`))
	})
	ts := httptest.NewServer(mux)
	t.Cleanup(ts.Close)

	f, out := test.NewFactory(true, nil, nil, "")
	f.AgentStudioClient = newClientForServer(t, ts)
	cmd := NewKeysCmd(f)
	result, err := test.Execute(cmd, "list", out)
	require.NoError(t, err)
	got := result.String()
	assert.Contains(t, got, "***")
	assert.NotContains(t, got, "sk-real")
}

func Test_runListCmd_TableShowSecretRevealsValue(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/1/secret-keys", func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte(`{"data":[{"id":"id1","name":"k1","value":"sk-real","createdAt":"2026-01-01T00:00:00Z","updatedAt":"2026-01-01T00:00:00Z","lastUsedAt":null,"isDefault":false,"agentIds":[]}],"pagination":{"page":1,"limit":10,"totalCount":1,"totalPages":1}}`))
	})
	ts := httptest.NewServer(mux)
	t.Cleanup(ts.Close)

	f, out := test.NewFactory(true, nil, nil, "")
	f.AgentStudioClient = newClientForServer(t, ts)
	cmd := NewKeysCmd(f)
	result, err := test.Execute(cmd, "list --show-secret", out)
	require.NoError(t, err)
	assert.Contains(t, result.String(), "sk-real")
}

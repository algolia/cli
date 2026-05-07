package domains

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/algolia/cli/test"
)

func writeTemp(t *testing.T, content string) string {
	t.Helper()
	p := filepath.Join(t.TempDir(), "in.json")
	require.NoError(t, os.WriteFile(p, []byte(content), 0o600))
	return p
}

func Test_runBulkInsertCmd_DomainAndFileMutex(t *testing.T) {
	p := writeTemp(t, `["a"]`)
	f, out := test.NewFactory(false, nil, nil, "")
	cmd := NewDomainsCmd(f)
	_, err := test.Execute(cmd, "bulk-insert agent-1 --domain a -F "+p, out)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "mutually exclusive")
}

func Test_runBulkInsertCmd_RequiresInput(t *testing.T) {
	f, out := test.NewFactory(false, nil, nil, "")
	cmd := NewDomainsCmd(f)
	_, err := test.Execute(cmd, "bulk-insert agent-1", out)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "provide at least one")
}

func Test_runBulkInsertCmd_DryRunListsDomains(t *testing.T) {
	f, out := test.NewFactory(false, nil, nil, "")
	cmd := NewDomainsCmd(f)
	result, err := test.Execute(cmd, "bulk-insert agent-1 --domain a --domain b --dry-run", out)
	require.NoError(t, err)
	got := result.String()
	assert.Contains(t, got, "Dry run: would POST /1/agents/agent-1/allowed-domains/bulk")
	assert.Contains(t, got, "domains (2): a, b")
}

func Test_runBulkInsertCmd_Live(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/1/agents/agent-1/allowed-domains/bulk", func(w http.ResponseWriter, r *http.Request) {
		var got map[string][]string
		require.NoError(t, json.NewDecoder(r.Body).Decode(&got))
		assert.Equal(t, []string{"a", "b"}, got["domains"])
		w.WriteHeader(http.StatusCreated)
		_, _ = w.Write([]byte(`{"domains":[]}`))
	})
	ts := httptest.NewServer(mux)
	t.Cleanup(ts.Close)

	f, out := test.NewFactory(false, nil, nil, "")
	f.AgentStudioClient = newClientForServer(t, ts)
	cmd := NewDomainsCmd(f)
	_, err := test.Execute(cmd, "bulk-insert agent-1 --domain a --domain b", out)
	require.NoError(t, err)
}

func Test_runBulkInsertCmd_FileWithEmptyArray(t *testing.T) {
	p := writeTemp(t, `[]`)
	f, out := test.NewFactory(false, nil, nil, "")
	cmd := NewDomainsCmd(f)
	_, err := test.Execute(cmd, "bulk-insert agent-1 -F "+p, out)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "array is empty")
}

func Test_runBulkDeleteCmd_NonTTYWithoutConfirmFails(t *testing.T) {
	f, out := test.NewFactory(false, nil, nil, "")
	cmd := NewDomainsCmd(f)
	_, err := test.Execute(cmd, "bulk-delete agent-1 --domain-id d1", out)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "--confirm required")
}

func Test_runBulkDeleteCmd_DryRun(t *testing.T) {
	f, out := test.NewFactory(false, nil, nil, "")
	cmd := NewDomainsCmd(f)
	result, err := test.Execute(cmd, "bulk-delete agent-1 --domain-id d1 --domain-id d2 --dry-run", out)
	require.NoError(t, err)
	got := result.String()
	assert.Contains(t, got, "Dry run: would DELETE /1/agents/agent-1/allowed-domains/bulk")
	assert.Contains(t, got, "ids (2): d1, d2")
}

func Test_runBulkDeleteCmd_Live(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/1/agents/agent-1/allowed-domains/bulk", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodDelete, r.Method)
		var got map[string][]string
		require.NoError(t, json.NewDecoder(r.Body).Decode(&got))
		assert.Equal(t, []string{"d1"}, got["domainIds"])
		w.WriteHeader(http.StatusNoContent)
	})
	ts := httptest.NewServer(mux)
	t.Cleanup(ts.Close)

	f, out := test.NewFactory(false, nil, nil, "")
	f.AgentStudioClient = newClientForServer(t, ts)
	cmd := NewDomainsCmd(f)
	_, err := test.Execute(cmd, "bulk-delete agent-1 --domain-id d1 -y", out)
	require.NoError(t, err)
}

// TTY mode triggers a pre-fetch so the success line can split requested
// IDs into removed (present at fetch time) vs already absent. Two of the
// requested IDs exist on the server, one does not.
func Test_runBulkDeleteCmd_TTYReportsRemovedAndAbsent(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/1/agents/agent-1/allowed-domains", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		_, _ = w.Write([]byte(`{"domains":[{"id":"d1","domain":"a.test"},{"id":"d2","domain":"b.test"}]}`))
	})
	mux.HandleFunc("/1/agents/agent-1/allowed-domains/bulk", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodDelete, r.Method)
		w.WriteHeader(http.StatusNoContent)
	})
	ts := httptest.NewServer(mux)
	t.Cleanup(ts.Close)

	f, out := test.NewFactory(true, nil, nil, "")
	f.AgentStudioClient = newClientForServer(t, ts)
	cmd := NewDomainsCmd(f)
	result, err := test.Execute(cmd, "bulk-delete agent-1 --domain-id d1 --domain-id d2 --domain-id d-missing -y", out)
	require.NoError(t, err)
	got := result.String()
	assert.Contains(t, got, "requested 3 (removed 2, already absent 1)")
}

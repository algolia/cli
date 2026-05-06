package providers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/algolia/cli/api/agentstudio"
	"github.com/algolia/cli/test"
)

func newClientForServer(t *testing.T, ts *httptest.Server) func() (*agentstudio.Client, error) {
	t.Helper()
	return func() (*agentstudio.Client, error) {
		return agentstudio.NewClient(agentstudio.Config{
			BaseURL:       ts.URL,
			ApplicationID: "APP123",
			APIKey:        "k",
			HTTPClient:    ts.Client(),
		})
	}
}

func writeTempJSON(t *testing.T, name, content string) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), name)
	require.NoError(t, os.WriteFile(path, []byte(content), 0o600))
	return path
}

// ---------------------------------------------------------------------
// list
// ---------------------------------------------------------------------

func Test_runListCmd_MasksSecretsByDefault(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/1/providers", func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte(`{
			"data":[{
				"id":"p1","name":"openai-prod","providerName":"openai",
				"input":{"apiKey":"sk-LEAKED-DATA"},
				"createdAt":"2026-01-01T00:00:00Z","updatedAt":"2026-01-01T00:00:00Z"
			}],
			"pagination":{"page":1,"limit":10,"totalCount":1,"totalPages":1}
		}`))
	})
	ts := httptest.NewServer(mux)
	t.Cleanup(ts.Close)

	f, out := test.NewFactory(false, nil, nil, "")
	f.AgentStudioClient = newClientForServer(t, ts)

	cmd := NewProvidersCmd(f)
	result, err := test.Execute(cmd, "list --output json", out)
	require.NoError(t, err)
	got := result.String()
	assert.NotContains(t, got, "sk-LEAKED-DATA", "raw apiKey must NOT appear in default output")
	assert.Contains(t, got, `"apiKey":"***"`, "apiKey must be replaced with masking sentinel")
}

func Test_runListCmd_ShowSecretRevealsRawKey(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/1/providers", func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte(`{
			"data":[{
				"id":"p1","name":"openai-prod","providerName":"openai",
				"input":{"apiKey":"sk-INTENTIONAL-EXPORT"},
				"createdAt":"2026-01-01T00:00:00Z","updatedAt":"2026-01-01T00:00:00Z"
			}],
			"pagination":{"page":1,"limit":10,"totalCount":1,"totalPages":1}
		}`))
	})
	ts := httptest.NewServer(mux)
	t.Cleanup(ts.Close)

	f, out := test.NewFactory(false, nil, nil, "")
	f.AgentStudioClient = newClientForServer(t, ts)

	cmd := NewProvidersCmd(f)
	result, err := test.Execute(cmd, "list --output json --show-secret", out)
	require.NoError(t, err)
	assert.Contains(t, result.String(), "sk-INTENTIONAL-EXPORT")
}

// ---------------------------------------------------------------------
// get
// ---------------------------------------------------------------------

func Test_runGetCmd_MasksByDefault(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/1/providers/p1", func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte(`{
			"id":"p1","name":"openai-prod","providerName":"openai",
			"input":{"apiKey":"sk-LEAKED"},
			"createdAt":"2026-01-01T00:00:00Z","updatedAt":"2026-01-01T00:00:00Z"
		}`))
	})
	ts := httptest.NewServer(mux)
	t.Cleanup(ts.Close)

	f, out := test.NewFactory(false, nil, nil, "")
	f.AgentStudioClient = newClientForServer(t, ts)

	cmd := NewProvidersCmd(f)
	result, err := test.Execute(cmd, "get p1", out)
	require.NoError(t, err)
	assert.NotContains(t, result.String(), "sk-LEAKED")
	assert.Contains(t, result.String(), `"apiKey":"***"`)
}

func Test_runGetCmd_RequiresProviderID(t *testing.T) {
	f, out := test.NewFactory(false, nil, nil, "")
	cmd := NewProvidersCmd(f)
	_, err := test.Execute(cmd, "get", out)
	require.Error(t, err)
}

// ---------------------------------------------------------------------
// create
// ---------------------------------------------------------------------

func Test_runCreateCmd_RoundTripsBody(t *testing.T) {
	body := `{"name":"openai-prod","providerName":"openai","input":{"apiKey":"sk-XYZ"}}`
	mux := http.NewServeMux()
	mux.HandleFunc("/1/providers", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		w.WriteHeader(http.StatusCreated)
		_, _ = w.Write([]byte(`{
			"id":"p1","name":"openai-prod","providerName":"openai",
			"input":{"apiKey":"sk-XYZ"},
			"createdAt":"2026-01-01T00:00:00Z","updatedAt":"2026-01-01T00:00:00Z"
		}`))
	})
	ts := httptest.NewServer(mux)
	t.Cleanup(ts.Close)

	specPath := writeTempJSON(t, "spec.json", body)

	f, out := test.NewFactory(false, nil, nil, "")
	f.AgentStudioClient = newClientForServer(t, ts)

	cmd := NewProvidersCmd(f)
	result, err := test.Execute(cmd, "create -F "+specPath, out)
	require.NoError(t, err)
	// Default success output is JSON; secrets masked.
	assert.NotContains(t, result.String(), "sk-XYZ")
	assert.Contains(t, result.String(), `"apiKey":"***"`)
}

func Test_runCreateCmd_DryRunSkipsAPI(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/1/providers", func(_ http.ResponseWriter, _ *http.Request) {
		t.Fatal("backend was called during --dry-run")
	})
	ts := httptest.NewServer(mux)
	t.Cleanup(ts.Close)

	specPath := writeTempJSON(t, "spec.json", `{"name":"x","providerName":"openai","input":{"apiKey":"sk-x"}}`)
	f, out := test.NewFactory(false, nil, nil, "")
	f.AgentStudioClient = newClientForServer(t, ts)

	cmd := NewProvidersCmd(f)
	result, err := test.Execute(cmd, "create -F "+specPath+" --dry-run", out)
	require.NoError(t, err)
	got := result.String()
	assert.Contains(t, got, "Dry run: would POST /1/providers")
	// Dry-run shows the unmodified body — masking does NOT apply to
	// dry-run because the user authored the file and is being shown
	// what THEY are about to send.
	assert.Contains(t, got, "sk-x")
}

func Test_runCreateCmd_RejectsInvalidJSON(t *testing.T) {
	specPath := writeTempJSON(t, "spec.json", `{not json`)
	f, out := test.NewFactory(false, nil, nil, "")
	cmd := NewProvidersCmd(f)
	_, err := test.Execute(cmd, "create -F "+specPath+" --dry-run", out)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "not valid JSON")
}

// ---------------------------------------------------------------------
// update
// ---------------------------------------------------------------------

func Test_runUpdateCmd_HitsCorrectPath(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/1/providers/p1", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPatch, r.Method)
		_, _ = w.Write([]byte(`{
			"id":"p1","name":"renamed","providerName":"openai",
			"input":{"apiKey":"sk-XYZ"},
			"createdAt":"2026-01-01T00:00:00Z","updatedAt":"2026-01-02T00:00:00Z"
		}`))
	})
	ts := httptest.NewServer(mux)
	t.Cleanup(ts.Close)

	patchPath := writeTempJSON(t, "patch.json", `{"name":"renamed"}`)
	f, out := test.NewFactory(false, nil, nil, "")
	f.AgentStudioClient = newClientForServer(t, ts)

	cmd := NewProvidersCmd(f)
	result, err := test.Execute(cmd, "update p1 -F "+patchPath, out)
	require.NoError(t, err)
	assert.Contains(t, result.String(), "renamed")
	assert.NotContains(t, result.String(), "sk-XYZ")
}

func Test_runUpdateCmd_DryRunIncludesProviderID(t *testing.T) {
	patchPath := writeTempJSON(t, "patch.json", `{"name":"renamed"}`)
	f, out := test.NewFactory(false, nil, nil, "")
	cmd := NewProvidersCmd(f)
	result, err := test.Execute(cmd, "update p1 -F "+patchPath+" --dry-run", out)
	require.NoError(t, err)
	assert.Contains(t, result.String(), "Dry run: would PATCH /1/providers/p1")
}

// ---------------------------------------------------------------------
// delete
// ---------------------------------------------------------------------

func Test_runDeleteCmd_DryRunPreFetchesAndPreviews(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/1/providers/p1", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method, "dry-run must GET, not DELETE")
		_, _ = w.Write([]byte(`{
			"id":"p1","name":"openai-prod","providerName":"openai",
			"input":{"apiKey":"sk-x"},
			"createdAt":"2026-01-01T00:00:00Z","updatedAt":"2026-01-01T00:00:00Z"
		}`))
	})
	ts := httptest.NewServer(mux)
	t.Cleanup(ts.Close)

	f, out := test.NewFactory(false, nil, nil, "")
	f.AgentStudioClient = newClientForServer(t, ts)

	cmd := NewProvidersCmd(f)
	result, err := test.Execute(cmd, "delete p1 --dry-run", out)
	require.NoError(t, err)
	got := result.String()
	assert.Contains(t, got, "Dry run: would DELETE /1/providers/p1")
	assert.Contains(t, got, "openai-prod")
	assert.Contains(t, got, "openai")
}

func Test_runDeleteCmd_NonTTYWithoutConfirmFails(t *testing.T) {
	f, out := test.NewFactory(false, nil, nil, "")
	cmd := NewProvidersCmd(f)
	_, err := test.Execute(cmd, "delete p1", out)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "--confirm required")
}

func Test_runDeleteCmd_PropagatesConflict(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/1/providers/p1", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			_, _ = w.Write([]byte(`{
				"id":"p1","name":"openai-prod","providerName":"openai",
				"input":{"apiKey":"sk-x"},
				"createdAt":"2026-01-01T00:00:00Z","updatedAt":"2026-01-01T00:00:00Z"
			}`))
		case http.MethodDelete:
			w.WriteHeader(http.StatusConflict)
			_, _ = w.Write([]byte(`{"detail":"Provider is in use by 3 agent(s)"}`))
		}
	})
	ts := httptest.NewServer(mux)
	t.Cleanup(ts.Close)

	f, out := test.NewFactory(false, nil, nil, "")
	f.AgentStudioClient = newClientForServer(t, ts)

	cmd := NewProvidersCmd(f)
	_, err := test.Execute(cmd, "delete p1 -y", out)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "in use by 3 agent")
}

// ---------------------------------------------------------------------
// models
// ---------------------------------------------------------------------

func Test_runModelsCmd_NoFlag_HitsCatalogRoute(t *testing.T) {
	mux := http.NewServeMux()
	hit := false
	mux.HandleFunc("/1/providers/models", func(w http.ResponseWriter, r *http.Request) {
		hit = true
		assert.Equal(t, http.MethodGet, r.Method)
		_, _ = w.Write([]byte(`{"openai":["gpt-4o","gpt-4o-mini"],"anthropic":["claude-3-5-sonnet"]}`))
	})
	ts := httptest.NewServer(mux)
	t.Cleanup(ts.Close)

	f, out := test.NewFactory(false, nil, nil, "")
	f.AgentStudioClient = newClientForServer(t, ts)

	cmd := NewProvidersCmd(f)
	result, err := test.Execute(cmd, "models --output json", out)
	require.NoError(t, err)
	assert.True(t, hit)

	var got map[string][]string
	require.NoError(t, json.Unmarshal([]byte(result.String()), &got))
	assert.Contains(t, got, "openai")
	assert.Contains(t, got["openai"], "gpt-4o")
}

func Test_runModelsCmd_WithFlag_HitsConfiguredProviderRoute(t *testing.T) {
	mux := http.NewServeMux()
	hit := false
	mux.HandleFunc("/1/providers/p1/models", func(w http.ResponseWriter, r *http.Request) {
		hit = true
		assert.Equal(t, http.MethodGet, r.Method)
		_, _ = w.Write([]byte(`["gpt-4o","my-fine-tune-1"]`))
	})
	ts := httptest.NewServer(mux)
	t.Cleanup(ts.Close)

	f, out := test.NewFactory(false, nil, nil, "")
	f.AgentStudioClient = newClientForServer(t, ts)

	cmd := NewProvidersCmd(f)
	result, err := test.Execute(cmd, "models --provider-id p1 --output json", out)
	require.NoError(t, err)
	assert.True(t, hit)
	assert.Contains(t, result.String(), "my-fine-tune-1")
}

func Test_runModelsCmd_PropagatesAPIError(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/1/providers/missing/models", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte(`{"detail":"Provider not found"}`))
	})
	ts := httptest.NewServer(mux)
	t.Cleanup(ts.Close)

	f, out := test.NewFactory(false, nil, nil, "")
	f.AgentStudioClient = newClientForServer(t, ts)

	cmd := NewProvidersCmd(f)
	_, err := test.Execute(cmd, "models --provider-id missing", out)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "Provider not found")
}
